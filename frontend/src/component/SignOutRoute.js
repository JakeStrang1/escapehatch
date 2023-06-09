import React, { useState }from 'react'
import { Redirect } from 'react-router-dom'
import AuthContext from '../authContext'
import Loading from "./Loading"



// SignOutRoute calls Sign Out on the backend and then redirects to the sign in page.
const SignOutRoute = () => {
    const forceUpdate = useForceUpdate()
    return (
        <AuthContext.Consumer>
            {auth => { // defined in index.js
                if (!auth.signOutAttempted) {
                    Promise.all([auth.SignOut(), wait(0)]) // add a wait time to force a minimum load time
                    .then(() => forceUpdate())
                    return (<Loading/>)
                }

                return (
                    <Redirect to={"/sign-in"}/>
                )
            }}
        </AuthContext.Consumer>
    )
}

// Source: https://stackoverflow.com/questions/46240647/react-how-to-force-a-function-component-to-render
function useForceUpdate() {
  const [, setValue] = useState(0) // integer state
  return () => setValue(value => ++value) // update the state to force render
}

function wait(ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

export default SignOutRoute