package test

import (
	"fmt"

	"github.com/JakeStrang1/escapehatch/http"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
)

// TestUserJourney is meant to be a test that runs through the happy path of all the major features.
// It can also be used to seed a database with test data.
func (s *Suite) TestUserJourney() {
	// Sign up user 1
	_, withUser1Cookie := s.CreateUser(s.NewSeedEmail("user1"))

	// Sign up user 2
	_, withUser2Cookie := s.CreateUser(s.NewSeedEmail("user2"))

	// Create post
	post := http.PostAPI{
		Body: lo.ToPtr("this is the body"),
	}
	response := s.Post("/posts", post, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("this is the body", gjson.Get(response.Body, "data.body").String())
	postID := gjson.Get(response.Body, "data.id").String()
	s.Assert().True(postID != "")

	// Get posts
	response = s.Get("/posts", withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("this is the body", gjson.Get(response.Body, "data.0.body").String())
	s.Assert().Equal(1, int(gjson.Get(response.Body, "pages.page").Int())) // Quick pagination assertions, move this to dedicated section eventually
	s.Assert().Equal(25, int(gjson.Get(response.Body, "pages.per_page").Int()))
	s.Assert().Equal(1, int(gjson.Get(response.Body, "pages.total_pages").Int()))

	// Get post by ID
	response = s.Get(fmt.Sprintf("/posts/%s", postID), withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("this is the body", gjson.Get(response.Body, "data.body").String())
	s.Assert().False(gjson.Get(response.Body, "pages").Exists()) // Another page test

	// Update post by ID
	post = http.PostAPI{
		Body: lo.ToPtr("this is the new body"),
	}
	response = s.Patch(fmt.Sprintf("/posts/%s", postID), post, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("this is the new body", gjson.Get(response.Body, "data.body").String())

	// Create comment
	comment := http.CommentAPI{
		Body: lo.ToPtr("this post is hilarious"),
	}
	response = s.Post(fmt.Sprintf("/posts/%s/comments", postID), comment, withUser2Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("this post is hilarious", gjson.Get(response.Body, "data.body").String())
	s.Assert().Equal(postID, gjson.Get(response.Body, "data.post_id").String())
	commentID := gjson.Get(response.Body, "data.id").String()
	s.Assert().True(commentID != "")

	// Get comments
	response = s.Get(fmt.Sprintf("/posts/%s/comments", postID), withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("this post is hilarious", gjson.Get(response.Body, "data.0.body").String())
	s.Assert().Equal(1, int(gjson.Get(response.Body, "pages.page").Int())) // Dummy pagination values since comments aren't actually paginated
	s.Assert().Equal(1, int(gjson.Get(response.Body, "pages.per_page").Int()))
	s.Assert().Equal(1, int(gjson.Get(response.Body, "pages.total_pages").Int()))

	// Get comment by ID
	response = s.Get(fmt.Sprintf("/posts/%s/comments/%s", postID, commentID), withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("this post is hilarious", gjson.Get(response.Body, "data.body").String())

	// Update comment by ID
	comment = http.CommentAPI{
		Body: lo.ToPtr("this post is NOT hilarious"),
	}
	response = s.Patch(fmt.Sprintf("/posts/%s/comments/%s", postID, commentID), comment, withUser2Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("this post is NOT hilarious", gjson.Get(response.Body, "data.body").String())

	// Get post with comments
	response = s.Get(fmt.Sprintf("/posts/%s", postID), withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("this post is NOT hilarious", gjson.Get(response.Body, "data.comments.0.body").String())

	// Delete post by ID
	response = s.Delete(fmt.Sprintf("/posts/%s", postID), withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("{}", gjson.Get(response.Body, "data").Raw)

	// Get user1 self
	response = s.Get("/users/me", withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(1, int(gjson.Get(response.Body, "data.number").Int()))
	userID1 := gjson.Get(response.Body, "data.id").String()
	s.Assert().True(userID1 != "")

	// Get user2 self
	response = s.Get("/users/me", withUser2Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(2, int(gjson.Get(response.Body, "data.number").Int()))
	userID2 := gjson.Get(response.Body, "data.id").String()
	s.Assert().True(userID2 != "")

	// Update user1
	userBody := http.UserAPI{
		Username: lo.ToPtr("stealth.dragon"),
		FullName: lo.ToPtr("John L. Userman"),
	}
	response = s.Patch("/users/me", userBody, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("stealth.dragon", gjson.Get(response.Body, "data.username").String())
	s.Assert().Equal("John L. Userman", gjson.Get(response.Body, "data.full_name").String())

	// Update user2
	userBody = http.UserAPI{
		Username: lo.ToPtr("crouching.sock"),
		FullName: lo.ToPtr("Amy del Taco von Trapp"),
	}
	response = s.Patch("/users/me", userBody, withUser2Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("crouching.sock", gjson.Get(response.Body, "data.username").String())
	s.Assert().Equal("Amy del Taco von Trapp", gjson.Get(response.Body, "data.full_name").String())

	// Get users
	response = s.Get("/users", withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(2, int(gjson.Get(response.Body, "data.#").Int()))
	s.Assert().True(gjson.Get(response.Body, "data.0.self").Bool())
	s.Assert().False(gjson.Get(response.Body, "data.1.self").Bool())

	// Follow user2
	response = s.Post(fmt.Sprintf("/users/%s/follow", userID2), nil, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().True(gjson.Get(response.Body, "data.followed_by_you").Bool())
	s.Assert().False(gjson.Get(response.Body, "data.follows_you").Bool())

	// Follow user1
	response = s.Post(fmt.Sprintf("/users/%s/follow", userID1), nil, withUser2Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().True(gjson.Get(response.Body, "data.followed_by_you").Bool())
	s.Assert().True(gjson.Get(response.Body, "data.follows_you").Bool())

	// Get followers
	response = s.Get("/users/me/followers", withUser2Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(1, int(gjson.Get(response.Body, "data.#").Int()))
	s.Assert().Equal("stealth.dragon", gjson.Get(response.Body, "data.0.follower_username").String())

	// Get following
	response = s.Get("/users/me/following", withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(1, int(gjson.Get(response.Body, "data.#").Int()))
	s.Assert().Equal("crouching.sock", gjson.Get(response.Body, "data.0.target_username").String())

	// Unfollow user2
	response = s.Post(fmt.Sprintf("/users/%s/unfollow", userID2), nil, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().False(gjson.Get(response.Body, "data.followed_by_you").Bool())
	s.Assert().True(gjson.Get(response.Body, "data.follows_you").Bool())

	// Remove user2
	response = s.Post(fmt.Sprintf("/users/%s/remove", userID2), nil, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().False(gjson.Get(response.Body, "data.followed_by_you").Bool())
	s.Assert().False(gjson.Get(response.Body, "data.follows_you").Bool())

	// Create book
	book := http.BookAPI{
		ItemAPI: http.ItemAPI{
			ImageURL: lo.ToPtr("https://images-na.ssl-images-amazon.com/images/S/compressed.photo.goodreads.com/books/1474169725i/15881.jpg"),
		},
		Title:          lo.ToPtr("Harry Potter and the Chamber of Secrets"),
		Author:         lo.ToPtr("J. K. Rowling"),
		PublishedYear:  lo.ToPtr("1998"),
		IsSeries:       lo.ToPtr(true),
		SeriesTitle:    lo.ToPtr("Harry Potter"),
		SequenceNumber: lo.ToPtr("2"),
	}
	response = s.Post("/books", book, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("Harry Potter and the Chamber of Secrets (Harry Potter #2)", gjson.Get(response.Body, "data.description").String())
	bookID := gjson.Get(response.Body, "data.id").String()
	s.Assert().True(bookID != "")

	// Create movie
	movie := http.MovieAPI{
		ItemAPI: http.ItemAPI{
			ImageURL: lo.ToPtr("https://media1.inlander.com/inlander/imager/u/slideshow/21189517/the-lord-of-the-rings-the-fellowship-of-the-ring-2001-4k-remaster"),
		},
		Title:          lo.ToPtr("The Fellowship of the Ring"),
		LeadActors:     []string{"Elijah Wood", "Ian McKellen"},
		PublishedYear:  lo.ToPtr("2001"),
		IsSeries:       lo.ToPtr(true),
		SeriesTitle:    lo.ToPtr("The Lord of the Rings"),
		SequenceNumber: lo.ToPtr("1"),
	}
	response = s.Post("/movies", movie, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("The Fellowship of the Ring (The Lord of the Rings #1)", gjson.Get(response.Body, "data.description").String())
	movieID := gjson.Get(response.Body, "data.id").String()
	s.Assert().True(movieID != "")

	// Create tv series
	tvSeries := http.TVSeriesAPI{
		ItemAPI: http.ItemAPI{
			ImageURL: lo.ToPtr("https://i.ebayimg.com/images/g/MagAAMXQGQRR82PV/s-l500.jpg"),
		},
		Title:             lo.ToPtr("The Office"),
		LeadActors:        []string{"Steve Carell", "Jenna Fischer"},
		TVSeriesStartYear: lo.ToPtr("2005"),
		TVSeriesEndYear:   lo.ToPtr("2013"),
	}
	response = s.Post("/tv-series", tvSeries, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("The Office (2005 - 2013)", gjson.Get(response.Body, "data.description").String())
	tvSeriesID := gjson.Get(response.Body, "data.id").String()
	s.Assert().True(tvSeriesID != "")

	// Add book to shelf
	response = s.Post(fmt.Sprintf("/items/%s/add", bookID), nil, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(1, int(gjson.Get(response.Body, "data.user_count").Int()))

	// Add movie to shelf
	response = s.Post(fmt.Sprintf("/items/%s/add", movieID), nil, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(1, int(gjson.Get(response.Body, "data.user_count").Int()))

	// Add tv series to shelf
	response = s.Post(fmt.Sprintf("/items/%s/add", tvSeriesID), nil, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(1, int(gjson.Get(response.Body, "data.user_count").Int()))

	// Get user1 self - confirm shelf
	response = s.Get("/users/me", withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(3, int(gjson.Get(response.Body, "data.shelves.#").Int()))
	s.Assert().Equal("Harry Potter and the Chamber of Secrets (Harry Potter #2)", gjson.Get(response.Body, "data.shelves.0.items.0.description").String())
	s.Assert().Equal("The Fellowship of the Ring (The Lord of the Rings #1)", gjson.Get(response.Body, "data.shelves.1.items.0.description").String())
	s.Assert().Equal("The Office (2005 - 2013)", gjson.Get(response.Body, "data.shelves.2.items.0.description").String())

	// Remove from shelf
	response = s.Post(fmt.Sprintf("/items/%s/remove", bookID), nil, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(0, int(gjson.Get(response.Body, "data.user_count").Int()))

	// Get user1 self - confirm shelf
	response = s.Get("/users/me", withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(0, int(gjson.Get(response.Body, "data.shelves.0.items.#").Int()))
}
