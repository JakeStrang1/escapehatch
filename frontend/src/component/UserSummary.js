import React from "react"
import NavBar from "./NavBar"
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

export default class UserSummary extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      user: this.props.user,
    }

  }

  render() {
    return (
      <>
        {/* Display on XS only */}
        <UserSummaryXS className="d-sm-none" user={this.props.user}/>

        {/* Display on SM only */}
        <UserSummarySM className="d-none d-sm-block d-md-none" user={this.props.user}/>

        {/* Display on MD and LG only */}
        <UserSummaryMD className="d-none d-md-block d-xl-none" user={this.props.user}/>

        {/* Display on XL only */}
        <UserSummaryXL className="d-none d-xl-block" user={this.props.user}/>
      </>
    )
  }
}

class UserSummaryXS extends React.Component {
    constructor(props){
      super(props);
  
      this.state = {
        user: this.props.user,
      }
  
    }
  
    render() {
      return (
        <>
          <Container className={this.props.className} fluid>
            <Row>
              <Col xs={12} style={{backgroundColor:"#222"}} className="pt-3 pl-4 pb-2">
                <Row>
                  <Col xs={9}>
                    <span id="username">{this.props.user.username}<span id="userNumber">#{this.props.user.number}</span></span>
                  </Col>
                  <Col xs={3} className="text-right">
                    <Dropdown>
                      <Dropdown.Toggle className="userDropDown" variant="secondary">
                      </Dropdown.Toggle>
                      <Dropdown.Menu>
                        <Dropdown.Item href="#/action-1">Log out</Dropdown.Item>
                      </Dropdown.Menu>
                    </Dropdown>
                  </Col>
                </Row>
                <Row>
                  <Col xs={6}>
                    <Row>
                      <Col>
                        <span id="fullName">{this.props.user.full_name}</span>
                      </Col>
                    </Row>
                    <Row>
                      <Col>
                        <Badge pill type="submit" className="userStatusBtn selfStatusBtn">You</Badge>
                      </Col>
                    </Row>
                  </Col>
                  <Col xs={3} className="followerColumn">
                    <Row>
                      <Col>
                        <span className="followerCount">{this.props.user.follower_count}</span>
                      </Col>
                    </Row>
                    <Row>
                      <Col>
                      <span className="followerLabel">Followers</span>
                      </Col>
                    </Row>
                  </Col>
                  <Col xs={3} className="followerColumn">
                    <Row>
                      <Col>
                      <span className="followerCount">{this.props.user.following_count}</span>
                      </Col>
                    </Row>
                    <Row>
                      <Col>
                      <span className="followerLabel">Following</span>
                      </Col>
                    </Row>
                  </Col>
                </Row>
              </Col>
            </Row>
          </Container>
        </>
      )
    }
  }

class UserSummarySM extends React.Component {
    constructor(props){
      super(props);
  
      this.state = {
        user: this.props.user,
      }
  
    }
  
    render() {
      return (
        <>
          <Container className={this.props.className} fluid>
            <Row>
              <Col xs={12} style={{backgroundColor:"#222"}} className="pt-3 pl-4 pb-2">
                <Row>
                  <Col xs={1}></Col>
                  <Col xs={8}>
                    <span id="username">{this.props.user.full_name}&emsp;<span style={{color: "#555"}}>&mdash;</span>&emsp;{this.props.user.username}<span id="userNumber">#{this.props.user.number}</span></span>
                  </Col>
                  <Col xs={2} className="text-right">
                    <Dropdown>
                      <Dropdown.Toggle className="userDropDown" variant="secondary">
                      </Dropdown.Toggle>
                      <Dropdown.Menu>
                        <Dropdown.Item href="#/action-1">Log out</Dropdown.Item>
                      </Dropdown.Menu>
                    </Dropdown>
                  </Col>
                  <Col xs={1}></Col>
                </Row>
                <Row className="mt-3">
                  <Col xs={1}></Col>
                  <Col>
                    <Row>
                      <Col>
                        <Badge pill type="submit" className="userStatusBtn selfStatusBtn">You</Badge>
                      </Col>
                    </Row>
                  </Col>
                  <Col className="followerColumn">
                    <Row>
                      <Col>
                        <span className="followerCount">{this.props.user.follower_count}</span>
                      </Col>
                    </Row>
                    <Row>
                      <Col>
                      <span className="followerLabel">Followers</span>
                      </Col>
                    </Row>
                  </Col>
                  <Col className="followerColumn">
                    <Row>
                      <Col>
                      <span className="followerCount">{this.props.user.following_count}</span>
                      </Col>
                    </Row>
                    <Row>
                      <Col>
                      <span className="followerLabel">Following</span>
                      </Col>
                    </Row>
                  </Col>
                  <Col xs={1}></Col>
                </Row>
              </Col>
            </Row>
          </Container>
        </>
      )
    }
}

