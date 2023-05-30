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

	// Update user
	userBody := http.UserAPI{
		Username: lo.ToPtr("stealth.dragon"),
		FullName: lo.ToPtr("John L. Userman"),
	}
	response = s.Patch("/users/me", userBody, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal("stealth.dragon", gjson.Get(response.Body, "data.username").String())
	s.Assert().Equal("John L. Userman", gjson.Get(response.Body, "data.full_name").String())

	// Get users
	response = s.Get("/users", withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(2, int(gjson.Get(response.Body, "data.#").Int()))
	s.Assert().True(gjson.Get(response.Body, "data.0.self").Bool())
	s.Assert().False(gjson.Get(response.Body, "data.1.self").Bool())

	// Follow user
	response = s.Post(fmt.Sprintf("/users/%s/follow", userID2), nil, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().True(gjson.Get(response.Body, "data.followed_by_you").Bool())

	// Get followers
	response = s.Get("/users/me/followers", withUser2Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().Equal(1, int(gjson.Get(response.Body, "data.#").Int()))
	s.Assert().Equal("stealth.dragon", gjson.Get(response.Body, "data.0.follower_username").String())

	// Unfollow user
	response = s.Post(fmt.Sprintf("/users/%s/unfollow", userID2), nil, withUser1Cookie)
	s.Assert().Equal(200, response.Status)
	s.Assert().False(gjson.Get(response.Body, "data.followed_by_you").Bool())
}
