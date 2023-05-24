import React from "react";
import api from "../api"

export default class NameForm extends React.Component {
  constructor(props) {
    super(props)
    this.state = {value: ''}

    this.handleChange = this.handleChange.bind(this)
    this.handleSubmit = this.handleSubmit.bind(this)
  }

  handleChange(event) {
    this.setState({value: event.target.value})
  }

  handleSubmit(event) {
    api.SignUp(this.state.value)
    .then(user => alert('Thanks for signing up! A verification email is on it\'s way.'))
    .catch(err => alert(err))
    event.preventDefault()
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit}>
        <label>
          Email address:
          <input type="text" value={this.state.value} onChange={this.handleChange} />
        </label>
        <input type="submit" value="Sign Up" />
      </form>
    )
  }
}