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
      showingSearchResult: false,
      authError: false,
      followers: [],
      following: [],
      newUsers: [],
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
        this.setState({
          followers: response.body.data,
          followerCount: response.body.pages.per_page * response.body.pages.total_pages,
          showingSearchResult: false
        })
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
        this.setState({
          following: response.body.data, 
          followingCount: response.body.pages.per_page * response.body.pages.total_pages,
          showingSearchResult: false
        })
        return
      }
      
      if (response.status == 401) {
        this.setState({authError: true})
      }
      console.log("Status: " + response.status + ", Code: " + response.errorCode + ", Message: " + response.errorMessage)
    })

    api.GetUsers() // TODO: Fetch more than just the first page of results!
    .then(response => {
      if (response.ok) {
        let filteredResults = response.body.data.filter(result => !result.followed_by_you && !result.self) // Only include users not already followed
        this.setState({
          newUsers: filteredResults,
          showingSearchResult: false
        })
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

    let req

    switch(window.location.pathname) {
      case "/followers":
        api.Data(api.GetFollowers("me", this.state.search)).then(results => this.setState({followers: results, showingSearchResult: true}))
      case "/following":
        api.Data(api.GetFollowing("me", this.state.search)).then(results => this.setState({following: results, showingSearchResult: true}))
      case "/find-users":
        api.Data(api.GetUsers(this.state.search)).then(results => this.setState({newUsers: results, showingSearchResult: true}))
      default:
        console.error("unknown url: " + window.location.pathname)
        return
    }
  }

  render() {
    if (this.state.authError) {
      return <Redirect to="/sign-out"/>
    }
    return (
      <>
        <NavBar followerCount={this.state.followerCount} followingCount={this.state.followingCount} friendsCurrent={true} searchText="Find users by name" handleSearchSubmit={this.handleSearchSubmit} handleSearchChange={this.handleSearchChange}/>
        {
          function () {
            switch(window.location.pathname) {
              case "/followers":
                return <FollowerResults results={this.state.followers} showingSearchResult={this.state.showingSearchResult}/>
              case "/following":
                return <FollowingResults results={this.state.following} showingSearchResult={this.state.showingSearchResult}/>
              case "/find-users":
                return <FindUsersResults results={this.state.newUsers} showingSearchResult={this.state.showingSearchResult}/>
              default:
                console.error("unknown url: " + window.location.pathname)
                return (<></>)
            }
          }.bind(this)()
        }
      </>
    )
  }
}

class FindUsersResults extends React.Component {
  render() {
    return (
      <Container className={this.props.className} fluid>
          <Row>
            <Col xs={12} className="">
              <Row>
                <Col xs={12} sm={10} md={8} lg={6} className="mx-auto mt-4">
                  <h3>
                    {this.props.showingSearchResult ? "Search results..." : "People You May Know"}
                  </h3>
                  {this.props.results.map(result => {
                    return <UserResult key={result.id} result={result}/>
                  })}
                </Col>
              </Row>
            </Col>
          </Row>
        </Container>
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
                <Col xs={12} sm={10} md={8} lg={6} className="mx-auto mt-4">
                  <h3>
                    {this.props.showingSearchResult ? "Search results in followers..." : "All Followers"}
                  </h3>
                  {this.props.results.map(result => {
                    return <FollowerResult key={result.follower_user_id} result={result}/>
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
                <Col xs={12} sm={10} md={8} lg={6} className="mx-auto mt-4">
                  <h3>
                    {this.props.showingSearchResult ? "Search results in following..." : "All Following"}
                  </h3>
                  {this.props.results.map(result => {
                    return <FollowingResult key={result.target_user_id} result={result}/>
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

export class UserResult extends React.Component {
  render() {
    return (
      <>
        <Row className="pt-3 pb-3">
          <Col as="a" href={"/users/" + this.props.result.id}>
            <Row>
              <Col>
                <div style={{color: "white"}}>{this.props.result.username}</div>
              </Col>
            </Row>
            <Row>
              <Col>
                <div style={{color: "#666"}}>{this.props.result.full_name}</div>
              </Col>
            </Row>
          </Col>
          <Col xs="auto">
            {
              function () {
                if (this.props.result.self) {
                  return (<ActionButton action="self" userId={this.props.result.id}/>)
                } else if (this.props.result.followed_by_you) {
                  return (<ActionButton action="unfollow" userId={this.props.result.id}/>)
                } else if (this.props.result.follows_you) {
                  return (<ActionButton action="follow back" userId={this.props.result.id}/>)
                } else {
                  return (<ActionButton action="follow" userId={this.props.result.id}/>)
                }
              }.bind(this)()
            }
          </Col>
        </Row>
        <Row>
          <Col className="search-result"></Col>
        </Row>
      </>
    )
  }
}

export class FollowerResult extends React.Component {
  render() {
    return (
      <>
        <Row className="pt-3 pb-3">
          <Col as="a" href={"/users/" + this.props.result.follower_user_id}>
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
          <Col xs="auto">
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
        <Col as="a" href={"/users/" + this.props.result.target_user_id}>
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
          <Col xs="auto">
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
    this.FollowBackButton = this.FollowBackButton.bind(this)
    this.SelfButton = this.SelfButton.bind(this)
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

  FollowBackButton = props => {
    return (
      <Form id={props.userId}   onSubmit={props.handleFollow}>
        <Button type="submit" className="orange-btn">
          Follow back
        </Button>
      </Form>
    )
  }

  SelfButton = props => {
    return (
      <Form id={props.userId} onSubmit={null}>
        <Button type="submit" disabled variant="dark">
          This is you
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
    if (this.state.action == "follow back") {
      return (<this.FollowBackButton {...this.props} handleFollow={this.handleFollow}/>)
    }
    if (this.state.action == "self") {
      return (<this.SelfButton {...this.props}/>)
    }
    return (<></>)
  }
})