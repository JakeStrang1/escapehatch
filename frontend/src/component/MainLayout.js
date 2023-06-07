import React from "react"
import {
  Route,
  Switch,
  Redirect
} from "react-router-dom";
import Home from "./Home"
import NewUser from "./NewUser"
import AuthRoute from "./AuthRoute"
import ResetAuthRoute from "./ResetAuthRoute"
import VerifyLink from "./VerifyLink"
import Verify from "./Verify"
import NotYou from "./NotYou"

export default class MainLayout extends React.Component {
  render() {
    return (
      <div>
        <Switch>
          <AuthRoute exact path="/" redirect="/sign-up" component={Home}/> 
          <AuthRoute path="/new-user" component={NewUser}/>
          <ResetAuthRoute path="/verify" component={Verify}/>
          <ResetAuthRoute path="/verify-link" component={VerifyLink}/>
          <Route path="/not-you" component={NotYou}/>
          <Redirect to={{pathname: "/oh-no", state: { errorCode: "not_found"}}}/>
        </Switch>
      </div>
    )
  }
}