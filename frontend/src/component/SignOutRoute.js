import React, { useState }from 'react'
import { Redirect } from 'react-router-dom'
import Loading from "./Loading"
import { clearAuth, setSignOutPending, setSignOutComplete } from "./../reducers/auth"
import { useSelector, useDispatch } from "react-redux";
import api from './../api'
import { updateUser } from "./../reducers/user"


// SignOutRoute calls Sign Out on the backend and then redirects to the sign in page.
const SignOutRoute = () => {
  const auth = useSelector(state => state.auth.value);
  const dispatch = useDispatch();

  if (auth.signOutStatus == "") {
      api.SignOut()
      .then(response => {
        dispatch(updateUser(null)) // Update user state
        dispatch(setSignOutComplete())
      })
      .catch(e => {
        console.log(e)
      })
      dispatch(setSignOutPending())
      dispatch(clearAuth())
      return (<Loading/>)
  }

  if (auth.signOutStatus == "PENDING") {
      return (<Loading/>)
  }

  if (auth.signOutStatus == "COMPLETE") {
      return (
        <Redirect to={"/sign-in"}/>
      )
  }

  console.error("unknown auth status: " + auth.status)  
}

export default SignOutRoute