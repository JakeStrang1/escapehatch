import React from "react"

export default class LinkButton extends React.Component {  
  render() {
    return (
      <button className="link-button" {...this.props}>
        {this.props.text}
      </button>
    )
  }
}
