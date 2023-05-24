import React from 'react'
import AuthContext from '../authContext'

// ResetAuthRoute resets the authentication state to "not attempted".
// This ensures that the next Auth check will call the server.
const ResetAuthRoute = ({component: Component, ...rest}) => {
    return (
        <AuthContext.Consumer>
            {auth => {
                auth.attempted = false                
                return (<Component {...rest}/>)
            }}
        </AuthContext.Consumer>
    )
}

export default ResetAuthRoute