import React from "react"
import NavBar from "./NavBar"
import UserSummary from "./UserSummary"
import Shelves from "./Shelves"
import Container from 'react-bootstrap/Container'
import Col from 'react-bootstrap/Col'
import InputGroup from 'react-bootstrap/InputGroup'
import FormControl from 'react-bootstrap/FormControl'
import Image from 'react-bootstrap/Image'
import Button from 'react-bootstrap/Button'
import Badge from 'react-bootstrap/Badge'
import Row from 'react-bootstrap/Row'
import Form from 'react-bootstrap/Form'
import { connect } from './../reducers'

import { Redirect, Link } from 'react-router-dom'
import api, {
  ERR_UNEXPECTED,
} from "../api"

export default class Search extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      results: [],
      search: ""
    }

    this.GetDefaultResults = this.GetDefaultResults.bind(this)
    this.handleSearchSubmit = this.handleSearchSubmit.bind(this)
    this.handleSearchChange = this.handleSearchChange.bind(this)

    this.GetDefaultResults()
  }

  GetDefaultResults() {
    api.Data(api.GetItems())
    .then(items => {
      this.setState({results: items})
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

    api.Data(api.Search(this.state.search))
    .then(items => {
      this.setState({results: items})
    })
  }

  render() {
    return (
      <>
        <NavBar searchCurrent={true} searchText="Find a show, movie, or book" handleSearchSubmit={this.handleSearchSubmit} handleSearchChange={this.handleSearchChange}/>
        <SearchResults results={this.state.results}/>
      </>
    )
  }
}

class SearchResults extends React.Component {
  render() {
    return (
      <Container className={this.props.className} fluid>
          <Row>
            <Col xs={12} className="">
              <Row>
                <Col>
                  {this.props.results.map((result, index) => {
                    return <SearchResult key={index} result={result}/>
                  })}
                  <ResultsFooter numResults={this.props.results.length}/>
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

export class SearchResult extends React.Component {
  render() {
    return (
      <>
        <Row className="pt-3 pb-3">
          <Col xs={4} className="d-flex text-right align-items-center justify-content-end">
            <Image src={this.props.result.image_url} className="search-result-image" fluid rounded/>
          </Col>
          <Col xs={8} sm={7}>
            {
              function () {
                switch(this.props.result.media_type) {
                  case "book":
                    return <BookResult result={this.props.result} isPreview={this.props.isPreview}/>
                  case "movie":
                    return <MovieResult result={this.props.result} isPreview={this.props.isPreview}/>
                  case "tv_series":
                    return <TVSeriesResult result={this.props.result} isPreview={this.props.isPreview}/>
                  default:
                    console.error("unknown type: " + this.props.result.media_type)
                    return (<></>)
                }
              }.bind(this)()
            }
          </Col>
        </Row>
        <Row>
          <Col xs={1} md={2} lg={3}></Col>
          <Col className="search-result"></Col>
          <Col xs={1} md={2} lg={3}></Col>
        </Row>
      </>
    )
  }
}

class BookResult extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      userCount: this.props.result.user_count,
    }

    this.incrementUserCount = this.incrementUserCount.bind(this)
  }

  incrementUserCount() {
    this.setState({userCount: this.state.userCount+1})
  }

  render() {
    return (
      <>
        <div className="result-content">
          <div>
            <div className="result-media-type">Book</div>
            <div className="result-title">{this.props.result.description}</div>
            <div className="result-details">{this.props.result.published_year}&nbsp;&nbsp;•&nbsp;&nbsp;{this.props.result.author}</div>
          </div>
          <div className="result-stat-box"> 
            {
              function() {
                if (!this.props.isPreview) {
                  return (
                    <>
                      <div className="result-user-count">Added by {this.state.userCount} {peopleOrPerson(this.state.userCount)}</div>
                      <AddButton itemId={this.props.result.id} incrementUserCount={this.incrementUserCount}/>
                    </>
                  )
                }
              }.bind(this)()
            }
          </div>
        </div>
      </>
    )
  }
}

class MovieResult extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      userCount: this.props.result.user_count,
    }

    this.incrementUserCount = this.incrementUserCount.bind(this)
  }

  incrementUserCount() {
    this.setState({userCount: this.state.userCount+1})
  }

  render() {
    return (
      <>
        <div className="result-content">
          <div>
            <div className="result-media-type">Movie</div>
            <div className="result-title">{this.props.result.description}</div>
            <div className="result-details">{this.props.result.published_year}&nbsp;&nbsp;•&nbsp;&nbsp;{formatActors(this.props.result.lead_actors)}</div>
          </div>
          <div className="result-stat-box">
            {
              function() {
                if (!this.props.isPreview) {
                  return (
                    <>
                      <div className="result-user-count">Added by {this.state.userCount} {peopleOrPerson(this.state.userCount)}</div>
                      <AddButton itemId={this.props.result.id} incrementUserCount={this.incrementUserCount}/>
                    </>
                  )
                }
              }.bind(this)()
            }
          </div>
        </div>
      </>
    )
  }
}

