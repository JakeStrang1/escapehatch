import React from "react"
import { useState, useMemo, useEffect } from 'react'
import debounce from 'lodash.debounce';
import { Redirect, Link } from 'react-router-dom'
import Container from 'react-bootstrap/Container'
import Image from 'react-bootstrap/Image'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import logo from './../assets/logo100.png'
import AuthContext from '../authContext'
import Loading from "./Loading"
import Fade from 'react-bootstrap/Fade'
import FormGroup from "react-bootstrap/esm/FormGroup"
import api, { 
  ERR_UNEXPECTED,
} from "../api"

export default class NewUser extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      user: this.props.user,
      page: 1,
    }

    this.handleNext = this.handleNext.bind(this)
    this.handlePrevious = this.handlePrevious.bind(this)
  }

  handleNext() {
    this.setState({page: this.state.page+1})
  }

  handlePrevious() {
    this.setState({page: this.state.page-1})
  }

  render() {
    if (completedSetup(this.state.user)) {
        return <Redirect to={{pathname: '/'}}/>
    }

    let content
    if (this.state.page == 1){
      content = <Page1 number={this.state.user.number} handleNext={this.handleNext}/>
    } else if (this.state.page == 2) {
      content = <Page2 number={this.state.user.number} handlePrevious={this.handlePrevious}/>
    }
    return (
      <>
        <Container fluid>
          <Row className="d-flex align-items-center landing-row text-center pt-3 head-space foot-space">
            <Col className="paragraph-column mx-auto">
              <Image src={logo} fluid/>
              {content}
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}

class Page1 extends React.Component {
  render() {
    return (
      <>
          <h4 style={{color:"white"}} className="mt-4">Welcome, User #{this.props.number}!</h4>       
          <p>
            Wow, user #{this.props.number} eh? That pretty much makes you a legend around here. Don't let it go to your head. 
          </p>

          <h5 style={{color:"#bbb"}} className="mt-5">A quick note...</h5>
          <p>
            You probably know this already User #{this.props.number}, but this is an <strong>unfinished app</strong>. So as you poke around don't be surprised if it blows up. Once things are more ready I'm sure you'll hear aaaall about it &#8212; coders are so needy, amiright??
          </p>

          <h5 style={{color:"#bbb"}} className="mt-5">Well, let's get you set up!</h5>
          <Button variant="primary" className="orange-btn mt-2" onClick={this.props.handleNext}>
              Next
          </Button>
        </>
    )
  }
}

class Page2 extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      user: this.props.user,
      defaultUsername: "",
      username: "",
      page: 1,
      validated: false,
      name: "",
      fullNameErrorMessage: "",
      fullNameSuccessMessage: "",
      usernameErrorMessage: "",
      usernameSuccessMessage: "",
      submitErrorMessage: "",
      success: false
    }

    this.usernameInputRef = null;
    this.fullNameInputRef = null;

    this.onFullNameChange = this.onFullNameChange.bind(this)
    this.handleSubmit = this.handleSubmit.bind(this)
    this.validateUsername = this.validateUsername.bind(this)
    this.validateFullName = this.validateFullName.bind(this)
    this.validateUsernameDebounce = debounce(this.validateUsername, 500)
  }

  onFullNameChange(e) {
    let fullName = e.target.value

    // Remove all non-alphanumeric characters (trim from the ends and replace with '.' for the middle)
    // Source: https://stackoverflow.com/a/20864946
    let newDefaultUsername = fullName.replace(/^[\W_]+/g, '') // Trim start
    newDefaultUsername = newDefaultUsername.replace(/[\W_]+$/g, '') // Trim end
    newDefaultUsername = newDefaultUsername.replace(/[\W_]+/g, '.').toLowerCase(); // Replace with . and lowercase
    if (newDefaultUsername != this.state.defaultUsername) {
      this.setState({defaultUsername: newDefaultUsername})
      this.validateUsernameDebounce({target: {value: newDefaultUsername}}) // invoke onInput event manually because changing default value doesn't trigger it 
    }

    this.validateFullName()
  }

  validateFullName() {
    let fullName = this.fullNameInputRef.value
    if (fullName) {
      this.setState({fullNameErrorMessage: "", fullNameSuccessMessage: "Looks good!"})
      this.fullNameInputRef.setCustomValidity("")
    } else {
      let errorMessage = "Full name is required."
      this.setState({fullNameErrorMessage: errorMessage, fullNameSuccessMessage: ""})
      this.fullNameInputRef.setCustomValidity(errorMessage)
      return
    }
  }

  validateUsername() {
    let username = this.usernameInputRef.value // Don't check e.target.value because if this is called during name change then it only contains the default value - which could be different
    this.setState({validated: true}) // as soon as the form has input, set validated
    
    if (!username) {
      let errorMessage = "Username is required."
      this.setState({usernameErrorMessage: errorMessage, usernameSuccessMessage: ""})
      this.usernameInputRef.setCustomValidity(errorMessage)
      return
    }

    if (username == this.state.username) {
      return // Do nothing if the value hasn't changed (e.g. the default value may have changed but the actual value stayed the same)
    }
    this.setState({username: username})

    api.ValidateUsername(username)
    .then(response => {
      if (response.ok) {
        this.setState({usernameErrorMessage: "", usernameSuccessMessage: "Username is available!"})
        this.usernameInputRef.setCustomValidity("")
        return
      }
      if (response.errorCode == "invalid") {
        let errorMessage = capitalizeFirstLetter(response.errorMessage) + "."
        this.setState({usernameErrorMessage: errorMessage, usernameSuccessMessage: ""})
        this.usernameInputRef.setCustomValidity(errorMessage)
        return
      }
      console.log("Status: " + response.status + ", Code: " + response.errorCode + ", Message: " + response.errorMessage)
    })
    .catch(e => {
      console.log(e)
    })
  }

  handleSubmit(e) {
    e.preventDefault() // prevent refresh

    let errorMessage = "Please fix the issues below."
    // If no inputs yet (obv. must fail, but we also want to show the validation error messages)
    if (this.fullNameInputRef.value == "" || this.usernameInputRef.value == "") {
      this.validateFullName()
      this.validateUsername()
      this.setState({submitErrorMessage: errorMessage})
      return
    }

    // Otherwise, validation will have already run on each input - we can user the error messages
    if (this.fullNameErrorMessage || this.usernameErrorMessage) {
      // Already showing form-level error messages, so just add in the submit message
      this.setState({submitErrorMessage: errorMessage})
      return
    }

    // Success!
    this.setState({submitErrorMessage: ""})
    api.UpdateUser(this.fullNameInputRef.value, this.usernameInputRef.value)
    .then(response => {
      if (response.ok) {
        this.setState({success: true})
      }
      console.log("Status: " + response.status + ", Code: " + response.errorCode + ", Message: " + response.errorMessage)
    })
    .catch(e => {
      console.log(e)
    })
  }

  render() {
    if (this.state.success) {
      return(
        <AuthContext.Consumer>
          {auth => { // defined in index.js
              auth.attempted = false   
              return (<Redirect to={"/new-user"}/>)
          }}
        </AuthContext.Consumer>
      )
    }
    return (
      <>
          <h4 style={{color:"white"}} className="mt-4">Fill in your profile...</h4>       
          
          <Form className="mt-5" noValidate validated={this.state.validated} onSubmit={this.handleSubmit}>
          <p style={{color:"red"}}>{this.state.submitErrorMessage}</p>
            <Form.Group controlId="formName" className="medium-input mx-auto text-left">
              <Form.Label style={{color:"white"}}>Full Name</Form.Label>
              <Form.Control type="text" placeholder="Your Full Name" onInput={this.onFullNameChange} ref={ref => this.fullNameInputRef = ref}/>
              <Form.Text className="text-muted">
                For the love of God, please capitalize properly.
              </Form.Text>
              <Form.Control.Feedback type="invalid">{this.state.fullNameErrorMessage}</Form.Control.Feedback>
              <Form.Control.Feedback>{this.state.fullNameSuccessMessage}</Form.Control.Feedback>
            </Form.Group>

            <Form.Group controlId="formUsername" className="medium-input mx-auto text-left">
              <Form.Label style={{color:"white"}}>Username</Form.Label>
              <Form.Control isInvalid={!!this.state.usernameErrorMessage} type="text" placeholder="user.name01" className="lowercase" defaultValue={this.state.defaultUsername} onInput={this.validateUsernameDebounce} ref={ref => this.usernameInputRef = ref}/>
              <Form.Text className="text-muted">
                No pressure, you can change this any time!
              </Form.Text>
              <Form.Control.Feedback type="invalid">{this.state.usernameErrorMessage}</Form.Control.Feedback>
              <Form.Control.Feedback>{this.state.usernameSuccessMessage}</Form.Control.Feedback>
              
            </Form.Group>
            <Button variant="primary" className="orange-btn mt-2" type="submit">
              Finish
            </Button>
          </Form>
          <Button variant="link" className="mt-3" onClick={this.props.handlePrevious}>
              Back
          </Button>
        </>
    )
  }
}

function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

function completedSetup(user) {
  return !user.username.startsWith("_") // If the username starts with an underscore, we know they haven't set up their name and username yet
}