class UserSummaryMD extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      user: this.props.user,
    }

  }

  render() {
    return (
      <>
        <Container className={this.props.className} fluid>
          <Row>
            <Col xs={12} style={{backgroundColor:"#222"}} className="pt-3 pl-4 pb-3">
              <Row>
                <Col xs={2}></Col>
                <Col>
                  <Row>
                    <Col xs={11}>
                      <span id="username">{this.props.user.full_name}&emsp;<span style={{color: "#555"}}>&mdash;</span>&emsp;{this.props.user.username}<span id="userNumber">#{this.props.user.number}</span></span>
                    </Col>
                    <Col xs={1} className="text-right">
                      <Dropdown>
                        <Dropdown.Toggle className="userDropDown" variant="secondary">
                        </Dropdown.Toggle>
                        <Dropdown.Menu>
                          <Dropdown.Item href="#/action-1">Log out</Dropdown.Item>
                        </Dropdown.Menu>
                      </Dropdown>
                    </Col>
                  </Row>
                  <Row className="mt-3">
                    <Col>
                      <Row>
                        <Col>
                          <Badge pill type="submit" className="userStatusBtn selfStatusBtn">You</Badge>
                        </Col>
                      </Row>
                    </Col>
                    <Col className="followerColumn">
                      <Row>
                        <Col>
                          <span className="followerCount">{this.props.user.follower_count}</span>
                        </Col>
                      </Row>
                      <Row>
                        <Col>
                        <span className="followerLabel">Followers</span>
                        </Col>
                      </Row>
                    </Col>
                    <Col style={{justifyContent: "right"}} className="d-flex">
                      <Row>
                        <Col className="followerColumn">
                          <Row>
                            <Col>
                            <span className="followerCount">{this.props.user.following_count}</span>
                            </Col>
                          </Row>
                          <Row>
                            <Col>
                            <span className="followerLabel">Following</span>
                            </Col>
                          </Row>
                        </Col>
                      </Row>
                    </Col>
                  </Row>
                </Col>
                <Col xs={2}></Col>
              </Row>
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}

class UserSummaryXL extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      user: this.props.user,
    }

  }

  render() {
    return (
      <>
        <Container className={this.props.className} fluid>
          <Row>
            <Col xs={12} style={{backgroundColor:"#222"}} className="pt-3 pl-4 pb-3">
              <Row>
                <Col xs={3}></Col>
                <Col>
                  <Row>
                    <Col xs={11}>
                      <span id="username">{this.props.user.full_name}&emsp;<span style={{color: "#555"}}>&mdash;</span>&emsp;{this.props.user.username}<span id="userNumber">#{this.props.user.number}</span></span>
                    </Col>
                    <Col xs={1} className="text-right">
                      <Dropdown>
                        <Dropdown.Toggle className="userDropDown" variant="secondary">
                        </Dropdown.Toggle>
                        <Dropdown.Menu>
                          <Dropdown.Item href="#/action-1">Log out</Dropdown.Item>
                        </Dropdown.Menu>
                      </Dropdown>
                    </Col>
                  </Row>
                  <Row className="mt-3">
                    <Col>
                      <Row>
                        <Col>
                          <Badge pill type="submit" className="userStatusBtn selfStatusBtn">You</Badge>
                        </Col>
                      </Row>
                    </Col>
                    <Col className="followerColumn">
                      <Row>
                        <Col>
                          <span className="followerCount">{this.props.user.follower_count}</span>
                        </Col>
                      </Row>
                      <Row>
                        <Col>
                        <span className="followerLabel">Followers</span>
                        </Col>
                      </Row>
                    </Col>
                    <Col style={{justifyContent: "right"}} className="d-flex">
                      <Row>
                        <Col className="followerColumn">
                          <Row>
                            <Col>
                            <span className="followerCount">{this.props.user.following_count}</span>
                            </Col>
                          </Row>
                          <Row>
                            <Col>
                            <span className="followerLabel">Following</span>
                            </Col>
                          </Row>
                        </Col>
                      </Row>
                    </Col>
                  </Row>
                </Col>
                <Col xs={3}></Col>
              </Row>
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}
