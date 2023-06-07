import React from "react"
import { Redirect, Link } from 'react-router-dom'
import api, { 
  ERR_EMAIL_UNAVAILABLE,
  ERR_INVALID,
} from "../api"
import Container from 'react-bootstrap/Container'
import Image from 'react-bootstrap/Image'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Fade from 'react-bootstrap/Fade'
import FormGroup from "react-bootstrap/esm/FormGroup"
import randomEmail from "../util/randomEmail"
import landingImage from './../assets/landing.png'
import RandomReview from "./RandomReview"

export default class SignUp extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      open: false,
      email: '',
      submit: false,
      validated: false,
      apiInvalidClass: '',
      ErrorComponent: EmailValidationError,
      errorPage: false,
      errorCode: '',
      errorMessage: '',
      HelpTextComponent: PrivacyStatement,
      trigger: false,
      placeholderEmail: randomEmail(),
    }

    // Linked from another page with props
    if (this.props.location.state) {
      this.state.email = this.props.location.state.email
      this.state.trigger = this.props.location.state.trigger
      if (this.props.location.state.placeholderEmail) {
        this.state.placeholderEmail = this.props.location.state.placeholderEmail
      }
    }

    this.handleChange = this.handleChange.bind(this)
    this.handleSubmit = this.handleSubmit.bind(this)
    this.formRef = React.createRef()
  }

  componentDidMount() {
    // Trigger sign up submit immediately
    if (this.state.trigger) {
      this.handleSubmit()
    }
  }

  setOpen(open) {
    this.setState({open: open})
  }

  handleChange(event) {
    this.setState({email: event.target.value})
  }

  handleSubmit(event) {
    this.setState({
      apiInvalidClass: '' // Prevent inconsistent spacing
    })
    if (this.formRef.current.checkValidity() === false) {
      if (event) {
        // Won't be defined when manual triggering
        event.preventDefault()
        event.stopPropagation()
      }

      this.setState({
        validated: true,
        ErrorComponent: EmailValidationError,
        HelpTextComponent: FormFeedback
      })
      return
    }

    this.setState({
      validated: true,
      HelpTextComponent: FormFeedback
    })
    api.SignUp(this.state.email)
    .then(res => {
      if (res.ok) {
        this.setState({apiInvalidClass: ""})
        this.setState({submit: true})
        return
      }
      this.handleError(res.errorCode, res.errorMessage)
    })
    .catch(e => {
      console.log(e)
    })
    if (event) {
      event.preventDefault()
    }
  }

  handleError(code, message) {
    switch (code) {
      case ERR_EMAIL_UNAVAILABLE:
        this.setState({
          validated: false, // Must be false when using custom validation
          apiInvalidClass: "is-invalid",
          ErrorComponent: EmailUnavailableError
        })
        break
      case ERR_INVALID:
        this.setState({
          validated: false, // Must be false when using custom validation
          apiInvalidClass: "is-invalid",
          ErrorComponent: EmailValidationBackendError
        })
        break
      default:
        this.setState({errorPage: true, errorCode: code, errorMessage: message})
        break
    }
  }

  render() {
    if (this.state.submit) {
      return <Redirect to={{pathname: '/verify', state: {email: this.state.email}}}/>
    }
    if (this.state.errorPage) {
      return <Redirect to={{pathname: '/oh-nooo', state: {errorCode: this.state.errorCode, errorMessage: this.state.errorMessage}}}/>
    }
    const ErrorMessage = this.state.ErrorComponent
    const HelpText = this.state.HelpTextComponent
    return (
      <>
        <Container fluid>
          <Row className="landing-row">
            <Col xs={12} lg={7} className="left align-self-center">
              <Image src={landingImage} fluid/>
            </Col>
            <Col xs={11} lg={4} className="right align-self-center">
              <RandomReview signUp={true}/>
              <Form className="signUpForm" noValidate validated={this.state.validated} onSubmit={this.handleSubmit} ref={this.formRef}>
                <h4 className="header">Create an account</h4>
                <Form.Group controlId="formEmail">
                  <Form.Label className="sr-only">
                    Email
                  </Form.Label>
                  <Form.Control className={this.state.apiInvalidClass} type="email" required placeholder={this.state.placeholderEmail}
                    onFocus={() => this.setOpen(true)}
                    onBlur={() => this.setOpen(false)}
                    onChange={this.handleChange}/>
                  <HelpText email={this.state.email} open={this.state.open} ErrorMessage={ErrorMessage}/>
                </Form.Group>
                <FormGroup controlId="formSubmit">
                  <Button variant="primary" type="submit" className="sr-only">
                    Submit
                  </Button>
                </FormGroup>
              </Form>
              <small className="text-muted">Already have an account?&nbsp;
                <Link to={{pathname: "/sign-in", state: { email: this.state.email, placeholderEmail: this.state.placeholderEmail }}}>
                  Sign in.
                </Link>
              </small>
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}

function PrivacyStatement(props) {
  return (
    <Fade in={props.open}>
      <Form.Text className="text-muted">
        We'll never share your email with anyone else.
      </Form.Text>
    </Fade>
  )
}

function FormFeedback(props) {
  const ErrorMessage = props.ErrorMessage
  return (
    <>
      <Form.Control.Feedback type="invalid">
        <ErrorMessage email={props.email}/>
      </Form.Control.Feedback>
      <Form.Control.Feedback type="valid">&nbsp;</Form.Control.Feedback> {/* Maintain consistent spacing */}
    </>
  )
}

function EmailValidationBackendError() {
  return (
    <p>Something went wrong on our end. Try something else?</p>
  )
}

function EmailValidationError() {
    return (
      <p>We hit a problem with that email address. Maybe try something else?</p>
    )
}

function EmailUnavailableError(props) {
  return (
    <p>Your account already exists!&nbsp;
      <Link to={{pathname: "/sign-in", state: { email: props.email, trigger: true}}}>
        Sign in with this email.
      </Link>
    </p>
  )
} 