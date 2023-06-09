import React from "react"
import NavBar from "./NavBar"
import UserSummary from "./UserSummary"
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

export default class Home extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      user: this.props.user,
    }

  }

  render() {
    return (
      <>
        <NavBar homeCurrent={true} />
        <UserSummary user={this.props.user}/>
      </>
    )
  }
}
