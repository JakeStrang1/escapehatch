package http

import (
	"github.com/JakeStrang1/escapehatch/internal"
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/services/users"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type UserAPI struct {
	DefaultModelAPI `json:",inline"`
	Email           *string    `json:"email"`
	Username        *string    `json:"username"`
	Number          *int       `json:"number"`
	FullName        *string    `json:"full_name"`
	Shelves         []ShelfAPI `json:"shelves"`
	FollowerCount   *int       `json:"follower_count"`
	FollowingCount  *int       `json:"following_count"`
}

type ShelfAPI struct {
	DefaultModelAPI `json:",inline"`
	Name            *string        `json:"name"`
	Items           []ShelfItemAPI `json:"items"`
}

type ShelfItemAPI struct {
	ItemID *string `json:"item_id"`
	Image  *string `json:"image"`
}

type UserQuery struct {
	PageQuery
	Search *string `form:"search"`
}

func GetUsers(c *gin.Context) {
	query := UserQuery{}
	err := c.BindQuery(&query)
	if err != nil {
		Error(c, &errors.Error{Code: errors.BadRequest, Message: "error in query parameters", Err: err})
		return
	}

	results := []users.User{}
	filter := users.Filter{
		Page:    query.Page,
		PerPage: query.PerPage,
		Search:  query.Search,
	}

	pageInfo, err := users.GetPage(filter, &results)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnMany(c, internal.Map(results, ToUserAPI), *pageInfo)
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
		DefaultModelAPI: ToDefaultModelAPI(user.DefaultModel),
		Email:           &user.Email,
		Username:        &user.Username,
		Number:          &user.Number,
		FullName:        &user.FullName,
		Shelves:         internal.Map(user.Shelves, ToShelfAPI),
		FollowerCount:   lo.ToPtr(len(user.Followers)),
		FollowingCount:  lo.ToPtr(len(user.Following)),
	}
}

func ToShelfAPI(shelf users.Shelf) ShelfAPI {
	return ShelfAPI{
		DefaultModelAPI: ToDefaultModelAPI(shelf.DefaultModel),
		Name:            &shelf.Name,
		Items:           internal.Map(shelf.Items, ToShelfItemAPI),
	}
}

func ToShelfItemAPI(item users.ShelfItem) ShelfItemAPI {
	return ShelfItemAPI{
		ItemID: &item.ItemID,
		Image:  &item.Image,
	}
}
