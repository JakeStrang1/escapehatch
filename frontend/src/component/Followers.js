import React from "react"
import NavBar from "./NavBar"
import UserSummary from "./UserSummary"
import Shelves from "./Shelves"
import Container from 'react-bootstrap/Container'
import Col from 'react-bootstrap/Col'
import Dropdown from 'react-bootstrap/Dropdown'
import Button from 'react-bootstrap/Button'
import Badge from 'react-bootstrap/Badge'
import Row from 'react-bootstrap/Row'
import Form from 'react-bootstrap/Form'
import { Redirect, Link } from 'react-router-dom'
import api, {
  ERR_UNEXPECTED,
} from "../api"
import { connect } from './../reducers'

export default class Followers extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      authError: false,
      followers: [],
      following: [],
      search: "",
      followerCount: 0,
      followingCount: 0,
    }

    this.GetDefaultResults = this.GetDefaultResults.bind(this)
    this.handleSearchSubmit = this.handleSearchSubmit.bind(this)
    this.handleSearchChange = this.handleSearchChange.bind(this)

    this.GetDefaultResults()
  }

  GetDefaultResults() {
    api.GetFollowers("me") // TODO: Fetch more than just the first page of results!
    .then(response => {
      if (response.ok) {
        this.setState({followers: response.body.data, followerCount: response.body.pages.per_page * response.body.pages.total_pages})
        return
      }
      
      if (response.status == 401) {
        this.setState({authError: true})
      }
      console.log("Status: " + response.status + ", Code: " + response.errorCode + ", Message: " + response.errorMessage)
    })

    api.GetFollowing("me") // TODO: Fetch more than just the first page of results!
    .then(response => {
      if (response.ok) {
        this.setState({following: response.body.data, followingCount: response.body.pages.per_page * response.body.pages.total_pages})
        return
      }
      
      if (response.status == 401) {
        this.setState({authError: true})
      }
      console.log("Status: " + response.status + ", Code: " + response.errorCode + ", Message: " + response.errorMessage)
    })
  }

  handleSearchChange(e) {
    this.setState({search: e.target.value})
  }

  handleSearchSubmit(e) {
    e.preventDefault()

    if (!this.state.search) {
      this.GetDefaultResults()
      return
    }

    // api.Data(api.Search(this.state.search))
    // .then(items => {
    //   this.setState({results: items})
    // })
  }

  render() {
    return (
      <>
        <NavBar followerCount={this.state.followerCount} followingCount={this.state.followingCount} friendsCurrent={true} searchText="Find users by name" handleSearchSubmit={this.handleSearchSubmit} handleSearchChange={this.handleSearchChange}/>
        {
          function () {
            switch(window.location.pathname) {
              case "/followers":
                return <FollowerResults results={this.state.followers}/>
              case "/following":
                return <FollowingResults results={this.state.following}/>
              case "/find-users":
                return (<></>)
              default:
                console.error("unknown url: " + window.location.pathname)
                return (<></>)
            }
          }.bind(this)()
        }
        {/* <SearchResults results={this.state.results}/> */}
      </>
    )
  }
}

class FollowerResults extends React.Component {
  render() {
    return (
      <Container className={this.props.className} fluid>
          <Row>
            <Col xs={12} className="">
              <Row>
                <Col xs={6} className="mx-auto mt-4">
                  <h3>Your Followers</h3>
                  {this.props.results.map((result, index) => {
                    return <FollowerResult key={index} result={result}/>
                  })}
                </Col>
              </Row>
            </Col>
          </Row>
        </Container>
    )
  }
}

class FollowingResults extends React.Component {
  render() {
    return (
      <Container className={this.props.className} fluid>
          <Row>
            <Col xs={12} className="">
              <Row>
                <Col xs={6} className="mx-auto mt-4">
                  <h3>Users You Follow</h3>
                  {this.props.results.map((result, index) => {
                    return <FollowingResult key={index} result={result}/>
                  })}
                </Col>
              </Row>
            </Col>
          </Row>
        </Container>
    )
  }
}

class ResultsFooter extends React.Component {
  NoResults = "There aren't any results for that..."
  NotIt = "Can't find what you're looking for?"

  render() {
    return (
      <>
          <Col className="pt-5 pb-5 d-flex flex-column text-center align-items-center justify-content-center">
            <h3 className="no-results-header">{this.props.numResults == 0 ? this.NoResults : this.NotIt}</h3> 
            <p className="orange mt-3 no-results-text"><a href="/add-new">You can add it here!</a></p>
            <p className="paragraph-column no-results-text">Our community gets stronger each time you add a new show, movie, or book. It only takes a minute to do! We also keep track of your contributions so you can get credit for them!</p>
          </Col>
      </>
    )
  }
}

