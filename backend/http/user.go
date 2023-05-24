package http

import (
	"github.com/JakeStrang1/saas-template/services/users"
	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	ID    *string `json:"id"`
	Email *string `json:"email"`
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(CtxKeyUserID)

	// Special case
	if id == "me" {
		id = userID
	}

	u := users.User{}
	err := users.GetByID(id, &u)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToUserAPI(u))
}

func ToUserAPI(user users.User) UserAPI {
	return UserAPI{
		ID:    &user.ID,
		Email: &user.Email,
	}
}
