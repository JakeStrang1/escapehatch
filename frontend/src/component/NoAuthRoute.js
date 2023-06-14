import React, { useState } from 'react'
import { Redirect } from 'react-router-dom'
import Loading from "./Loading"
import { useSelector, useDispatch } from "react-redux";
import { updateUser } from "./../reducers/user"
import { setAuthPending, setAuthComplete, clearSignOut } from "./../reducers/auth"
import api from './../api'

// NoAuthRoute ensures that authentication has been attempted.
// If authentication was successful, redirects to the home page.
// If authentication fails, proceeds to the given route.
const NoAuthRoute = ({component: Component, ...rest}) => {
  const auth = useSelector(state => state.auth.value);
  const dispatch = useDispatch();

  if (auth.status == "") {
    api.GetUser()
    .then(response => {
      if (response.ok) {
        dispatch(updateUser(response.body.data)) // Update user state
      } else {
        console.log("Status: " + response.status + ", Code: " + response.errorCode + ", Message: " + response.errorMessage)
      }
        dispatch(setAuthComplete())
      }
    )
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
  var child = props.children
  if (user) {
      child = <Redirect to="/"/>
  }
  return (
      <>
          {child}
      </>
  )
}

export default NoAuthRoute