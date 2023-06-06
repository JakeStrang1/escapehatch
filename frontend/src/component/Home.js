import React from "react"
import NavBar from "./NavBar"
import Container from 'react-bootstrap/Container'

export default class Home extends React.Component {
  render() {
    return (
      <>
        <NavBar homeCurrent={true} />
      </>
    )
  }
}