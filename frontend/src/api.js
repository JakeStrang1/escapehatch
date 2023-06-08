const api = {}
const host = process.env.REACT_APP_API_HOST

// Error Codes
export const ERR_BAD_REQUEST = "bad_request"
export const ERR_INTERNAL = "internal"
export const ERR_INVALID = "invalid"
export const ERR_NOT_FOUND = "not_found"
export const ERR_UNEXPECTED = "unexpected"

// Auth error codes
export const ERR_CHALLENGE_ALREADY_USED = "challenge_already_used"
export const ERR_CHALLENGE_EXPIRED = "challenge_expired"
export const ERR_CHALLENGE_FAILED = "challenge_failed"
export const ERR_CHALLENGE_INVALIDATED = "challenge_invalidated"
export const ERR_CHALLENGE_MAX_ATTEMPTS = "challenge_max_attempts"
export const ERR_DO_NOT_CONTACT = "do_not_contact"
export const ERR_EMAIL_FLAGGED = "email_flagged"
export const ERR_EMAIL_INVALID = "email_invalid"
export const ERR_EMAIL_UNAVAILABLE = "email_unavailable"
export const ERR_IP_FLAGGED = "ip_flagged"
export const ERR_SESSION_EXPIRED = "session_expired"
export const ERR_SESSION_INVALID = "session_invalid"
export const ERR_UNAUTHENTICATED = "unauthenticated"

api.SignUp = email => {
    return fetch(host + '/auth/sign-up', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({ email: email })
    })
    .then(FormatResponse)
}

api.SignIn = email => {
    return fetch(host + '/auth/sign-in', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({ email: email })
    })
    .then(FormatResponse)
}

api.Verify = (secret, email, isHash) => {
    let body = {
        secret: secret
    }
    
    if (isHash) {
        body.email_hash = email
    } else {
        body.email = email
    }
    
    return fetch(host + '/auth/verify', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include', // Needed to properly set cookie from response
        body: JSON.stringify(body)
    })
    .then(FormatResponse)
}

api.NotYou = (secret, emailHash, doNotContact) => {
    let body = {
        secret: secret,
        email_hash: emailHash
    }
    if (doNotContact === true) {
        body.do_not_contact = true
    }

    return fetch(host + '/auth/not-you', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include', // Needed to properly set cookie from response
        body: JSON.stringify(body)
    })
    .then(FormatResponse)
}

api.GetUser = () => {
    return fetch(host + '/users/me', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
    })
    .then(FormatResponse)
}

api.UpdateUser = (fullName, username) => {
    let body = {
        full_name: fullName,
        username: username
    }
    return fetch(host + '/users/me', {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify(body)
    })
    .then(FormatResponse)
}

api.ValidateUsername = (username) => {
    return fetch(host + '/validate-username?u=' + encodeURIComponent(username), {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
    })
    .then(FormatResponse)
}

// If request is successful, FormatResponse returns:
// {
//     ok: true,
//     status: <status>,
//     body: <body>
// }
//
// Otherwise returns:
// {
//     ok: false,
//     status: <status>,
//     errorCode: <code>,
//     errorMessage: <message>
// }
async function FormatResponse(res) {
    let apiResponse = {
        ok: res.ok,
        status: res.status,
    }

    let body = await res.json()

    if (!res.ok) {
        if (!body.code || !body.error) {
            throw new Error("Malformed API error: " + JSON.stringify(body))
        }
        apiResponse.errorCode = body.code
        apiResponse.errorMessage = body.error
        return apiResponse
    }

    apiResponse.body = body
    return apiResponse
}

export default api