import React from "react"
import { Redirect, Link } from 'react-router-dom'
import Container from 'react-bootstrap/Container'
import Image from 'react-bootstrap/Image'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Fade from 'react-bootstrap/Fade'
import FormGroup from "react-bootstrap/esm/FormGroup"
import api, { 
  ERR_UNEXPECTED,
} from "../api"

export default class NewUser extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      user: this.props.location?.state?.user,
    }



    this.GetUser = this.GetUser.bind(this)
    
    if (!this.state.user) {
        this.GetUser()
    }
    
  }

  GetUser() {
    api.GetUser()
    .then(res => {
      if (res.ok) {
        this.setState({
          user: res.body.data
        })
        return
      }
      this.setState({errorPage: true, errorCode: ERR_UNEXPECTED, errorMessage: res.errorMessage})
    })
    .catch(e => {
      console.log(e)
    })
  }

  render() {
    if (this.state.user.id && !needsUsername(this.state.user)) {
        return <Redirect to={{pathname: '/', state: {user: this.state.user}}}/>
    }
    return (
    <>
        <Container fluid>
          <Row className="d-flex align-items-center landing-row head-space">
            <Col>
              <Row className="d-flex align-items-center text-center pt-3 foot-space">
                <Col>
                  <h4 className="orange">Welcome!</h4> 


                  {/* TODO: have the user fill out name and username */}

                  
                </Col>
              </Row>
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}

function needsUsername(user) {
  return user.username.startsWith("_")
}