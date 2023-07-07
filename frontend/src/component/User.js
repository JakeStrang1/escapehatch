import React from "react"
import {
  Redirect
} from "react-router-dom";
import NavBar from "./NavBar"
import UserSummary from "./UserSummary"
import Shelves from "./Shelves"
import { connect } from "react-redux";
import api, {
  ERR_UNEXPECTED,
} from "../api"

class User extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      user: null,
      error: false
    }

    api.GetUser(props.match.params.id)
    .then(response => {
        if (response.ok) {
          this.setState({user: response.body.data})
          return
        }
        console.log("Status: " + response.status + ", Code: " + response.errorCode + ", Message: " + response.errorMessage)
        this.setState({error: true})
    })
    .catch(e => {
        console.log(e)
    })
  }

  render() {
    if (this.state.error) {
      return (<Redirect to={{pathname: "/oh-no", state: { errorCode: "not_found"}}}/>)
    }
    return (
      <>
        <NavBar/>
        {
          this.state.user && (
            <>
              <UserSummary user={this.state.user}/>
              <Shelves user={this.state.user}/>
            </>
          )
        }
      </>
    )
  }
}

function mapStateToProps(state) {
  const user = state.user.value
  return {
    user
  };
}

export default connect(mapStateToProps)(User);
