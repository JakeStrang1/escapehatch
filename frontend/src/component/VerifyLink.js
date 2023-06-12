import React from "react";
import api from "../api"
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

class VerifyLink extends React.Component {
  constructor(props) {
    super(props)
    this.state = {value: ''}

    if (this.props.auth.status != "") {
      this.props.clearAuth()
      this.props.updateUser(null)
    }

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

export default connect(mapStateToProps, mapDispatchToProps)(VerifyLink)