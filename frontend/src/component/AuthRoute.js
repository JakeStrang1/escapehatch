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
            {auth => {
                if (!auth.attempted) {
                    Promise.all([auth.User(), wait(2000)])
                    .then(() => forceUpdate())
                    return (<Loading/>)
                }

                return (
                <AuthChecker auth={auth.user} redirect={rest.redirect}>
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
        var redirect = this.props.redirect ? this.props.redirect : "/sign-in"
        var child = <Redirect to={redirect}/>
        if (this.props.auth) {
            child = this.props.children
        }
        return (
            <>
                {child}
            </>
        )
    }
}

export default AuthRoute