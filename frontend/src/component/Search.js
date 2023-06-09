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

export default class Search extends React.Component {
  constructor(props){
    super(props);
  }

  render() {
    return (
      <>
        <NavBar searchCurrent={true} />
      </>
    )
  }
}
