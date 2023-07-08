import React from "react"
import Col from 'react-bootstrap/Col'
import Row from 'react-bootstrap/Row'
import Container from 'react-bootstrap/Container'
import Navbar from 'react-bootstrap/Navbar'
import Nav from 'react-bootstrap/Nav'
import Image from 'react-bootstrap/Image'
import emptyShelfImage from './../assets/empty_shelf.png'
import searchImage from './../assets/search.png'
import scrollLeft from './../assets/scroll-left.png'
import scrollRight from './../assets/scroll-right.png'
import scrollLeftFade from './../assets/scroll-left-fade.png'
import scrollRightFade from './../assets/scroll-right-fade.png'

export default class Shelves extends React.Component {  
  constructor(props){
    super(props);
  }

  render() {
    if (!this.props.user.shelves || this.props.user.shelves.length == 0) {
        return (
            <Empty user={this.props.user}/>
        )
    }

    return (
      <>
        <Container className={this.props.className} fluid>
          <Row>
            <Col xs={12} className="">
              <Row>
                <Col className="ml-0 mr-0 pl-0 pr-0 mt-4">
                  {this.props.user.shelves.filter(shelf => shelf.items.length > 0).map((shelf, index) => {
                    return <Shelf key={index} shelf={shelf}/>
                  })}
                </Col>
              </Row>
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}

class Shelf extends React.Component {
  constructor(props) {
    super(props)

    this.scrollLeftImageRef = React.createRef()
    this.scrollRightImageRef = React.createRef()
    this.scrollLeftLinkRef = React.createRef()
    this.scrollRightLinkRef = React.createRef()
    this.shelfRef = React.createRef()

    this.state = {
      leftMax: true,
      rightMax: true,
      mouseOverShelf: false,
      previousWidth: 0
    }

    this.scrollLeftMouseOver = this.scrollLeftMouseOver.bind(this)
    this.scrollLeftMouseOut = this.scrollLeftMouseOut.bind(this)
    this.scrollRightMouseOver = this.scrollRightMouseOver.bind(this)
    this.scrollRightMouseOut = this.scrollRightMouseOut.bind(this)
    this.shelfMouseOver = this.shelfMouseOver.bind(this)
    this.shelfMouseOut = this.shelfMouseOut.bind(this)
    this.onLeftScrollClick = this.onLeftScrollClick.bind(this)
    this.onRightScrollClick = this.onRightScrollClick.bind(this)
    this.onScroll = this.onScroll.bind(this)
  }

  componentDidMount() {
    this.setState({previousWidth: this.shelfRef.current.scrollWidth})
    this.onScroll() // initialize scroll button state
  }

  componentDidUpdate() {
    // Otherwise the side scroll buttons may not appear on page load
    if (this.state.previousWidth != this.shelfRef.current.scrollWidth) {
      this.setState({previousWidth: this.shelfRef.current.scrollWidth})
      this.onScroll() // recalculate scroll button state
    }
  }

  scrollLeftMouseOver() {
    if (this.scrollLeftImageRef.current) {
      this.scrollLeftImageRef.current.src = scrollLeft
    }
  }

  scrollLeftMouseOut() {
    if (this.scrollLeftImageRef.current) {
      this.scrollLeftImageRef.current.src = scrollLeftFade
    }
  }

  scrollRightMouseOver() {
    if (this.scrollRightImageRef.current) {
      this.scrollRightImageRef.current.src = scrollRight
    }
  }

  scrollRightMouseOut() {
    if (this.scrollRightImageRef.current) {
      this.scrollRightImageRef.current.src = scrollRightFade
    }
  }

  shelfMouseOver(e) {
    this.setState({mouseOverShelf: true})
  }

  shelfMouseOut(e) {
    this.setState({mouseOverShelf: false})
  }

  onLeftScrollClick(e) {
    let newScrollLeft = this.shelfRef.current.scrollLeft - (this.shelfRef.current.offsetWidth - 200)
    newScrollLeft = newScrollLeft < 0 ? 0 : newScrollLeft
    this.shelfRef.current.scrollTo({left: newScrollLeft, behavior: "smooth"})
  }

  onRightScrollClick(e) {
    let newScrollLeft = this.shelfRef.current.scrollLeft + (this.shelfRef.current.offsetWidth - 200)
    newScrollLeft = newScrollLeft > this.shelfRef.current.scrollWidth ? this.shelfRef.current.scrollWidth : newScrollLeft
    this.shelfRef.current.scrollTo({left: newScrollLeft, behavior: "smooth"})
  }

  onScroll() {
    let leftMax = false
    let rightMax = false
    if (this.shelfRef.current.scrollLeft <= 0) {
      leftMax = true
    }

    if (this.shelfRef.current.scrollWidth - this.shelfRef.current.scrollLeft <= this.shelfRef.current.offsetWidth) {
      rightMax = true
    }

    this.setState({leftMax: leftMax, rightMax: rightMax})
  }

  render() {
    return (
      <>
        <div className="pt-3 pl-5 ml-3 shelf-name">
          <h4>{this.props.shelf.name}&nbsp;&nbsp;•&nbsp;&nbsp;{this.props.shelf.items.length}</h4>
        </div>
        <div className="shelf-container" onMouseEnter={this.shelfMouseOver} onMouseLeave={this.shelfMouseOut}>
          <Row className="d-flex flex-row pt-3 pb-3 pl-5 pr-5 justify-content-start shelf" ref={this.shelfRef} onScroll={this.onScroll}>
            {this.props.shelf.items.map((item, index) => {
              return <Item key={index} item={item}/>
            })}
          </Row>
          <Row className="d-flex flex-row shelf-overlay">
            {
              !this.state.leftMax && (
                <div className="d-flex flex-row shelf-scroll-left align-items-center justify-content-center mr-auto" onClick={this.onLeftScrollClick} onMouseOver={this.scrollLeftMouseOver} onMouseOut={this.scrollLeftMouseOut} ref={this.scrollLeftLinkRef}>
                  {
                    this.state.mouseOverShelf && (
                      <Image className="scroll-arrow" src={scrollLeftFade} fluid ref={this.scrollLeftImageRef}/>
                    )
                  }
                </div>
              )
            }
            {
              !this.state.rightMax && (
                <div className="d-flex flex-row shelf-scroll-right align-items-center justify-content-center ml-auto" onClick={this.onRightScrollClick} onMouseOver={this.scrollRightMouseOver} onMouseOut={this.scrollRightMouseOut} ref={this.scrollRightLinkRef}>
                  {
                    this.state.mouseOverShelf && (
                      <Image className="scroll-arrow" src={scrollRightFade} fluid ref={this.scrollRightImageRef}/>
                    )
                  }
                </div>
              )
            }
          </Row>
        </div>
        <Row className="ml-0 mr-0">
          <Col xs={1} md={2} lg={3}></Col>
          <Col className="search-result"></Col>
          <Col xs={1} md={2} lg={3}></Col>
        </Row>
      </>
    )
  }
}

class Item extends React.Component {
  render() {
    return (
      <>
        <div className="d-flex pl-3 shelf-item">
          <Image src={this.props.item.image} className="" fluid rounded/>
        </div>
      </>
    )
  }
}

class Empty extends React.Component {
  render() {
    return (
      <>
        <Col className="align-items-center mt-5">
          <Row className="d-flex">
            <Col xs={1}></Col>
            <Col>
              <Row className="d-flex align-items-center text-center pt-2">
                <Col xs={2}></Col>
                <Col>
                  <Image src={emptyShelfImage} fluid/>
                </Col>
                <Col xs={2}></Col>
              </Row>
              <Row className="d-flex align-items-center text-center pt-4">
                <Col md={10} lg={8} xl={6} className="mx-auto">
                  <h4 className="orange">{this.props.user.self ? "Your shelf is empty!" : "There's nothing here! :("}</h4>
                  <p>{this.props.user.self ? 
                    "Your shelf is where you showcase the shows, movies, and books that you love to escape into." : 
                    (<><em>Is your friend okay?</em> When someone doesn't have favourite shows and movies they may be spending unhealthy amounts of time in the real world.</>)
                    }
                  </p>
                  { this.props.user.self &&
                    (<div className="d-inline-block mt-3">
                      <a href="/search">
                        <div className="outline-link pt-3 pr-3 pl-3 pb-2">
                          <Image src={searchImage}/>
                          <span style={{color: "white"}} className="d-block mt-2">Add to shelf</span>
                        </div>
                      </a>
                    </div>)
                  }
                </Col>
              </Row>
            </Col>
            <Col xs={1}></Col>
          </Row>
        </Col>
      </>
    )
  }
}

// TODO: fix mouse out/over events for shelf to hide/display arrows - the event refires each time the mouse crosses a div boundary (maybe a debounce will help? - or maybe checking the event target?)