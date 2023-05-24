package test

import (
	api "github.com/JakeStrang1/escapehatch/http"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
)

// TestCreateUser tests user creation
func (s *Suite) TestCreateUser() {
	// Sign up
	_, _ = s.CreateUser(s.NewSeedEmail("user1"))

	// Test duplicate email fails
	signUpBody := api.SignInBody{
		Email: lo.ToPtr(s.NewSeedEmail("user1")),
	}
	response := s.Post("/auth/sign-up", signUpBody)
	s.Assert().Equal(422, response.Status)
	s.Assert().Equal("email_unavailable", gjson.Get(response.Body, "code").String())
}