class TVSeriesResult extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      userCount: this.props.result.user_count,
    }

    this.incrementUserCount = this.incrementUserCount.bind(this)
  }

  incrementUserCount() {
    this.setState({userCount: this.state.userCount+1})
  }

  render() {
    return (
      <>
        <div className="result-content">
          <div>
            <div className="result-media-type">TV Series</div>
            <div className="result-title">{this.props.result.title}</div>
            <div className="result-details">{this.props.result.tv_series_start_year} &ndash; {this.props.result.tv_series_end_year}&nbsp;&nbsp;•&nbsp;&nbsp;{formatActors(this.props.result.lead_actors)}</div>
          </div>
          <div className="result-stat-box">
            {
              function() {
                if (!this.props.isPreview) {
                  return (
                    <>
                      <div className="result-user-count">Added by {this.state.userCount} {peopleOrPerson(this.state.userCount)}</div>
                      <AddButton itemId={this.props.result.id} incrementUserCount={this.incrementUserCount}/>
                    </>
                  )
                }
              }.bind(this)()
            }
          </div>
        </div>
      </>
    )
  }
}

const AddButton = connect(class AddButton extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      loading: false
    }

    this.handleAddToShelf = this.handleAddToShelf.bind(this)
    this.LoadButton = this.LoadButton.bind(this)
    this.AddButton = this.AddButton.bind(this)
    this.AddedButton = this.AddedButton.bind(this)
  }

  handleAddToShelf(e) {
    e.preventDefault()
    api.Data(api.AddItem(e.target.id))
    .then(() => api.Data(api.GetUser()))
    .then(user => {
      this.setState({loading: false})
      this.props.updateUser(user)
      this.props.incrementUserCount()
    })
    this.setState({loading: true})
  }

  LoadButton = props => {
    return (
      <Form id={props.itemId} onSubmit={null}>
          <Button type="submit" disabled className="orange-btn">
            Loading...
          </Button>
      </Form>
    )
  }

  AddButton = props => {
    return (
      <Form id={props.itemId} onSubmit={props.handleAddToShelf}>
        <Button type="submit" className="orange-btn">
          Add to shelf
        </Button>
      </Form>
    )
  }

  AddedButton = props => {
    return (
      <Form id={props.itemId} onSubmit={null}>
        <Button type="submit" disabled className="orange-btn">
          Added to shelf
        </Button>
      </Form>
    )
  }

  render() {
    if (this.state.loading) {
      return (<this.LoadButton itemId={this.props.itemId}/>)
    }
    if (userHasItem(this.props.user, this.props.itemId)) {
      return (<this.AddedButton itemId={this.props.itemId}/>)
    }
    return (<this.AddButton itemId={this.props.itemId} handleAddToShelf={this.handleAddToShelf}/>)
  }
})

function formatActors(actors) {
  return actors.join(", ")
}

function peopleOrPerson(num) {
  if (num == 1) {
    return "person"
  }
  return "people"
}

function userHasItem(user, itemId) {
  if (!user.shelves) {
    return false
  }

  for (let i = 0; i < user.shelves.length; i++) {
    if (!user.shelves[i]?.items) {
      continue
    }
    for (let j = 0; j < user.shelves[i].items.length; j++) {
      if (user.shelves[i].items[j].item_id == itemId) {
        return true
      }
    }
  }

  return false
}