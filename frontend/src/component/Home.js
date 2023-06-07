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
      user: this.props.location?.state?.user,
    }

    this.GetUser = this.GetUser.bind(this)
    if (!this.state.user) {
      this.GetUser()
    }
  }

  GetUser() {
    api.GetUser()
    .then(res => {
      if (res.ok) {
        this.setState({
          user: res.body.data
        })
        return
      }
      this.setState({errorPage: true, errorCode: ERR_UNEXPECTED, errorMessage: res.errorMessage})
    })
    .catch(e => {
      console.log(e)
    })
  }

  render() {
    if (this.state.user && needsUsername(this.state.user)) {
      return <Redirect to={{pathname: '/new-user', state: {user: this.state.user}}}/>
    }
    return (
      <>
        <NavBar homeCurrent={true} />
      </>
    )
  }
}

function needsUsername(user) {
  return user.username.startsWith("_")
}