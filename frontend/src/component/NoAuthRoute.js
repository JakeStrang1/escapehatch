import React, { useState } from 'react'
import { Redirect } from 'react-router-dom'
import AuthContext from '../authContext'
import Loading from "./Loading"

// NoAuthRoute ensures that authentication has been attempted.
// If authentication was successful, redirects to the home page.
// If authentication fails, proceeds to the given route.
const NoAuthRoute = ({component: Component, ...rest}) => {
    const forceUpdate = useForceUpdate()
    return (
        <AuthContext.Consumer>
            {auth => {
                if (!auth.attempted) {
                    Promise.all([auth.User(), wait(0)]) // add a wait time to force a minimum load time
                    .then(() => forceUpdate())
                    return (<Loading/>)
                }

                return (
                <AuthChecker auth={auth.user}>
                    <Component {...rest}/>
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
        var child = this.props.children
        if (this.props.auth) {
            child = <Redirect to="/"/>
        }
        return (
            <>
                {child}
            </>
        )
    }
}

export default NoAuthRoute