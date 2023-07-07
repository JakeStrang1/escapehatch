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

api.SignOut = () => {
    return fetch(host + '/auth/sign-out', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include' // Needed to properly set cookie from response
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

api.GetUsers = (search) => {
    let queries = []
  
    if (search) {
      queries.push("search=" + encodeURIComponent(search))
    }

    let queryString = ""
    if (queries.length > 0) {
      queryString = "?" + queries.join("&")
    }

    return fetch(host + '/users' + queryString, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
    })
    .then(FormatResponse)
}

api.GetUser = (userId) => {
    if (!userId) {
        userId = "me" // Fetch self user
    }
    return fetch(host + '/users/' + userId, {
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

api.GetItems = () => {
    return fetch(host + '/items', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
    })
    .then(FormatResponse)
}

api.AddItem = itemId => {
    return fetch(host + '/items/' + itemId + "/add", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include'
    })
    .then(FormatResponse)
}

api.FollowUser = userId => {
    return fetch(host + '/users/' + userId + "/follow", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include'
    })
    .then(FormatResponse)
}

api.UnfollowUser = userId => {
    return fetch(host + '/users/' + userId + "/unfollow", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include'
    })
    .then(FormatResponse)
}

api.RemoveUser = userId => {
    return fetch(host + '/users/' + userId + "/remove", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include'
    })
    .then(FormatResponse)
}

api.GetFollowers = (userId, search) => {
    let searchQuery = ""
    if (search) {
        searchQuery = "?search=" + encodeURIComponent(search)
    }
    return fetch(host + '/users/' + userId + '/followers' + searchQuery, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
    })
    .then(FormatResponse)
}

api.GetFollowing = (userId, search) => {
    let searchQuery = ""
    if (search) {
        searchQuery = "?search=" + encodeURIComponent(search)
    }
    return fetch(host + '/users/' + userId + '/following' + searchQuery, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
    })
    .then(FormatResponse)
}

api.CreateBook = (body) => {
    if (body.image_url) {
        return fetch(host + '/books', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify(body)
        })
        .then(FormatResponse)
    } else {
        const formData  = new FormData()
        formData.append("image_file", body.image_file)
        formData.append("title", body.title)
        formData.append("author", body.author)
        formData.append("published_year", body.published_year)
        formData.append("is_series", body.is_series)
        formData.append("series_title", body.series_title)
        formData.append("sequence_number", body.sequence_number)

        return fetch(host + '/books', {
            method: 'POST',
            credentials: 'include',
            body: formData
        })
        .then(FormatResponse)
    }
}

api.CreateMovie = (body) => {
    if (body.image_url) {
        return fetch(host + '/movies', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify(body)
        })
        .then(FormatResponse)
    } else {
        const formData  = new FormData()
        formData.append("image_file", body.image_file)
        formData.append("title", body.title)
        formData.append("lead_actors", body.lead_actors[0])
        formData.append("lead_actors", body.lead_actors[1])
        formData.append("published_year", body.published_year)
        formData.append("is_series", body.is_series)
        formData.append("series_title", body.series_title)
        formData.append("sequence_number", body.sequence_number)

        return fetch(host + '/movies', {
            method: 'POST',
            credentials: 'include',
            body: formData
        })
        .then(FormatResponse)
    }
}

api.CreateTvSeries = (body) => {
    if (body.image_url) {
        return fetch(host + '/tv-series', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify(body)
        })
        .then(FormatResponse)
    } else {
        const formData  = new FormData()
        formData.append("image_file", body.image_file)
        formData.append("title", body.title)
        formData.append("lead_actors", body.lead_actors[0])
        formData.append("lead_actors", body.lead_actors[1])
        formData.append("tv_series_start_year", body.tv_series_start_year)
        formData.append("tv_series_end_year", body.tv_series_end_year)

        return fetch(host + '/tv-series', {
            method: 'POST',
            credentials: 'include',
            body: formData
        })
        .then(FormatResponse)
    }
}

api.Search = (searchText) => {
    return fetch(host + '/search?search=' + encodeURIComponent(searchText), {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
    })
    .then(FormatResponse)
}

// Use this to wrap a request function if you only want the "data" field.
// E.g. api.Data(api.GetUser()).then(user => { do something })
api.Data = async (promise) => {
    return promise.then(response => {
        if (response.ok) {
            return response.body.data
        }
        console.log("Status: " + response.status + ", Code: " + response.errorCode + ", Message: " + response.errorMessage)
    })
    .catch(e => {
        console.log(e)
    })
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