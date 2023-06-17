import React, { forwardRef } from "react"
import NavBar from "./NavBar"
import { SearchResult } from "./Search"
import UserSummary from "./UserSummary"
import Shelves from "./Shelves"
import Container from 'react-bootstrap/Container'
import Col from 'react-bootstrap/Col'
import Row from 'react-bootstrap/Row'
import Form from 'react-bootstrap/Form'
import Alert from 'react-bootstrap/Alert'
import InputGroup from 'react-bootstrap/InputGroup'
import FormControl from 'react-bootstrap/FormControl'
import Image from 'react-bootstrap/Image'
import Button from 'react-bootstrap/Button'
import { connect } from "./../reducers";
import invalidImage from './../assets/invalid.png'
import { cloneDeep } from 'lodash-es'
import api, { 
  ERR_UNEXPECTED,
} from "../api"


class AddNewItem extends React.Component {
  constructor(props){
    super(props);

    this.defaultPreview = {
      id: "preview",
      media_type: "tv_series",
      image_url: invalidImage,
      user_count: 1,
      description: "Title",
      title: "Title",
      author: "Author",
      published_year: "2000",
      is_series: false,
      series_title: "Series Title",
      sequence_number: "1",
      lead_actors: ["Actor 1", "Actor 2"],
      tv_series_start_year: "2000",
      tv_series_end_year: "Present",
    }

    this.state = {
      type: "TV Series",
      imageType: "url",
      validate: false,
      errorMessage: "",
      isSeries: false,
      imageInvalid: false,
      submitting: false,
      success: false,
      imageFile: null,
      preview: cloneDeep(this.defaultPreview)
    }

    this.formRefs = {
      mediaType: React.createRef(),
      title: React.createRef(),
      isSeries: React.createRef(),
      seriesTitle: React.createRef(),
      seriesNumber: React.createRef(),
      author: React.createRef(),
      publishedYear: React.createRef(),
      leadActor1: React.createRef(),
      leadActor2: React.createRef(),
      startYear: React.createRef(),
      endYear: React.createRef(),
      imageUrl: React.createRef(),
      imageFile: React.createRef(),
      imageResult: React.createRef()
    }

    this.handleCreate = this.handleCreate.bind(this)
    this.handleCreateAndAdd = this.handleCreateAndAdd.bind(this)
    this.handleChange = this.handleChange.bind(this)
    this.handleTypeChange = this.handleTypeChange.bind(this)
    this.handleImageType = this.handleImageType.bind(this)
    this.handleImageFile = this.handleImageFile.bind(this)
    this.handleImageURL = this.handleImageURL.bind(this)
    this.onImageError = this.onImageError.bind(this)
    this.handleSeriesChange = this.handleSeriesChange.bind(this)
    this.tvSeriesBody = this.tvSeriesBody.bind(this)
    this.bookBody = this.bookBody.bind(this)
    this.movieBody = this.movieBody.bind(this)
    this.clearForm = this.clearForm.bind(this)
    this.getPreview = this.getPreview.bind(this)
  }

  // handleChange creates a "preview" result and keeps it updated on every change to the form
  handleChange() {
    this.setState({preview: this.getPreview()})
  }

