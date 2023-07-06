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
import { Redirect, Link } from 'react-router-dom'
import api, {
  ERR_UNEXPECTED,
} from "../api"

export default class Followers extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      authError: false,
      results: [],
      search: ""
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
        this.setState({results: response.body.data})
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

    api.Data(api.Search(this.state.search))
    .then(items => {
      this.setState({results: items})
    })
  }

  render() {
    return (
      <>
        <NavBar friendsCurrent={true} searchText="Find users by name" handleSearchSubmit={this.handleSearchSubmit} handleSearchChange={this.handleSearchChange}/>
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
          <Col xs={8} sm={7}>
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