export class FollowerResult extends React.Component {
  render() {
    return (
      <>
        <Row className="pt-3 pb-3">
          <Col>
            <Row>
              <Col>
                <div style={{color: "white"}}>{this.props.result.follower_username}</div>
              </Col>
            </Row>
            <Row>
              <Col>
                <div style={{color: "#666"}}>{this.props.result.follower_full_name}</div>
              </Col>
            </Row>
          </Col>
          <Col className="text-right">
            <ActionButton action="remove" userId={this.props.result.follower_user_id}/>
          </Col>
        </Row>
        <Row>
          <Col className="search-result"></Col>
        </Row>
      </>
    )
  }
}

export class FollowingResult extends React.Component {
  render() {
    return (
      <>
        <Row className="pt-3 pb-3">
          <Col>
            <Row>
              <Col>
                <div style={{color: "white"}}>{this.props.result.target_username}</div>
              </Col>
            </Row>
            <Row>
              <Col>
                <div style={{color: "#666"}}>{this.props.result.target_full_name}</div>
              </Col>
            </Row>
          </Col>
          <Col className="text-right">
            <ActionButton action="unfollow" userId={this.props.result.target_user_id}/>
          </Col>
        </Row>
        <Row>
          <Col className="search-result"></Col>
        </Row>
      </>
    )
  }
}

const ActionButton = connect(class ActionButton extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      action: this.props.action,
      loading: false
    }

    this.handleUnfollow = this.handleUnfollow.bind(this)
    this.handleFollow = this.handleFollow.bind(this)
    this.handleRemove = this.handleRemove.bind(this)
    this.UnfollowLoadButton = this.UnfollowLoadButton.bind(this)
    this.FollowLoadButton = this.FollowLoadButton.bind(this)
    this.UnfollowButton = this.UnfollowButton.bind(this)
    this.FollowButton = this.FollowButton.bind(this)
    this.RemoveButton = this.RemoveButton.bind(this)
    this.RemovedButton = this.RemovedButton.bind(this)
  }

  handleUnfollow(e) {
    e.preventDefault()
    api.Data(api.UnfollowUser(e.target.id))
    .then(() => {
      this.setState({loading: false, action: "follow"})
    })
    this.setState({loading: true})
  }

  handleFollow(e) {
    e.preventDefault()
    api.Data(api.FollowUser(e.target.id))
    .then(() => {
      this.setState({loading: false, action: "unfollow"})
    })
    this.setState({loading: true})
  }

  handleRemove(e) {
    e.preventDefault()
    api.Data(api.RemoveUser(e.target.id))
    .then(() => {
      this.setState({loading: false, action: "removed"})
    })
    this.setState({loading: true})
  }

  UnfollowLoadButton = props => {
    return (
      <Form id={props.userId} onSubmit={null}>
          <Button type="submit" disabled variant="dark">
            Loading...
          </Button>
      </Form>
    )
  }

  FollowLoadButton = props => {
    return (
      <Form id={props.userId} onSubmit={null}>
          <Button type="submit" disabled className="orange-btn">
            Loading...
          </Button>
      </Form>
    )
  }

  UnfollowButton = props => {
    return (
      <Form id={props.userId} onSubmit={props.handleUnfollow}>
        <Button type="submit" variant="dark">
          Unfollow
        </Button>
      </Form>
    )
  }

  RemoveButton = props => {
    return (
      <Form id={props.userId} onSubmit={props.handleRemove}>
        <Button type="submit" variant="dark">
          Remove
        </Button>
      </Form>
    )
  }

  RemovedButton = props => {
    return (
      <Form id={props.userId} onSubmit={null}>
        <Button type="submit" disabled variant="dark">
          Removed
        </Button>
      </Form>
    )
  }

  FollowButton = props => {
    return (
      <Form id={props.userId}   onSubmit={props.handleFollow}>
        <Button type="submit" className="orange-btn">
          Follow
        </Button>
      </Form>
    )
  }

  render() {
    if ((this.state.action == "unfollow" || this.state.action == "remove") && this.state.loading) {
      return (<this.UnfollowLoadButton {...this.props}/>)
    }
    if (this.state.action == "follow" && this.state.loading) {
      return (<this.FollowLoadButton {...this.props}/>)
    }
    if (this.state.action == "unfollow") {
      return (<this.UnfollowButton {...this.props} handleUnfollow={this.handleUnfollow}/>)
    }
    if (this.state.action == "remove") {
      return (<this.RemoveButton {...this.props} handleRemove={this.handleRemove}/>)
    }
    if (this.state.action == "removed") {
      return (<this.RemovedButton {...this.props}/>)
    }
    if (this.state.action == "follow") {
      return (<this.FollowButton {...this.props} handleFollow={this.handleFollow}/>)
    }
    return (<></>)
  }
})