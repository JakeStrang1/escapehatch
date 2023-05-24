package http

import (
	"net/http"
	"os"
	"time"

	"github.com/JakeStrang1/escapehatch/services/auth"
	"github.com/JakeStrang1/escapehatch/services/auth/session"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// SignInBody is used for both Sign In and Sign Up
type SignInBody struct {
	Email      *string           `json:"email" validate:"required,email"`
	Metadata   map[string]string `json:"metadata"`
	RememberMe *bool             `json:"remember_me"`
}

func SignUp(c *gin.Context) {
	var body SignInBody
	err := Body(c, &body)
	if err != nil {
		return
	}

	err = auth.SignUp(lo.FromPtr(body.Email), body.Metadata, c.ClientIP())
	if err != nil {
		Error(c, err)
		return
	}
	c.JSON(200, struct{}{})
}

func SignIn(c *gin.Context) {
	var body SignInBody
	err := Body(c, &body)
	if err != nil {
		return
	}

	err = auth.SignIn(lo.FromPtr(body.Email), body.Metadata, lo.FromPtr(body.RememberMe), c.ClientIP())
	if err != nil {
		Error(c, err)
		return
	}
	c.JSON(200, struct{}{})
}

type VerifyBody struct {
	Email     *string `json:"email" validate:"required_without=EmailHash"`
	Secret    *string `json:"secret" validate:"required"`
	EmailHash *string `json:"email_hash" validate:"required_without=Email"`
}

type VerifyResponse struct {
	Metadata map[string]string `json:"metadata"`
}

func Verify(c *gin.Context) {
	var body VerifyBody
	err := Body(c, &body)
	if err != nil {
		return
	}

	session, metadata, err := auth.Verify(lo.FromPtr(body.Email), lo.FromPtr(body.Secret), lo.FromPtr(body.EmailHash))
	if err != nil {
		Error(c, err)
		return
	}

	response := VerifyResponse{Metadata: metadata}

	secure := false
	httpOnly := false
	if os.Getenv("PRODUCTION") != "false" {
		// Cookie doesn't save in production without this. Must be called before SetCookie().
		// Secure flag must also be set when using this, so only use on production
		c.SetSameSite(http.SameSiteNoneMode)
		secure = true
		httpOnly = true
	}
	c.SetCookie("SID", session.Token, maxAge(*session), "/", host(c.Request.Host), secure, httpOnly)

	c.JSON(200, response)
}

// maxAge returns 0 if session.RememberMe=false, otherwise returns the number of seconds between Now and ExpiresAt
func maxAge(session session.Session) int {
	now := time.Now()
	if !session.RememberMe || session.ExpiresAt.Before(now) {
		return 0
	}
	return int(session.ExpiresAt.Sub(now).Seconds())
}

type NotYouBody struct {
	Secret       *string `json:"secret" validate:"required"`
	EmailHash    *string `json:"email_hash" validate:"required"`
	DoNotContact *bool   `json:"do_not_contact"`
}

func NotYou(c *gin.Context) {
	var body NotYouBody
	err := Body(c, &body)
	if err != nil {
		return
	}

	err = auth.NotYou(lo.FromPtr(body.Secret), lo.FromPtr(body.EmailHash), lo.FromPtr(body.DoNotContact))
	if err != nil {
		Error(c, err)
		return
	}

	c.JSON(200, struct{}{})
}

func host(host string) string {
	if os.Getenv("PRODUCTION") != "false" {
		return host
	}
	return "localhost"
}
