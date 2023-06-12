import React from "react"
import NavBar from "./NavBar"
import UserSummary from "./UserSummary"
import Shelves from "./Shelves"
import { connect } from "react-redux";

class Home extends React.Component {
  constructor(props){
    super(props);
  }

  render() {
    return (
      <>
        <NavBar homeCurrent={true} />
        <UserSummary user={this.props.user}/>
        <Shelves user={this.props.user}/>
      </>
    )
  }
}

function mapStateToProps(state) {
  const user = state.user.value
  return {
    user
  };
}

export default connect(mapStateToProps)(Home);
