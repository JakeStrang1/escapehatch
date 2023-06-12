import React from "react"
import { Redirect } from 'react-router-dom'
import NavBar from "./NavBar"
import Container from 'react-bootstrap/Container'
import Col from 'react-bootstrap/Col'
import Dropdown from 'react-bootstrap/Dropdown'
import Button from 'react-bootstrap/Button'
import Badge from 'react-bootstrap/Badge'
import Row from 'react-bootstrap/Row'
import api, {
  ERR_UNEXPECTED,
} from "../api"


class UserSummary extends React.Component {
  constructor(props){
    super(props);
  }

  render() {
    return (
      <>
        {/* Display on XS only */}
        <UserSummaryXS className="d-sm-none" user={this.props.user}/>

        {/* Display on SM and MD only */}
        <UserSummaryMD className="d-none d-sm-block d-lg-none" user={this.props.user}/>

        {/* Display on LG only */}
        <UserSummaryLG className="d-none d-lg-block d-xl-none" user={this.props.user}/>

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
            <Col xs={12} style={{backgroundColor:"#222"}} className="pt-3 pl-4 pb-3">
              <Row>
                <Col>
                  <div className="d-flex flex-row justify-content-between">
                    <UserLabelXS user={this.props.user}/>
                    <RightSummaryXS user={this.props.user}/>
                  </div>
                </Col>
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
                <Col>
                  <div className="d-flex flex-row align-items-center justify-content-between">
                    <UserLabel user={this.props.user}/>
                    <RightSummary user={this.props.user}/>
                  </div>
                </Col>
              </Row>
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}

class UserSummaryLG extends React.Component {
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
                <Col xs={1}></Col>
                <Col>
                  <div className="d-flex flex-row align-items-center justify-content-between">
                    <UserLabel user={this.props.user}/>
                    <RightSummary user={this.props.user}/>
                  </div>
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
                <Col xs={2}></Col>
                <Col>
                  <div className="d-flex flex-row align-items-center justify-content-between">
                    <UserLabel user={this.props.user}/>
                    <RightSummary user={this.props.user}/>
                  </div>
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

class UserLabelXS extends React.Component {
  render() {
    return (
      <>
        <Col className="d-flex flex-column">
          <Row>
            <Col>
              <span id="username-xl">{this.props.user.username}</span>&nbsp;&nbsp;&nbsp;
              <span id="userNumber-xl">#{this.props.user.number}</span>
            </Col>
          </Row>
          <Row className="flex-grow-1">
            <Col className="d-flex flex-column justify-content-center">
              <span id="fullName-xs">{this.props.user.full_name}</span>
            </Col>
          </Row>
        </Col>
      </>
    )
  }
}

class UserLabel extends React.Component {
  render() {
    return (
      <>
      <div className="d-inline">
        <span id="username-xl">{this.props.user.username}</span>&emsp;
        <span id="userNumber-xl">#{this.props.user.number}&emsp;•&emsp;</span>
        <span id="fullName-xl">{nbsp(this.props.user.full_name)}</span>
        
      </div>
    </>
    )
  }
}

class FollowerStatXS extends React.Component {
  render() {
    return (
      <div className="text-center">
        <Row>
          <Col>
            <span className="followerCount-xs">{this.props.count}</span>
          </Col>
        </Row>
        <Row>
          <Col>
          <span className="followerLabel-xs">{this.props.label}</span>
          </Col>
        </Row>
      </div>
    )
  }
}

class FollowerStat extends React.Component {
  render() {
    return (
      <div className="text-center">
        <Row>
          <Col>
            <span className="followerCount">{this.props.count}</span>
          </Col>
        </Row>
        <Row>
          <Col>
          <span className="followerLabel">{this.props.label}</span>
          </Col>
        </Row>
      </div>
    )
  }
}

class UserDropDown extends React.Component {
  render() {
    return (
      <Dropdown>
        <Dropdown.Toggle className="userDropDown" variant="secondary">
        </Dropdown.Toggle>
        <Dropdown.Menu>
          <Dropdown.Item href="/sign-out">Log out</Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    )
  }
}

class RightSummaryXS extends React.Component {
  render() {
    return (
      <>
        <Col className="rightSummary-xs">
          <Row>
            <Col className="text-right">
              <UserDropDown/>
            </Col>
          </Row>
          <Row className="justify-content-end mt-2">
            <Col xs={12} className="d-flex flex-row justify-content-between">
              <FollowerStatXS count={this.props.user.follower_count} label="Followers"/>
              <FollowerStatXS count={this.props.user.following_count} label="Following"/>
            </Col>
          </Row>
        </Col>
      </>
    )
  }
}

class RightSummary extends React.Component {
  render() {
    return (
      // <Row>
        <div className="d-flex flex-row rightSummary justify-content-between">
          <FollowerStat count={this.props.user.follower_count} label="Followers"/>
          <FollowerStat count={this.props.user.following_count} label="Following"/>
          <UserDropDown/>
        </div>
      // </Row>
    )
  }
}

function nbsp(str) {
  return str.replace(/ /g, "\u00a0")
}

export default UserSummary;