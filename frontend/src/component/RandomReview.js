import React from "react"
import LandingPageReview from "./LandingPageReview"
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

export default class RandomReview extends React.Component {  
  constructor(props){
    super(props);

    
    let storage = {
      read: [],
      signUp: true
    }
    let storageStr = sessionStorage.getItem("reviews");
    
    if (storageStr) {
      storage = JSON.parse(storageStr)
    }

    let unreadReviews = []

    // If first visit, and on sign up page, show an "explanatory" review
    if (this.props.signUp && storage.read.length == 0) {
      unreadReviews = reviews.filter((review) => {
        return review.explanatory
      })
    } else {
      // Otherwise choose from any unread review
      unreadReviews = reviews.filter((review, i) => {
        return !storage.read.includes(i)
      })
    }

    // Check if this is a "Sign In <-> Sign Up" page situation - in which case we want to keep the same review on screen.
    if (storage.signUp != this.props.signUp && storage.read.length > 0) {
      storage.signUp = this.props.signUp
      sessionStorage.setItem("reviews", JSON.stringify(storage));
      this.randomReview = reviews[storage.read[storage.read.length-1]] // get last read review
      if (unreadReviews.length == 0) {
        this.randomReview = lastReview
      }
      return
    }
    storage.signUp = this.props.signUp

    if (unreadReviews.length == 0) {
      this.randomReview = lastReview
    } else {
      let random = Math.floor(Math.random() * unreadReviews.length)
      this.randomReview = unreadReviews[random]
      storage.read.push(reviews.findIndex(el => el == this.randomReview))
      sessionStorage.setItem("reviews", JSON.stringify(storage));
    }
  }
   
  render() {
    return (
      <LandingPageReview stars={this.randomReview.stars} body={this.randomReview.body} author={this.randomReview.author}/>
    )
  }
}

var lastReview = {
  stars: 4.5,
  body: "I don't even know what this app is, but I just sat here refreshing the page over and over cause these reviews are wild. I read every single one. I can't believe they are real!!1!",
  author: "You",
  explanatory: false
}

var reviews = [
  {
    stars: 4.5,
    body: "So, basically you say what shows you're into so your friends can see. It could be a big thing.",
    author: "Creator",
    explanatory: true
  },
  {
    stars: 1,
    body: "I found out my crush watches Nathan For You. I am devastated.",
    author: "User",
    explanatory: false
  },
  {
    stars: 5,
    body: "It was cool. I looked up a show I like and I clicked 'like' on it. The whole thing worked.",
    author: "User",
    explanatory: true
  },
  {
    stars: 5,
    body: "I thought it was just for shows so I was surprised to see movies and book there too. Way to over deliver! Now I don't know where to begin. Maybe with books?",
    author: "User",
    explanatory: false
  },
  {
    stars: 5,
    body: "I went in and added all the shows I ever watched. It was like seeing my life flash in my eyes. I feel content but also empty inside.",
    author: "User",
    explanatory: true
  },
  {
    stars: 4,
    body: "It's neat that you can add in missing shows and books, especially cause they're pretty much all missing at this point.",
    author: "User",
    explanatory: false
  },
  {
    stars: 2,
    body: "There's no 'enter' button on the Sign Up page, is anyone else seeing this? I don't know how to go to the next part.",
    author: "User",
    explanatory: false
  },
  {
    stars: 1,
    body: "all These People need to stop Escaping into mass media and Come To God ... only He will Save",
    author: "User",
    explanatory: false
  },
  {
    stars: 5,
    body: "At first I thought, \"okaaaaaay... what do we have here??\". Then I learned about it and my questions were answered.",
    author: "User",
    explanatory: true
  },
  {
    stars: 5,
    body: "I love anything where I can make lists, and this is basically a giant list of my shows. I kinda want to make a list of all the sites where you make lists.",
    author: "User",
    explanatory: true
  },
  {
    stars: 5,
    body: "Every time I finish a new show I get to check it off on this site. I've been feeling a lot more accomplished since signing up for this.",
    author: "User",
    explanatory: true
  },
  {
    stars: 3,
    body: "It's just Letterboxd but way worse. The only good quality is there's no ads. I don't think the creator knows how to do them yet. 3 stars I guess.",
    author: "User",
    explanatory: true
  },
  {
    stars: 2,
    body: "Clearly just built by one dude in his basement saying \"I want to be rich! I'll make a website for tv shows!\". Lame. Do something meaningful with your life.",
    author: "User",
    explanatory: false
  },
  {
    stars: 3,
    body: "Landing page was fun but the rest of the site is totally lacking features. I thought I would be able to make lists of TV shows, is this like half finished?",
    author: "User",
    explanatory: false
  },
  {
    stars: 5,
    body: "I'm trying to get my friends to sign up too so we can see each other's favourite shows. Eric if you are seeing this put your email where it says Create an account and press Enter.",
    author: "User",
    explanatory: true
  },
  {
    stars: 5,
    body: "I already added all my books in Goodreads, now I have to add them all over again here? Okay I will.",
    author: "User",
    explanatory: false
  },
  {
    stars: 3,
    body: "Essentially what we have here is just a classic \"social cataloging\" site based around movies, shows, and books. This is what I wish someone had explained to me, so now I'm telling you.",
    author: "User",
    explanatory: true
  },
  {
    stars: 5,
    body: "Super open to having more escape hatch followers so definitely reach out if you see this!! Name is below. Or if it says user I'm so sorry.",
    author: "User",
    explanatory: false
  },
  {
    stars: 5,
    body: "It made sense to me right away. Look up a show, add it to your shelf, and then any time you look at your shelf you can see the show there. It's like the idea just crystallized so hard in my mind.",
    author: "User",
    explanatory: true
  },
  {
    stars: 5,
    body: "I love escape hatch. I don't give a f*** what Keanu is saying about it. Tbh I lost some respect for him when he said that.",
    author: "User",
    explanatory: false
  },
  {
    stars: 5,
    body: "To the people saying this is like Letterboxd, it's not. It is a totally different and unique thing. Also 5 stars.",
    author: "Creator",
    explanatory: false
  },
]
  
