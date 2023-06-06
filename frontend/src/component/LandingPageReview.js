import React from "react"
import Col from 'react-bootstrap/Col'
import Navbar from 'react-bootstrap/Navbar'
import Nav from 'react-bootstrap/Nav'
import Image from 'react-bootstrap/Image'
import banner from './../assets/banner.png'
import stars1 from './../assets/stars-1.png'
import stars2 from './../assets/stars-2.png'
import stars3 from './../assets/stars-3.png'
import stars4 from './../assets/stars-4.png'
import stars45 from './../assets/stars-4-5.png'
import stars5 from './../assets/stars-5.png'

export default class NavBar extends React.Component {  
  constructor(props){
    super(props);
    this.Stars = this.Stars.bind(this)
  }

  Stars() {
    if (this.props.stars == 1) {
        return <Image src={stars1}/>
    }
    if (this.props.stars == 2) {
        return <Image src={stars2}/>
    }
    if (this.props.stars == 3) {
        return <Image src={stars3}/>
    }
    if (this.props.stars == 4) {
        return <Image src={stars4}/>
    }
    if (this.props.stars == 4.5) {
        return <Image src={stars45}/>
    }
    return <Image src={stars5}/>
  }

  render() {
    return (
      <>
        <div className="review-box">
            <this.Stars/>
            <h2 className="review">{this.props.body}</h2>
            <p className="review-author">- {this.props.author}</p>
        </div>
      </>
    )
  }
}
