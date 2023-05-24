import React from "react"
import { withRouter } from 'react-router'
import LinkButton from "./LinkButton"

// Source: https://stackoverflow.com/questions/30915173/react-router-go-back-a-page-how-do-you-configure-history
// Source: https://stackoverflow.com/questions/46681387/react-router-v4-how-to-go-back
class BackButton extends React.Component {  
  constructor(props){
    super(props);
    this.goBack = this.goBack.bind(this)
  }
 
  goBack(){
    this.props.history.goBack()
  }

  render() {
    return (
      <LinkButton onClick={this.props.history.goBack} {...this.props}/>
    )
  }
}

export default withRouter(BackButton)