  getPreview() {
    let preview = cloneDeep(this.defaultPreview)

    if (this.formRefs.imageResult?.current?.src && this.formRefs.imageResult.current.src != window.location.href) {
      preview.image_url = this.formRefs.imageResult.current.src
    } else {
      preview.image_url = invalidImage
    }

    if (this.formRefs.title?.current?.value.trim()) {
      preview.title = this.formRefs.title.current.value.trim()
    }

    if (this.formRefs.author?.current?.value.trim()) {
      preview.author = this.formRefs.author.current.value.trim()
    }

    if (this.formRefs.publishedYear?.current?.value.trim()) {
      preview.published_year = this.formRefs.publishedYear.current.value.trim()
    }

    if (this.formRefs.isSeries?.current?.checked) {
      preview.is_series = this.formRefs.isSeries.current.checked
    }

    if (this.formRefs.seriesTitle?.current?.value.trim()) {
      preview.series_title = this.formRefs.seriesTitle.current.value.trim()
    }

    if (this.formRefs.seriesNumber?.current?.value.trim()) {
      preview.sequence_number = this.formRefs.seriesNumber.current.value.trim()
    }

    if (this.formRefs.leadActor1?.current?.value.trim()) {
      preview.lead_actors[0] = this.formRefs.leadActor1.current.value.trim()
    }

    if (this.formRefs.leadActor2?.current?.value.trim()) {
      preview.lead_actors[1] = this.formRefs.leadActor2.current.value.trim()
    }

    if (this.formRefs.startYear?.current?.value.trim()) {
      preview.tv_series_start_year = this.formRefs.startYear.current.value.trim()
    }

    if (this.formRefs.endYear?.current?.value.trim()) {
      preview.tv_series_end_year = this.formRefs.endYear.current.value.trim()
    }

    if (this.formRefs.mediaType?.current) {
      let mediaType = this.formRefs.mediaType.current.value.trim()
      switch (mediaType) {
        case "TV Series":
          preview.media_type = "tv_series"
          preview.description = preview.title
          break
        case "Movie":
          preview.media_type = "movie"
          if (!preview.is_series) {
            preview.description = preview.title
          } else if (!preview.sequence_number) {
            preview.description = preview.title + " (" + preview.series_title + ")"
          } else if (!isNaN(Number(preview.sequence_number))) {
            preview.description = preview.title + " (" + preview.series_title + " #" + preview.sequence_number + ")"
          } else {
            preview.description = preview.title + " (" + preview.series_title + " " + preview.sequence_number + ")"
          }
          break
        case "Book":
          preview.media_type = "book"
          if (!preview.is_series) {
            preview.description = preview.title
          } else if (!preview.sequence_number) {
            preview.description = preview.title + " (" + preview.series_title + ")"
          } else if (!isNaN(Number(preview.sequence_number))) {
            preview.description = preview.title + " (" + preview.series_title + " #" + preview.sequence_number + ")"
          } else {
            preview.description = preview.title + " (" + preview.series_title + " " + preview.sequence_number + ")"
          }
          break
        default:
          console.error("unknown type: " + mediaType)
          break
      }
    }

    return preview
  }

  handleCreate(e, add) {
    e.preventDefault()
    this.setState({submitting: true, success: false, errorMessage: "", validate: true})

    let result
    let body
    let err
    switch(this.formRefs.mediaType.current.value) {
      case "TV Series":
        [body, err] = this.tvSeriesBody()
        if (err) {
          return
        }
        result = api.CreateTvSeries(body)
        break
      case "Movie":
        [body, err] = this.movieBody()
        if (err) {
          return
        }
        result = api.CreateMovie(body)
        break
      case "Book":
        [body, err] = this.bookBody()
        if (err) {
          return
        }
        result = api.CreateBook(body)
        break
      default:
        console.error("unknown media type: " + this.formRefs.mediaType.current.value)
        return
    }

    result.then(res => {
      window.scrollTo(0, 0)

      if (res.ok) {
        this.clearForm()

        this.setState({
          submitting: false,
          success: true,
          errorMessage: "",
          validate: false,
          preview: this.getPreview() // Form must be cleared first
        })
        
        if (add) {
          api.AddItem(res.body.data.id)
        }
        return
      }

      this.setState({
        submitting: false,
        success: false,
        errorMessage: res.errorMessage
      })
      console.log("Status: " + res.status + ", Code: " + res.errorCode + ", Message: " + res.errorMessage)
    })
  }

  handleCreateAndAdd(e) {
    this.handleCreate(e, true)
  }

