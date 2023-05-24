import React from 'react'
import ReactDOM from 'react-dom'
import {
  BrowserRouter as Router,
  Route,
  Switch
} from "react-router-dom";
import MainLayout from "./component/MainLayout"
import NoAuthRoute from "./component/NoAuthRoute"
import ResetAuthRoute from "./component/ResetAuthRoute"
import SignIn from "./component/SignIn"
import SignUp from "./component/SignUp"
import NotYou from "./component/NotYou"
import Verify from "./component/Verify"
import ErrorPage from "./component/ErrorPage"
import AuthContext from "./authContext"
// import Cookies from 'js-cookie';
import api from './api'
import 'bootstrap/dist/css/bootstrap.min.css'
import './stylesheet.css'

const auth = {
  attempted: false,
  user: null,
  User: async function() {
    this.attempted = true
    return api.GetUser()
    .then(response => {
      if (response.ok) {
        this.user = response.body
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
          <Route path="/not-you" component={NotYou}/>
          <Route path="/oh-no" component={ErrorPage}/>
          <Route path="/" component={MainLayout}/>
        </Switch>
      </AuthContext.Provider>
    </Router>,
  document.getElementById('root')
)
