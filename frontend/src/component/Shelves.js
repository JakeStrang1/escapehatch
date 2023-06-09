import React from "react"
import Col from 'react-bootstrap/Col'
import Row from 'react-bootstrap/Row'
import Container from 'react-bootstrap/Container'
import Navbar from 'react-bootstrap/Navbar'
import Nav from 'react-bootstrap/Nav'
import Image from 'react-bootstrap/Image'
import emptyShelfImage from './../assets/empty_shelf.png'
import searchImage from './../assets/search.png'

import banner from './../assets/banner.png'

export default class Shelves extends React.Component {  
  constructor(props){
    super(props);
  }

  render() {
    if (!this.props.user.shelves || this.props.user.shelves.length == 0) {
        return (
            <Empty/>
        )
    }

    return (
      <>

      </>
    )
  }
}

class Empty extends React.Component {
  render() {
    return (
      <>
        <Col className="align-items-center mt-5">
          <Row className="d-flex">
            <Col xs={1}></Col>
            <Col>
              <Row className="d-flex align-items-center text-center pt-2">
                <Col xs={2}></Col>
                <Col>
                  <Image src={emptyShelfImage} fluid/>
                </Col>
                <Col xs={2}></Col>
              </Row>
              <Row className="d-flex align-items-center text-center pt-4">
                <Col>
                  <h4 className="orange">Your shelf is empty!</h4>
                  <p>Your shelf is where you showcase the shows, movies, and books that you love to escape into.</p>
                  <div className="d-inline-block mt-3">
                    <a href="/search">
                      <div className="outline-link pt-3 pr-3 pl-3 pb-2">
                        <Image src={searchImage}/>
                        <span style={{color: "white"}} className="d-block mt-2">Add to shelf</span>
                      </div>
                    </a>
                  </div>
                </Col>
              </Row>
            </Col>
            <Col xs={1}></Col>
          </Row>
        </Col>
      </>
    )
  }
}




