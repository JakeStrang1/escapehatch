package http

import (
	"github.com/JakeStrang1/escapehatch/services/users"
	"github.com/gin-gonic/gin"
)

type HealthCheckAPI struct {
	Success bool `json:"success"`
}

// Public endpoint
func HealthCheck(c *gin.Context) {

	// Verify database connection is working
	_, err := users.GetCount(users.Filter{})
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, HealthCheckAPI{Success: true})
}
