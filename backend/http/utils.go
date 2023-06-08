package http

import (
	"github.com/JakeStrang1/escapehatch/services/users"
	"github.com/gin-gonic/gin"
)

func ValidateUsername(c *gin.Context) {
	username := c.Query("u")
	err := users.ValidateUsername(username)
	if err != nil {
		Error(c, err)
		return
	}
	ReturnOne(c, struct{}{})
}
