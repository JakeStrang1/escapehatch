import React from "react"
import Container from 'react-bootstrap/Container'
import InputGroup from 'react-bootstrap/InputGroup'
import FormControl from 'react-bootstrap/FormControl'
import Col from 'react-bootstrap/Col'
import Row from 'react-bootstrap/Row'
import Navbar from 'react-bootstrap/Navbar'
import Nav from 'react-bootstrap/Nav'
import Image from 'react-bootstrap/Image'
import banner from './../assets/banner.png'
import searchImage from './../assets/search_small.png'

export default class NavBar extends React.Component {  
  constructor(props){
    super(props);

    this.homeClass = this.props.homeCurrent ? "home-nav-current" : "home-nav" // Orange vs gray icon
    this.searchClass = this.props.searchCurrent ? "search-nav-current" : "search-nav"
    this.friendsClass = this.props.friendsCurrent ? "friends-nav-current" : "friends-nav"

    this.navSearchClass = this.props.searchCurrent ? "navbar-search" : ""
  }

  render() {
    return (
      <>
        <Navbar fixed="top" className={this.navSearchClass + " sticky-nav"} style={{backgroundColor: "black"}}>
          <Col>
            <Row>
              <Col xs={6} sm={6} md={5} lg={4} className="mt-auto mb-auto brand-column">
                <Navbar.Brand href="/">
                  <Image src={banner} fluid/>
                </Navbar.Brand>
              </Col>
              <Col xs={6} sm={6} md={5} lg={4}>
                <Row id="nav-link-row">
                    <Col className="bare-column"><Nav.Link href="/" className={this.homeClass}></Nav.Link></Col>
                    <Col className="bare-column"><Nav.Link href="/search" className={this.searchClass}></Nav.Link></Col>
                    <Col className="bare-column"><Nav.Link href="/followers" className={this.friendsClass}></Nav.Link></Col>
                </Row>
              </Col>
            </Row>
            {
              function () {
                if (this.props.searchCurrent) {
                  return <SearchBar/>
                }
              }.bind(this)()
            }
          </Col>
        </Navbar>
      </>
    )
  }
}

class SearchBar extends React.Component {
  render() {
    return (
      // <Container className={this.props.className} fluid>
        <Row>
          <Col xs={12} style={{backgroundColor:"black"}} className="pl-4">
            <Row>
              <Col sm={1} md={2} lg={3} xl={4}></Col>
              <Col>
                <div className="d-flex flex-row align-items-center justify-content-center">
                  <InputGroup className="mt-1 mb-1">
                    <InputGroup.Prepend>
                      <InputGroup.Text id="search" style={{backgroundColor: "white"}}>
                        <Image src={searchImage} fluid/>
                      </InputGroup.Text>
                    </InputGroup.Prepend>
                    <FormControl
                      placeholder="Find a show, movie, or book"
                      aria-label="Find a show, movie, or book"
                      aria-describedby="search"
                    />
                  </InputGroup>
                </div>
              </Col>
              <Col sm={1} md={2} lg={3} xl={4}></Col>
            </Row>
          </Col>
        </Row>
      // </Container>
    )
  }
}
