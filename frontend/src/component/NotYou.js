import React from "react"
import { Redirect, Link } from 'react-router-dom'
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import LinkButton from "./LinkButton"
import api, {
  ERR_UNEXPECTED
} from "../api"

export default class NotYou extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      email: '',
      secret: '',
      content: 'not-you',
      submitDoNotContact: false,
      errorPage: false,
      errorCode: '',
      errorMessage: '',
      reportSuspicious: false,
      submitReportSuspicious: false
    }
    let queryParams = new URLSearchParams(this.props.location.search)
    if (queryParams.get("emailHash")) {
      this.state.email = queryParams.get("emailHash")
    }
    if (queryParams.get("secret")) {
      this.state.secret = queryParams.get("secret")
    }
    this.doNotContact = this.doNotContact.bind(this)
    this.toggleReportSuspicious = this.toggleReportSuspicious.bind(this)
    this.reportSuspicious = this.reportSuspicious.bind(this)
  }

  doNotContact(event) {
    api.NotYou(this.state.secret, this.state.email, true)
    .then(res => {
      if (res.ok) {
        this.setState({
          submitDoNotContact: true,
          reportSuspicious: false,
          submitReportSuspicious: false
        })
        return
      }
      this.setState({errorPage: true, errorCode: ERR_UNEXPECTED, errorMessage: res.errorMessage})
    })
    .catch(e => {
      console.log(e)
    })
    event.preventDefault()
  }

  toggleReportSuspicious() {
    this.setState({
      reportSuspicious: true,
      submitDoNotContact: false,
      submitReportSuspicious: false
    })
  }

  reportSuspicious(event) {
    api.NotYou(this.state.secret, this.state.email, false)
    .then(res => {
      if (res.ok) {
        this.setState({
          submitReportSuspicious: true,
          submitDoNotContact: false,
          reportSuspicious: false
        })
        return
      }
      this.setState({errorPage: true, errorCode: ERR_UNEXPECTED, errorMessage: res.errorMessage})
    })
    .catch(e => {
      console.log(e)
    })
    event.preventDefault()
  }

  render() {
    const NotYouComponent = (
      <>
        <h2>Sorry about that!</h2>
        <br/>
        <p>
          You can <LinkButton text="click here" onClick={this.doNotContact}/> to permanently disable all emails from our site.
          <br/>
          To keep using the site, you can <LinkButton text="report suspicious activity" onClick={this.toggleReportSuspicious}/> instead.
        </p>
        <small>
          Nevermind.&nbsp;
          <Link to="/">
            Take me back.  
          </Link>
        </small>
      </>
    )

    const SubmitDoNotContactComponent = (
      <>
        <h2>You won't hear from us again.</h2>
        <br/>
        <i className="memo">
          All the best to you!
        </i>
      </>
    )

    const ReportSuspiciousComponent = (
      <>
        <h2>Is someone pretending to be you?</h2>
        <br/>
        <p>
          If you click the link below, we'll log you out on all devices, and prevent any new sign-ins for the next 30 minutes.
          <br/>We'll also temporarily block the device that triggered the email you received.
          <br/>
          <br/>
          <LinkButton text="Click here" onClick={this.reportSuspicious}/> to log out all devices for 30 minutes.
        </p>
        <small>
          <Link to="/">
            Actually, nevermind.
          </Link>
        </small>
      </>
    )

    const SubmitReportSuspiciousComponent = (
      <>
        <h2>Thanks, we've taken care of it.</h2>
        <br/>
        <p>
          Remember, you won't be able to log in for 30 minutes.
        </p>
        <small>
          <Link to="/">
            Back to site.  
          </Link>
        </small>
      </>
    )

    let Content = NotYouComponent
    if (this.state.errorPage) {
      return <Redirect push to={{pathname: '/oh-nooo', state: {errorCode: this.state.errorCode, errorMessage: this.state.errorMessage}}}/>
    }
    if (this.state.submitDoNotContact) {
      Content = SubmitDoNotContactComponent
    }
    if (this.state.reportSuspicious) {
      Content = ReportSuspiciousComponent
    }
    if (this.state.submitReportSuspicious) {
      Content = SubmitReportSuspiciousComponent
    }
    return (
      <>
        <Container fluid>
          <Row className="d-flex align-items-center text-center pt-3">
              <Col>
                {Content}
              </Col>
          </Row>
        </Container>
      </>
    )
  }
}
