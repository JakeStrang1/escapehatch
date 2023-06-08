import React, { useState }from 'react'
import { Redirect } from 'react-router-dom'
import AuthContext from '../authContext'
import Loading from "./Loading"



// AuthRoute ensures that authentication has been attempted.
// If authentication was successful, proceeds to the given route.
// If authentication fails, redirects to the sign in page.
const AuthRoute = ({component: Component, ...rest}) => {
    const forceUpdate = useForceUpdate()
    return (
        <AuthContext.Consumer>
            {auth => { // defined in index.js
                if (!auth.attempted) {
                    Promise.all([auth.User(), wait(0)]) // add a wait time to force a minimum load time
                    .then(() => forceUpdate())
                    return (<Loading/>)
                }

                return (
                <AuthChecker user={auth.user} redirect={rest.redirect}>
                    <Component user={auth.user} {...rest}/>
                </AuthChecker>
            )}}
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

class AuthChecker extends React.Component {
  render() {
    var noAuthRedirect = this.props.redirect ? this.props.redirect : "/sign-in"
    var child = this.props.children
    if (!this.props.user) {
      child = <Redirect to={noAuthRedirect}/>
    } else if (!completedSetup(this.props.user) && window.location.pathname != "/new-user") { // Avoid infinite redirect loop
      child = <Redirect to={"/new-user"}/>
    }
    return (
      <>
        {child}
      </>
    )
  }
}

function completedSetup(user) {
  return !user.username.startsWith("_") // If the username starts with an underscore, we know they haven't set up their name and username yet
}

export default AuthRoute