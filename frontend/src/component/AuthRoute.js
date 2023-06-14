import React, { useState }from 'react'
import { Redirect } from 'react-router-dom'
import Loading from "./Loading"
import { useSelector, useDispatch } from "react-redux";
import { updateUser } from "./../reducers/user"
import { setAuthPending, setAuthComplete, clearSignOut } from "./../reducers/auth"
import api from './../api'

// AuthRoute ensures that authentication has been attempted.
// If authentication was successful, proceeds to the given route.
// If authentication fails, redirects to the sign in page.
const AuthRoute = ({component: Component, ...rest}) => {
  const auth = useSelector(state => state.auth.value);
  const dispatch = useDispatch();

  // let refresh = (window.performance && performance.navigation.type == 1) // Detect page refresh vs. navigation
  // TODO: how to avoid an infinite loop situation while using this?

  if (auth.status == "") {
    api.GetUser()
    .then(response => {
      if (response.ok) {
        dispatch(updateUser(response.body.data)) // Update user state
      } else {
        console.log("Status: " + response.status + ", Code: " + response.errorCode + ", Message: " + response.errorMessage)
      }
      dispatch(setAuthComplete())
    })
    .catch(e => {
      console.log(e)
    })
    dispatch(setAuthPending())
    dispatch(clearSignOut())

    return (<Loading/>)
  }

  if (auth.status == "PENDING") {
    return (<Loading/>)
  }

  if (auth.status == "COMPLETE") {
    return (
      <AuthChecker redirect={rest.redirect}>
          <Component {...rest}/>
      </AuthChecker>
    )
  }

  console.error("unknown auth status: " + auth.status)  
}

const AuthChecker = props => {
  const user = useSelector(state => state.user.value);
  var noAuthRedirect = props.redirect ? props.redirect : "/sign-in"
  var child = props.children
  if (!user) {
    child = <Redirect to={noAuthRedirect}/>
  } else if (!completedSetup(user) && window.location.pathname != "/new-user") { // Avoid infinite redirect loop
    child = <Redirect to={"/new-user"}/>
  }
  return (
    <>
      {child}
    </>
  )
}

function completedSetup(user) {
  return !user.username.startsWith("_") // If the username starts with an underscore, we know they haven't set up their name and username yet
}

export default AuthRoute