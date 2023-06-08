import React from 'react'
import ReactDOM from 'react-dom'
import {
  BrowserRouter as Router,
  Route,
  Switch,
  Redirect
} from "react-router-dom";
import NoAuthRoute from "./component/NoAuthRoute"
import AuthRoute from "./component/AuthRoute"
import NewUser from "./component/NewUser"
import ResetAuthRoute from "./component/ResetAuthRoute"
import SignIn from "./component/SignIn"
import SignUp from "./component/SignUp"
import NotYou from "./component/NotYou"
import Verify from "./component/Verify"
import VerifyLink from "./component/VerifyLink"
import Home from "./component/Home"
import ErrorPage from "./component/ErrorPage"
import AuthContext from "./authContext"
// import Cookies from 'js-cookie';
import api from './api'
import 'bootstrap/dist/css/bootstrap.min.css'
import './assets/stylesheet.css'

const auth = {
  attempted: false,
  user: null,
  User: async function() {
    this.attempted = true
    return api.GetUser()
    .then(response => {
      if (response.ok) {
        this.user = response.body.data
        return
      }
      console.log("Status: " + response.status + ", Code: " + response.errorCode + ", Message: " + response.errorMessage)
    })
    .catch(e => {
      console.log(e)
    })
  }
}

ReactDOM.render(
    <Router>
      <AuthContext.Provider value={auth}>
        <Switch>
          <NoAuthRoute path="/sign-up" component={SignUp}/>
          <NoAuthRoute path="/sign-in" component={SignIn}/>
          <ResetAuthRoute path="/verify" component={Verify}/>
          <ResetAuthRoute path="/verify-link" component={VerifyLink}/>
          <Route path="/not-you" component={NotYou}/>
          <Route path="/oh-no" component={ErrorPage}/>
          <AuthRoute path="/new-user" component={NewUser}/>
          <AuthRoute exact path="/" redirect="/sign-up" component={Home}/>
          <Redirect to={{pathname: "/oh-no", state: { errorCode: "not_found"}}}/>
        </Switch>
      </AuthContext.Provider>
    </Router>,
  document.getElementById('root')
)
