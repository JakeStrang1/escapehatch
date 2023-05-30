package app

import (
	"github.com/JakeStrang1/escapehatch/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/****************************************************************************************
 * router.go
 *
 * This file is intended to:
 * - Define all API routes and middleware for the app
 ****************************************************************************************/

// DefaultAccessPolicy is to allow signed in users only and users cannot act on behalf of other users
var DefaultAccessPolicy = []gin.HandlerFunc{http.AccessPolicyAuthenticatedUsersOnly, http.AccessPolicyUsersCannotOverrideSelf}

// Router sets up the router for the app
func Router(config Config) *gin.Engine {
	r := http.NewEngine()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.CORSAllowOrigin},
		AllowMethods:     []string{"OPTIONS", "POST", "GET", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
	}))

	/*******************************************
	 * Public routes
	 *******************************************/

	r.POST("/auth/sign-up", http.SignUp)
	r.POST("/auth/sign-in", http.SignIn)
	r.POST("/auth/verify", http.Verify)
	r.POST("/auth/not-you", http.NotYou)

	r.Use(http.Authenticate) // Identify caller if possible

	/*******************************************
	 * Access-controlled routes past this point
	 *******************************************/

	// Users
	r.GET("/users", DefaultAccessPolicy, http.GetUsers)
	r.GET("/users/:id", DefaultAccessPolicy, http.AccessPolicyUsersCannotOverrideID, http.GetUser)
	r.PATCH("/users/:id", DefaultAccessPolicy, http.AccessPolicyUsersCannotOverrideID, http.UpdateUser)
	r.GET("/users/:id/followers", DefaultAccessPolicy, http.GetUserFollowers)

	// Posts
	r.POST("/posts", DefaultAccessPolicy, http.CreatePost)
	r.GET("/posts", DefaultAccessPolicy, http.GetPosts)
	r.GET("/posts/:id", DefaultAccessPolicy, http.AccessPolicyPost, http.GetPost)
	r.PATCH("/posts/:id", DefaultAccessPolicy, http.AccessPolicyPost, http.UpdatePost)
	r.DELETE("/posts/:id", DefaultAccessPolicy, http.AccessPolicyPost, http.DeletePost)

	// Comments
	r.POST("/posts/:id/comments", DefaultAccessPolicy, http.CreateComment)
	r.GET("/posts/:id/comments", DefaultAccessPolicy, http.GetComments)
	r.GET("/posts/:id/comments/:commentID", DefaultAccessPolicy, http.GetComment)
	r.PATCH("/posts/:id/comments/:commentID", DefaultAccessPolicy, http.UpdateComment)
	r.DELETE("/posts/:id/comments/:commentID", DefaultAccessPolicy, http.DeleteComment)

	return r.Engine
}
