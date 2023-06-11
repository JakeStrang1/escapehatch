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
                </Col>
              </Row>
            </Col>
          </Row>
        </Container>
    )
  }
}

class SearchResult extends React.Component {
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
                    return <BookResult result={this.props.result}/>
                  case "movie":
                    return <MovieResult result={this.props.result}/>
                  case "tv_series":
                    return <TVSeriesResult result={this.props.result}/>
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
            <div className="result-user-count">Added by {this.props.result.user_count} people</div>
            <AddButton/>
          </div>
        </div>
      </>
    )
  }
}

class MovieResult extends React.Component {
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
            <div className="result-user-count">Added by {this.props.result.user_count} people</div>
            <AddButton/>
          </div>
        </div>
      </>
    )
  }
}

class TVSeriesResult extends React.Component {
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
            <div className="result-user-count">Added by {this.props.result.user_count} people</div>
            <AddButton/>
          </div>
        </div>
      </>
    )
  }
}

class AddButton extends React.Component {
  render() {
    return (
      <>
        <Button className="orange-btn">Add to shelf</Button>
      </>
    )
  }
}

function formatActors(actors) {
  return actors.join(", ")
}