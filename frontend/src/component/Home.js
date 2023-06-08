import React from "react"
import NavBar from "./NavBar"
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
      </>
    )
  }
}
