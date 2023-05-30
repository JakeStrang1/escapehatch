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
	Self            *bool      `json:"self"`
	FollowsYou      *bool      `json:"follows_you"`
	FollowedByYou   *bool      `json:"followed_by_you"`
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
	userID := c.GetString(CtxKeyUserID)
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

	ReturnMany(c, ToUserAPIs(userID, results), *pageInfo)
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

	ReturnOne(c, ToUserAPI(userID, u))
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(CtxKeyUserID)

	// Special case
	if id == "me" {
		id = userID
	}

	body := UserAPI{}
	err := Body(c, &body)
	if err != nil {
		Error(c, err)
		return
	}

	result := users.User{}
	update := ToUserUpdate(body)
	err = users.Update(id, update, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToUserAPI(userID, result))
}

func ToUserAPIs(selfID string, dbUsers []users.User) []UserAPI {
	results := []UserAPI{}
	for _, dbUser := range dbUsers {
		results = append(results, ToUserAPI(selfID, dbUser))
	}
	return results
}

func ToUserAPI(selfID string, user users.User) UserAPI {
	return UserAPI{
		DefaultModelAPI: ToDefaultModelAPI(user.DefaultModel),
		Email:           &user.Email,
		Username:        &user.Username,
		Number:          &user.Number,
		FullName:        &user.FullName,
		Shelves:         internal.Map(user.Shelves, ToShelfAPI),
		FollowerCount:   lo.ToPtr(len(user.Followers)),
		FollowingCount:  lo.ToPtr(len(user.Following)),
		Self:            lo.ToPtr(user.ID == selfID),
		FollowsYou:      lo.ToPtr(user.Follows(selfID)),
		FollowedByYou:   lo.ToPtr(user.FollowedBy(selfID)),
	}
}

func ToUserUpdate(userAPI UserAPI) users.UserUpdate {
	return users.UserUpdate{
		Username: userAPI.Username,
		FullName: userAPI.FullName,
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
