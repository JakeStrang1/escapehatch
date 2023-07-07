import React from "react"
import Container from 'react-bootstrap/Container'
import InputGroup from 'react-bootstrap/InputGroup'
import FormControl from 'react-bootstrap/FormControl'
import Tabs from 'react-bootstrap/Tabs'
import Tab from 'react-bootstrap/Tab'
import Col from 'react-bootstrap/Col'
import Row from 'react-bootstrap/Row'
import Navbar from 'react-bootstrap/Navbar'
import Nav from 'react-bootstrap/Nav'
import Image from 'react-bootstrap/Image'
import banner from './../assets/banner.png'
import searchImage from './../assets/search_small.png'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import FormGroup from "react-bootstrap/esm/FormGroup"

export default class NavBar extends React.Component {  
  constructor(props){
    super(props);

    this.homeClass = this.props.homeCurrent ? "home-nav-current" : "home-nav" // Orange vs gray icon
    this.searchClass = this.props.searchCurrent ? "search-nav-current" : "search-nav"
    this.friendsClass = this.props.friendsCurrent ? "friends-nav-current" : "friends-nav"
    this.searchBar = this.props.searchCurrent
    this.friendsBar = this.props.friendsCurrent

    this.navSearchClass = ""
    if (this.searchBar) {
      this.navSearchClass = "navbar-search"
    } else if (this.friendsBar) {
      this.navSearchClass = "navbar-friends"
    }
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
                if (this.searchBar) {
                  return <SearchBar {... this.props}/>
                } else if (this.friendsBar) {
                  return <FriendsBar {... this.props}/>
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
        <Row style={{backgroundColor:"black"}} >
          <Col xs={12} style={{backgroundColor:"black"}}  className="pl-4">
            <Row>
              <Col sm={1} md={2} lg={3} xl={4}></Col>
              <Col>
                <div className="d-flex flex-row align-items-center justify-content-center">
                  <Form className="search-form" onSubmit={this.props.handleSearchSubmit} onChange={this.props.handleSearchChange}>
                    <Form.Group className="mb-0 mt-2">
                      <Form.Label className="sr-only">
                        Search
                      </Form.Label>
                      <InputGroup>
                        <InputGroup.Prepend>
                          <InputGroup.Text id="search" style={{backgroundColor: "white"}}>
                            <Image src={searchImage} fluid/>
                          </InputGroup.Text>
                        </InputGroup.Prepend>
                        <FormControl
                          placeholder={this.props.searchText}
                          aria-label={this.props.searchText}
                          aria-describedby="search"
                        />
                      </InputGroup>
                    </Form.Group>
                    <FormGroup controlId="searchSubmit" className="mb-0">
                      <Button variant="primary" type="submit" className="sr-only">
                        Submit
                      </Button>
                    </FormGroup>
                  </Form>
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

class FriendsBar extends React.Component {
  render() {
    return (
      <Row className="full-width">
        <Col xs={12}>
          <Row>
            <Col style={{backgroundColor:"#222"}} className="pt-2">
              <Row>
                <Col xs={12} md={10} lg={8} xl={6} className="mx-auto">
                  <Nav justify className="justify-content-center friend-tabs" variant="tabs" defaultActiveKey={window.location.pathname}>
                    <Nav.Item>
                      <Nav.Link href="/followers" className="friend-tab">{this.props.followerCount + " Followers"}</Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                      <Nav.Link href="/following" className="friend-tab">{this.props.followingCount + " Following"}</Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                      <Nav.Link href="/find-users" className="friend-tab">Discover</Nav.Link>
                    </Nav.Item>
                  </Nav>
                </Col>
              </Row>
            </Col>
          </Row>
          <SearchBar searchText={this.props.searchText} handleSearchSubmit={this.props.handleSearchSubmit} handleSearchChange={this.props.handleSearchChange}/>
        </Col>
      </Row>
    )
  }
}
