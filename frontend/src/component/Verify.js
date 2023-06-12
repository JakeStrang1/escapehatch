import React from "react"
import Container from 'react-bootstrap/Container'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import FormGroup from "react-bootstrap/esm/FormGroup"
import Image from 'react-bootstrap/Image'
import keyImage from './../assets/key.png'
import { Redirect, Link } from 'react-router-dom'
import api, { 
  ERR_CHALLENGE_FAILED,
  ERR_NOT_FOUND,
  ERR_UNEXPECTED 
} from "../api"
import { connect } from "react-redux";
import { clearAuth } from "./../reducers/auth"
import { updateUser } from "./../reducers/user"

function mapStateToProps(state) {
  const user = state.user.value
  const auth = state.auth.value
  return {
    auth,
    user
  };
}

const mapDispatchToProps = (dispatch) => {
  return {
    clearAuth: () => dispatch(clearAuth()),
    updateUser: () => dispatch(updateUser())
  }
};

class Verify extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      wordOne: '',
      wordTwo: '',
      wordThree: '',
      wordFour: '',
      submit: false,
      error: '',
      email: '',
      isHash: false,
      fromSignIn: false,
      isValidClass: '',
      placeholderOneDefault: 'super',
      placeholderTwoDefault: 'duper',
      placeholderThreeDefault: 'secret',
      placeholderFourDefault: 'phrase',
      placeholderOne: '',
      placeholderTwo: '',
      placeholderThree: '',
      placeholderFour: '',
      FeedbackComponent: EmptyFeedback
    }

    if (this.props.auth.status != "") {
      this.props.clearAuth()
      this.props.updateUser(null)
    }

    this.state.placeholderOne = this.state.placeholderOneDefault
    this.state.placeholderTwo = this.state.placeholderTwoDefault
    this.state.placeholderThree = this.state.placeholderThreeDefault
    this.state.placeholderFour = this.state.placeholderFourDefault

    // Linked from another page with props
    if (this.props.location.state) {
      this.state.email = this.props.location.state.email
      this.state.fromSignIn = this.props.location.state.fromSignIn // For navigation link to go back
    }

    this.handleChangeOne = this.handleChangeOne.bind(this)
    this.handleChangeTwo = this.handleChangeTwo.bind(this)
    this.handleChangeThree = this.handleChangeThree.bind(this)
    this.handleChangeFour = this.handleChangeFour.bind(this)
    this.handleSubmit = this.handleSubmit.bind(this)
    this.handleError = this.handleError.bind(this)
    this.validate = this.validate.bind(this)
    this.inputOneRef = React.createRef()
    this.inputTwoRef = React.createRef()
    this.inputThreeRef = React.createRef()
    this.inputFourRef = React.createRef()

    let queryParams = new URLSearchParams(this.props.location.search)
    if (queryParams.get("emailHash")) {
      this.state.email = queryParams.get("emailHash")
      this.state.isHash = true
    }
  }

  componentDidMount() {
    if (!this.state.email) {
      return // Show error page instead
    }
    let queryParams = new URLSearchParams(this.props.location.search)
    if (queryParams.get("secret")) {
      this.inputOneRef.current.value = queryParams.get("secret")
      this.handleChangeOne()
    }
  }

  componentDidUpdate(prevProps, prevState) {
    // If any field has a value, clear the placeholders for all fields
    if (prevState.wordOne === "" && prevState.wordTwo === "" && prevState.wordThree === "" && prevState.wordFour === ""
      && (this.state.wordOne !== "" || this.state.wordTwo !== "" || this.state.wordThree !== "" || this.state.wordFour !== "")) {
        this.setState({
          placeholderOne: '',
          placeholderTwo: '',
          placeholderThree: '',
          placeholderFour: ''
        })
    }

    // If all fields don't have a value, add the placeholders back
    if (this.state.wordOne === "" && this.state.wordTwo === "" && this.state.wordThree === "" && this.state.wordFour === ""
      && (prevState.wordOne !== "" || prevState.wordTwo !== "" || prevState.wordThree !== "" || prevState.wordFour !== "")) {
        this.setState({
          placeholderOne: this.state.placeholderOneDefault,
          placeholderTwo: this.state.placeholderTwoDefault,
          placeholderThree: this.state.placeholderThreeDefault,
          placeholderFour: this.state.placeholderFourDefault,
        })
    }
  }

  handleChangeOne() {
    let value = this.inputOneRef.current.value
    let spaceIndex = value.indexOf(" ")
    if (spaceIndex === -1) {
      this.setState({wordOne: value})
      this.validate()
      return
    }
    
    this.setState({wordOne: value.substring(0, spaceIndex)})
    this.inputTwoRef.current.value = value.substring(spaceIndex+1) + this.state.wordTwo
    this.inputTwoRef.current.focus()
    this.handleChangeTwo()
    return
  }

  handleChangeTwo() {
    let value = this.inputTwoRef.current.value
    let spaceIndex = value.indexOf(" ")
    if (spaceIndex === -1) {
      this.setState({wordTwo: value})
      this.validate()
      return
    }
    
    this.setState({wordTwo: value.substring(0, spaceIndex)})
    this.inputThreeRef.current.value = value.substring(spaceIndex+1) + this.state.wordThree
    this.inputThreeRef.current.focus()
    this.handleChangeThree()
    return
  }

  handleChangeThree() {
    let value = this.inputThreeRef.current.value
    let spaceIndex = value.indexOf(" ")
    if (spaceIndex === -1) {
      this.setState({wordThree: value})
      this.validate()
      return
    }
    
    this.setState({wordThree: value.substring(0, spaceIndex)})
    this.inputFourRef.current.value = value.substring(spaceIndex+1) + this.state.wordFour
    this.inputFourRef.current.focus()
    this.handleChangeFour()
    return
  }

  handleChangeFour() {
    let value = this.inputFourRef.current.value
    let spaceIndex = value.indexOf(" ")
    if (spaceIndex === -1) {
      this.setState({wordFour: value})
    } else {
      this.setState({wordFour: value.substring(0, spaceIndex)})
    }

    this.validate()
  }

  handleSubmit(event) {
    let valid = this.validate()
    if (!valid) {
      this.setState({
        FeedbackComponent: RequiredError
      })
      event.preventDefault()
      event.stopPropagation()
      return
    }

    let words = [this.state.wordOne, this.state.wordTwo, this.state.wordThree, this.state.wordFour]
    let secret = words.join(" ")
    api.Verify(secret, this.state.email, this.state.isHash)
    .then(res => {
      if (res.ok) {
        this.setState({
          FeedbackComponent: EmptyFeedback,
          submit: true,
          metadata: res.body.metadata
        })
        return
      }
      this.handleError(res.errorCode, res.errorMessage)
    })
    .catch(e => {
      console.log(e)
    })
    event.preventDefault()
  }

  handleError(code, message) {
    switch (code) {
      case ERR_CHALLENGE_FAILED:
        this.setState({
          isValidClass: "is-invalid",
          FeedbackComponent: ChallengeFailedError
        })
        break
      case ERR_NOT_FOUND:
        // If the email isn't found, it should show an Unexpected error instead
        this.setState({errorPage: true, errorCode: ERR_UNEXPECTED, errorMessage: message})
        break
      default:
        this.setState({errorPage: true, errorCode: code, errorMessage: message})
        break
    }
  }

  validate() {
    if (
      this.inputOneRef.current.checkValidity() === true
      && this.inputTwoRef.current.checkValidity() === true
      && this.inputThreeRef.current.checkValidity() === true
      && this.inputFourRef.current.checkValidity() === true
    ) {
      this.setState({
        isValidClass: "is-valid" // Sets the validation icon on last input field
      })
      return true
    }

    // Only set to invalid if there was already a validation state defined
    this.setState((state) => {
      if (state.isValidClass !== "") {
        return {isValidClass: "is-invalid"}
      }
    })
    return false
  }

  render() {
    if (this.state.submit) {
      return <Redirect to="/"/>
    }
    if (this.state.errorPage) {
      return <Redirect push to={{pathname: '/oh-no', state: {errorCode: this.state.errorCode, errorMessage: this.state.errorMessage}}}/>
    }
    let backLink = "/sign-up"
    let backText = "Back to sign up page."
    if (this.state.fromSignIn) {
      backLink = "/sign-in"
      backText = "Back to sign in page."
    }
    const Feedback = this.state.FeedbackComponent
    const PassPhraseContent = (
      <>
        <Form className="verifyForm" noValidate onSubmit={this.handleSubmit} autoComplete="off">
          <Form.Group controlId="form" className="text-left">
            <Form.Row>
              <Col>
                <Form.Label className="sr-only">
                  First word:
                </Form.Label>
                <Form.Control type="text" placeholder={this.state.placeholderOne} ref={this.inputOneRef} value={this.state.wordOne} required
                  onChange={this.handleChangeOne} className="lowercase"/>
              </Col>
              <Col>
                <Form.Label className="sr-only">
                Second word:
                </Form.Label>
                <Form.Control type="text" placeholder={this.state.placeholderTwo} ref={this.inputTwoRef} value={this.state.wordTwo} required
                  onChange={this.handleChangeTwo} className="lowercase"/>
              </Col>
              <Col>
                <Form.Label className="sr-only">
                Third word:
                </Form.Label>
                <Form.Control type="text" placeholder={this.state.placeholderThree} ref={this.inputThreeRef} value={this.state.wordThree} required
                  onChange={this.handleChangeThree} className="lowercase"/>
              </Col>
              <Col>
                <Form.Label className="sr-only">
                Fourth word:
                </Form.Label>
                <Form.Control className={"no-validation-border lowercase" + this.state.isValidClass} type="text" placeholder={this.state.placeholderFour} ref={this.inputFourRef} value={this.state.wordFour} required
                  onChange={this.handleChangeFour}/>
              </Col>
            </Form.Row>
            <small className="text-danger">
              <Feedback/>
            </small>
          </Form.Group>
          <FormGroup controlId="formSubmit">
            <Button variant="primary" type="submit" className="orange-btn">
              Submit
            </Button>
          </FormGroup>
        </Form>
        <small>
          <Link to={backLink} className="gray">
            {backText} 
          </Link>
        </small>
      </>
    )
    const ErrorContent = (
      <>
        <br/>
        <h5>
          Sorry, looks like something went wrong!
        </h5>
        <br/>
        <small>
          <Link to="/">
            Back to site.
          </Link>
        </small>
      </>
    )
    let Content = PassPhraseContent
    if (!this.state.email) {
      Content = ErrorContent
    }
    return (
      <>
        <Container fluid>
          <Row className="d-flex align-items-center landing-row head-space">
            <Col>
              <Row className="d-flex align-items-center text-center pt-3">
                <Col xs={12}>
                  <Image src={keyImage} fluid/>
                </Col>
              </Row>
              <Row className="d-flex align-items-center text-center pt-3 foot-space">
                <Col>
                  <h4 className="orange">We've sent you an email with a secret phrase:</h4>
                  {Content}
                </Col>
              </Row>
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}

function EmptyFeedback() {
  return (
    <p>&nbsp;</p>
  )
}

function RequiredError() {
  return (
    <p>This can't be right, you're missing some words!</p>
  )
}

function ChallengeFailedError() {
  return (
    <p>Oops! That's not the right phrase. Give it another try.</p>
  )
}

export default connect(mapStateToProps, mapDispatchToProps)(Verify)