import React from "react"
import Col from 'react-bootstrap/Col'
import Row from 'react-bootstrap/Row'
import Navbar from 'react-bootstrap/Navbar'
import Nav from 'react-bootstrap/Nav'
import Image from 'react-bootstrap/Image'
import banner from './../assets/banner.png'

export default class NavBar extends React.Component {  
  constructor(props){
    super(props);
    this.homeClass = this.props.homeCurrent ? "home-nav-current" : "home-nav" // Orange vs gray icon
    this.searchClass = this.props.searchCurrent ? "search-nav-current" : "search-nav"
    this.friendsClass = this.props.friendsCurrent ? "friends-nav-current" : "friends-nav"
  }

  render() {
    return (
      <>
        <Navbar fixed="top" className="sticky-nav">
            <Col xs={6} sm={6} md={5} lg={4} align-self-center className="mt-auto mb-auto brand-column">
              <Navbar.Brand href="#home">
                <Image src={banner} fluid/>
              </Navbar.Brand>
            </Col>
            <Col xs={6} sm={6} md={5} lg={4}>
              <Row id="nav-link-row">
                  <Col className="bare-column"><Nav.Link href="#home" className={this.homeClass}></Nav.Link></Col>
                  <Col className="bare-column"><Nav.Link href="#link" className={this.searchClass}></Nav.Link></Col>
                  <Col className="bare-column"><Nav.Link href="#link" className={this.friendsClass}></Nav.Link></Col>
              </Row>
            </Col>
        </Navbar>
      </>
    )
  }
}