  tvSeriesBody() {
    if (this.formRefs.title.current.value.trim() == "") {
      this.setState({errorMessage: "Series title is required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.formRefs.leadActor1.current.value.trim() == "") {
      this.setState({errorMessage: "Two lead actors are required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.formRefs.leadActor2.current.value.trim() == "") {
      this.setState({errorMessage: "Two lead actors are required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.formRefs.startYear.current.value.trim() == "") {
      this.setState({errorMessage: "Start year is required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.formRefs.endYear.current.value.trim() == "") {
      this.setState({errorMessage: "End year is required (can be 'Present')", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.state.imageInvalid) {
      this.setState({errorMessage: "Image is required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    let imageUrl = ""
    let imageFile = ""
    if (this.state.imageType == "url") {
      imageUrl = this.formRefs.imageUrl.current.value.trim()
    } else {
      imageFile = this.state.imageFile
    }

    return [{
      title: this.formRefs.title.current.value.trim(),
      lead_actors: [this.formRefs.leadActor1.current.value.trim(), this.formRefs.leadActor2.current.value.trim()],
      tv_series_start_year: this.formRefs.startYear.current.value.trim(),
      tv_series_end_year: this.formRefs.endYear.current.value.trim(),
      image_url: imageUrl,
      image_file: imageFile,
    }, null]
  }

  movieBody() {
    if (this.formRefs.title.current.value.trim() == "") {
      this.setState({errorMessage: "Movie title is required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.formRefs.leadActor1.current.value.trim() == "") {
      this.setState({errorMessage: "Two lead actors are required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.formRefs.leadActor2.current.value.trim() == "") {
      this.setState({errorMessage: "Two lead actors are required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.formRefs.publishedYear.current.value.trim() == "") {
      this.setState({errorMessage: "Release year is required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.state.imageInvalid) {
      this.setState({errorMessage: "Image is required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    let seriesTitle = ""
    let seriesNumber = ""
    if (this.formRefs.isSeries.current.checked) {
      seriesTitle = this.formRefs.seriesTitle.current.value.trim()
      seriesNumber = this.formRefs.seriesNumber.current.value.trim()

      if (seriesTitle == "") {
        this.setState({errorMessage: "Series title is required", success: false, submitting: false})
        window.scrollTo(0, 0)
        return [null, true]
      }
    }

    let imageUrl = ""
    let imageFile = ""
    if (this.state.imageType == "url") {
      imageUrl = this.formRefs.imageUrl.current.value.trim()
    } else {
      imageFile = this.state.imageFile
    }

    return [{
      title: this.formRefs.title.current.value.trim(),
      lead_actors: [this.formRefs.leadActor1.current.value.trim(), this.formRefs.leadActor2.current.value.trim()],
      published_year: this.formRefs.publishedYear.current.value.trim(),
      is_series: this.formRefs.isSeries.current.checked,
      series_title: seriesTitle,
      sequence_number: seriesNumber,
      image_url: imageUrl,
      image_file: imageFile,
    }, null]
  }

  bookBody() {
    if (this.formRefs.title.current.value.trim() == "") {
      this.setState({errorMessage: "Book title is required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.formRefs.author.current.value.trim() == "") {
      this.setState({errorMessage: "Author is required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.formRefs.publishedYear.current.value.trim() == "") {
      this.setState({errorMessage: "Published year is required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    if (this.state.imageInvalid) {
      this.setState({errorMessage: "Image is required", success: false, submitting: false})
      window.scrollTo(0, 0)
      return [null, true]
    }

    let seriesTitle = ""
    let seriesNumber = ""
    if (this.formRefs.isSeries.current.checked) {
      seriesTitle = this.formRefs.seriesTitle.current.value.trim()
      seriesNumber = this.formRefs.seriesNumber.current.value.trim()

      if (seriesTitle == "") {
        this.setState({errorMessage: "Series title is required", success: false, submitting: false})
        window.scrollTo(0, 0)
        return [null, true]
      }
    }

    let imageUrl = ""
    let imageFile = ""
    if (this.state.imageType == "url") {
      imageUrl = this.formRefs.imageUrl.current.value.trim()
    } else {
      imageFile = this.state.imageFile
    }

    return [{
      title: this.formRefs.title.current.value.trim(),
      author: this.formRefs.author.current.value.trim(),
      published_year: this.formRefs.publishedYear.current.value.trim(),
      is_series: this.formRefs.isSeries.current.checked,
      series_title: seriesTitle,
      sequence_number: seriesNumber,
      image_url: imageUrl,
      image_file: imageFile,
    }, null]
  }

  clearForm() {
    // Only clear the fields that wouldn't make sense to keep if submitting multiple items in a row
    // E.g. keep series title, author, lead actors, the same
    if (this.formRefs.title?.current) {
      this.formRefs.title.current.value = ""
    }

    if (this.formRefs.seriesNumber?.current) {
      this.formRefs.seriesNumber.current.value = ""
    }

    if (this.formRefs.publishedYear?.current) {
      this.formRefs.publishedYear.current.value = ""
    }

    if (this.formRefs.startYear?.current) {
      this.formRefs.startYear.current.value = ""
    }

    if (this.formRefs.endYear?.current) {
      this.formRefs.endYear.current.value = ""
    }

    if (this.formRefs.imageResult?.current) {
      this.formRefs.imageResult.current.src = ""
    }

    if (this.formRefs.imageUrl?.current) {
      this.formRefs.imageUrl.current.value = ""
    }

    if (this.formRefs.imageFile?.current) {
      this.formRefs.imageFile.current.value = null
    }

    if (this.formRefs.mediaType.current.value == "TV Series") {
      // Clear lead actors for TV series but not for movie - since movies could be in a series
      if (this.formRefs.leadActor1?.current) {
        this.formRefs.leadActor1.current.value = ""
      }

      if (this.formRefs.leadActor2?.current) {
        this.formRefs.leadActor2.current.value = ""
      }
    }
  }

  handleSeriesChange(e) {
    this.handleChange()
    this.setState({isSeries: e.target.checked})
  }

  handleTypeChange(e) {
    this.setState({type: e.target.value})
    this.clearForm()
    this.handleChange()
  }

  handleImageType(e) {
    this.setState({imageType: e.target.value})
    this.formRefs.imageResult.current.src = ""
  }

  handleImageFile(input) {
    this.setState({imageInvalid: false})
    if (input.target.files && input.target.files[0]) {
        var reader = new FileReader();
        reader.onload = (e) => {
          this.formRefs.imageResult.current.src = e.target.result
          this.setState({imageFile: input.target.files[0]})
          this.handleChange()
        };
        reader.readAsDataURL(input.target.files[0]);
    }
  }

  handleImageURL(input) {
    this.setState({imageInvalid: false})
    this.formRefs.imageResult.current.src = input.target.value
    this.handleChange()
  }

  onImageError(e) {
    this.setState({imageInvalid: true})
    if (this.formRefs.imageResult.current.src == window.location.href) {
      return // image is empty
    }
    this.formRefs.imageResult.current.src = invalidImage
    this.handleChange()
  }

  render() {
    return (
      <>
        <NavBar/>
        <Container className={this.props.className} fluid>
          <Row>
            <Col xs={12} style={{backgroundColor:"#222"}} className="pt-3 pb-3">
              <Form className="" noValidate validated={this.state.validate} onSubmit={null}>
              <Row className="pt-3 pb-3">
                <Col xs={0} lg={4} className="d-flex text-right align-items-center justify-content-end"></Col>
                <Col xs={12} lg={7} className="d-flex flex-column">
                  <h3 className="mb-4">Add New Item</h3>
                  {
                    function() {
                      if (this.state.success) {
                        return (
                          <Alert variant="success">
                            Item saved successfully. Do another! :O
                          </Alert>
                        )
                      }
                      if (this.state.errorMessage) {
                        return (
                          <Alert variant="danger">
                            {this.state.errorMessage}
                          </Alert>
                        )
                      }
                    }.bind(this)()
                  }
                  
                  
                    {/* <p style={{color:"red"}}>{this.state.errorMessage}</p> */}
                    <Form.Group controlId="newItem.mediaType" className="medium-input">
                      <Form.Label style={{color:"white"}}>Type</Form.Label>
                      <Form.Control as="select" onChange={this.handleTypeChange} ref={this.formRefs.mediaType}>
                        <option>TV Series</option>
                        <option>Movie</option>
                        <option>Book</option>
                      </Form.Control>
                    </Form.Group>
                    {
                      function () {
                        switch(this.state.type) {
                          case "TV Series":
                            return (
                              <>
                                <Form.Group controlId="newTvSeries.title" className="large-input text-left">
                                  <Form.Label style={{color:"white"}}>Series Title</Form.Label>
                                  <Form.Control required type="text" placeholder="What We Do in the Shadows" onChange={this.handleChange} ref={this.formRefs.title}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                                <Form.Group controlId="newTvSeries.leadActor1" className="medium-input text-left">
                                  <Form.Label style={{color:"white"}}>Lead Actor 1</Form.Label>
                                  <Form.Control required type="text" placeholder="Kayvan Novak" onChange={this.handleChange} ref={this.formRefs.leadActor1}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                                <Form.Group controlId="newTvSeries.leadActor2" className="medium-input text-left">
                                  <Form.Label style={{color:"white"}}>Lead Actor 2</Form.Label>
                                  <Form.Control required type="text" placeholder="Matt Berry" onChange={this.handleChange} ref={this.formRefs.leadActor2}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                                <Form.Group controlId="newTvSeries.startYear" className="medium-input text-left">
                                  <Form.Label style={{color:"white"}}>Start Year</Form.Label>
                                  <Form.Control required type="text" placeholder="2019" onChange={this.handleChange} ref={this.formRefs.startYear}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                                <Form.Group controlId="newTvSeries.endYear" className="medium-input text-left">
                                  <Form.Label style={{color:"white"}}>End Year</Form.Label>
                                  <Form.Control required type="text" placeholder="Present" onChange={this.handleChange} ref={this.formRefs.endYear}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                              </>
                            )
                          case "Movie":
                            return (
                              <>
                                <Form.Group controlId="newMovie.title" className="large-input text-left">
                                  <Form.Label style={{color:"white"}}>Movie Title</Form.Label>
                                  <Form.Control required type="text" placeholder="Terminator 2: Judgment Day" onChange={this.handleChange} ref={this.formRefs.title}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                                <Form.Group controlId="newMovie.isSeries" className="medium-input text-left">
                                  <Form.Check style={{color: "white"}} type="checkbox" label="Is Series?" onChange={this.handleSeriesChange} ref={this.formRefs.isSeries}/>
                                </Form.Group>
                                {
                                  function () {
                                    if (this.state.isSeries) {
                                      return (
                                        <>
                                          <Form.Group controlId="newMovie.seriesTitle" className="large-input text-left">
                                            <Form.Label style={{color:"white"}}>Series Title</Form.Label>
                                            <Form.Control required={this.state.isSeries} type="text" placeholder="Terminator" onChange={this.handleChange} ref={this.formRefs.seriesTitle}/>
                                            <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                          </Form.Group>
                                          <Form.Group controlId="newMovie.number" className="small-input text-left">
                                            <Form.Label style={{color:"white"}}>Number</Form.Label>
                                            <Form.Control type="text" placeholder="2" onChange={this.handleChange} ref={this.formRefs.seriesNumber}/>
                                          </Form.Group>
                                        </>
                                      )
                                    } else {
                                      return (<></>)
                                    }
                                  }.bind(this)()
                                }
                                <Form.Group controlId="newMovie.leadActor1" className="medium-input text-left">
                                  <Form.Label style={{color:"white"}}>Lead Actor 1</Form.Label>
                                  <Form.Control required type="text" placeholder="Arnold Schwarzenegger" onChange={this.handleChange} ref={this.formRefs.leadActor1}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                                <Form.Group controlId="newMovie.leadActor2" className="medium-input text-left">
                                  <Form.Label style={{color:"white"}}>Lead Actor 2</Form.Label>
                                  <Form.Control required type="text" placeholder="Linda Hamilton" onChange={this.handleChange} ref={this.formRefs.leadActor2}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                                <Form.Group controlId="newMovie.year" className="medium-input text-left">
                                  <Form.Label style={{color:"white"}}>Release Year</Form.Label>
                                  <Form.Control required type="text" placeholder="1991" onChange={this.handleChange} ref={this.formRefs.publishedYear}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                              </>
                            )
                          case "Book":
                            return (
                              <>
                                <Form.Group controlId="newBook.title" className="large-input text-left">
                                  <Form.Label style={{color:"white"}}>Book Title</Form.Label>
                                  <Form.Control required type="text" placeholder="The Fellowship of the Ring" onChange={this.handleChange} ref={this.formRefs.title}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                                <Form.Group controlId="newBook.isSeries" className="medium-input text-left">
                                  <Form.Check style={{color: "white"}} type="checkbox" label="Is Series?" onChange={this.handleSeriesChange} ref={this.formRefs.isSeries}/>
                                </Form.Group>
                                {
                                  function () {
                                    if (this.state.isSeries) {
                                      return (
                                        <>
                                          <Form.Group controlId="newBook.seriesTitle" className="large-input text-left">
                                            <Form.Label style={{color:"white"}}>Series Title</Form.Label>
                                            <Form.Control required={this.state.isSeries} type="text" placeholder="The Lord of the Rings" onChange={this.handleChange} ref={this.formRefs.seriesTitle}/>
                                            <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                          </Form.Group>
                                          <Form.Group controlId="newBook.number" className="small-input text-left">
                                            <Form.Label style={{color:"white"}}>Number</Form.Label>
                                            <Form.Control type="text" placeholder="1" onChange={this.handleChange} ref={this.formRefs.seriesNumber}/>
                                          </Form.Group>
                                        </>
                                      )
                                    } else {
                                      return (<></>)
                                    }
                                  }.bind(this)()
                                }
                                <Form.Group controlId="newBook.author" className="medium-input text-left">
                                  <Form.Label style={{color:"white"}}>Author</Form.Label>
                                  <Form.Control required type="text" placeholder="J. R. R. Tolkien" onChange={this.handleChange} ref={this.formRefs.author}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                                <Form.Group controlId="newBook.year" className="medium-input text-left">
                                  <Form.Label style={{color:"white"}}>Published Year</Form.Label>
                                  <Form.Control required type="text" placeholder="1954" onChange={this.handleChange} ref={this.formRefs.publishedYear}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                              </>
                            )
                          default:
                            console.error("unknown type: " + this.state.type)
                            return (<></>)
                        }
                      }.bind(this)()
                    }
                    <Form.Group controlId="newBook.imageType" className="text-left">
                      <Form.Label style={{color:"white"}}>Cover Image:&nbsp;&nbsp;&nbsp;&nbsp;</Form.Label>
                      <Form.Check name="imageRadio" defaultChecked inline type="radio" style={{color: "white"}} label="From URL" onClick={this.handleImageType} value="url"/>
                      <Form.Check name="imageRadio" inline type="radio" style={{color: "white"}} label="From file" onClick={this.handleImageType} value="file"/>
                    </Form.Group>
                    {
                      function () {
                        switch(this.state.imageType) {
                          case "url":
                            return (
                              <Form.Group controlId="newItem.imageUrl" className="large-input text-left">
                                <Form.Label style={{color:"white"}}>Image URL</Form.Label>
                                <Form.Control required={this.state.imageType == "url"} type="text" placeholder="https://www.example.com/image123.jpg" onChange={this.handleImageURL} ref={this.formRefs.imageUrl}/>
                                <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                              </Form.Group>
                            )
                          case "file":
                            return (
                              <>
                                <Form.Group controlId="newItem.imageFile">
                                  <Form.File required style={{color:"white"}} label="Image from file" onChange={this.handleImageFile} ref={this.formRefs.imageFile}/>
                                  <Form.Control.Feedback type="invalid">Required</Form.Control.Feedback>
                                </Form.Group>
                              </>
                            )
                          default:
                            console.error("unknown image type: " + this.state.imageType)
                            return (<></>)
                        }
                      }.bind(this)()
                    }
                    <div className="d-none image-area mt-4"><Image id="imageResult" src="" alt="" fluid className="rounded shadow-sm d-block" ref={this.formRefs.imageResult} onError={this.onImageError}/></div>
                </Col>
              </Row>
              <Row>
                <Col xs={0} lg={4} className="d-flex text-right align-items-center justify-content-end">
                </Col>
                <Col xs={12} lg={7}>
                  <h4 className="mt-3 mb-4">Preview</h4>
                </Col>
              </Row>
              <Row style={{backgroundColor: "#000"}}>
                <Col>
                  <SearchResult key={0} result={this.state.preview} isPreview/>
                </Col>
              </Row>
              <Row>
                <Col xs={0} lg={4} className="d-flex text-right align-items-center justify-content-end">
                </Col>
                <Col xs={12} lg={7}>
                  <Button variant="primary" className="orange-btn mt-4" type="submit" disabled={this.state.submitting} onClick={!this.state.submitting ? this.handleCreate : null}>
                    {this.state.submitting ? 'Loading…' : 'Create'}
                  </Button>
                  <Button variant="secondary" className="mt-4 ml-3" type="submit" disabled={this.state.submitting} onClick={!this.state.submitting ? this.handleCreateAndAdd : null}>
                    {this.state.submitting ? 'Loading…' : 'Create and Add To Shelf'}
                  </Button>
                </Col>
              </Row>
              </Form>
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}

export default connect(AddNewItem);

// TODO - Add in TV and Movie forms