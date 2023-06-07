import React from "react"
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Container from 'react-bootstrap/Container'

export default class Loading extends React.Component {
  render() {
    return (
      <Container fluid>
        <Row className="landing-row head-space foot-space">
          <Col xs={6} className="">
            <Row className="align-items-center d-flex flex-row-reverse">
                <p className="typewriter-test">loading loading loading loading loading loading</p>
            </Row>
          </Col>
          <Col xs={6}></Col>
        </Row>
      </Container>
    )
  }
}