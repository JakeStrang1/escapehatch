import React from "react"
import { Link } from 'react-router-dom'
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import BackButton from "./BackButton"
import { 
  ERR_NOT_FOUND,
  ERR_DO_NOT_CONTACT,
  ERR_EMAIL_FLAGGED,
  ERR_IP_FLAGGED,
  ERR_CHALLENGE_ALREADY_USED,
  ERR_CHALLENGE_EXPIRED,
  ERR_CHALLENGE_INVALIDATED,
  ERR_CHALLENGE_MAX_ATTEMPTS
} from "../api"

export default class ErrorPage extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      errorCode: '',
      errorMessage: ''
    }
    if (this.props.location.state) {
      this.state = {
        errorCode: this.props.location.state.errorCode,
        errorMessage: this.props.location.state.errorMessage,
      }
    }
  }

  render() {
    const ErrorComponent = errorComponent(this.state.errorCode)
    console.log(this.state.errorCode + ": " + this.state.errorMessage)
    return (
      <>
        <Container fluid>
          <Row className="d-flex align-items-center text-center pt-3">
              <Col>
                <ErrorComponent/>
              </Col>
          </Row>
        </Container>
      </>
    )
  }
}

function errorComponent(code) {
  switch(code) {
    case ERR_NOT_FOUND:
      return NotFoundError
    case ERR_DO_NOT_CONTACT:
      return DoNotContactError
    case ERR_EMAIL_FLAGGED:
      return EmailFlaggedError
    case ERR_IP_FLAGGED:
      return IPFlaggedError
    case ERR_CHALLENGE_ALREADY_USED:
      return ChallengeAlreadyUsedError
    case ERR_CHALLENGE_EXPIRED:
      return ChallengeExpiredError
    case ERR_CHALLENGE_INVALIDATED:
      return ChallengeInvalidatedError
    case ERR_CHALLENGE_MAX_ATTEMPTS:
      return ChallengeMaxAttemptsError
    default:
      return UnexpectedError
  }
}

function NotFoundError() {
  return (
      <>
          <h2>Whoops! Not Found</h2>
          <br/>
            <Link to="/">
              Take me back.
            </Link>
      </>
  )
}

function DoNotContactError() {
  return (
      <>
          <h2>Your email address is on our do-not-contact list</h2>
          <p>
            <Link to="/support/unblock-email">Click here</Link> to get sorted out.
          </p>
          <small>
            <Link to="/">
              Back to the site.
            </Link>
          </small>
      </>
  )
}

function EmailFlaggedError() {
  return (
      <>
          <h2>Your email address has been flagged for suspicious activity</h2>
          <p>Please try again later.</p>
          <small>
            <Link to="/">
              Back to the site.
            </Link>
          </small>
      </>
  )
}

function IPFlaggedError() {
  return (
      <>
          <h2>Your IP address has been flagged for suspicious activity</h2>
          <p>Please try again later.</p>
          <small>
            <Link to="/">
              Back to the site.
            </Link>
          </small>
      </>
  )
}

function UnexpectedError() {
    return (
        <>
            <h2>Welp! That was unexpected.</h2>
            <p>Tell you what, let's try that again later.</p>
            <small>
              <Link to="/">
                Back to the site.
              </Link>
            </small>
        </>
    )
}

function ChallengeAlreadyUsedError() {
  return (
      <>
          <h2>You already used that code!</h2>
          <p>Double check your inbox, that code has already been used up.</p>
          <small>
            <BackButton text="Go back and try again."/>
          </small>
      </>
  )
}

function ChallengeExpiredError() {
  return (
      <>
          <h2>Your login code is expired</h2>
          <p>You'll need to sign in again.</p>
          <small>
            <Link to="/">
              Back to the site.
            </Link>
          </small>
      </>
  )
}

function ChallengeInvalidatedError() {
  return (
      <>
          <h2>Your login code is expired</h2>
          <p>You'll need to sign in again.</p>
          <small>
            <Link to="/">
              Back to the site.
            </Link>
          </small>
      </>
  )
}

function ChallengeMaxAttemptsError() {
  return (
      <>
          <h2>Too many attempts</h2>
          <p>You'll need to sign in again.</p>
          <small>
            <Link to="/">
              Back to the site.
            </Link>
          </small>
      </>
  )
}