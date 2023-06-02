package http

import (
	"github.com/JakeStrang1/escapehatch/internal"
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/internal/pages"
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
	ItemID      *string `json:"item_id"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image"`
}

type FollowerAPI struct {
	DefaultModelAPI  `json:",inline"`
	TargetUserID     *string `json:"target_user_id"`
	TargetUsername   *string `json:"target_username"`
	TargetFullName   *string `json:"target_full_name"`
	FollowerUserID   *string `json:"follower_user_id"`
	FollowerUsername *string `json:"follower_username"`
	FollowerFullName *string `json:"follower_full_name"`
}

type UserQuery struct {
	PageQuery
	Search *string `form:"search"`
}

type UserFollowerQuery struct {
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

func FollowUser(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(CtxKeyUserID)

	result := users.User{}
	err := users.Follow(id, userID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToUserAPI(userID, result))
}

func UnfollowUser(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(CtxKeyUserID)

	result := users.User{}
	err := users.Unfollow(id, userID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToUserAPI(userID, result))
}

func RemoveUser(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(CtxKeyUserID)

	result := users.User{}
	err := users.Remove(id, userID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToUserAPI(userID, result))
}

func GetUserFollowers(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(CtxKeyUserID)

	// Special case
	if id == "me" {
		id = userID
	}

	query := UserFollowerQuery{}
	err := c.BindQuery(&query)
	if err != nil {
		Error(c, &errors.Error{Code: errors.BadRequest, Message: "error in query parameters", Err: err})
		return
	}

	results := []users.Follower{}
	filter := users.FollowerFilter{
		TargetUserID: lo.ToPtr(id),
		Search:       query.Search,
	}
	err = users.GetManyFollowers(filter, &results)
	if err != nil {
		Error(c, err)
		return
	}

	pageInfo := pages.PageResult{
		Page:       1,
		PerPage:    len(results),
		TotalPages: 1,
	}

	ReturnMany(c, internal.Map(results, ToFollowerAPI), pageInfo)
}

func GetUserFollowing(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(CtxKeyUserID)

	// Special case
	if id == "me" {
		id = userID
	}

	query := UserFollowerQuery{}
	err := c.BindQuery(&query)
	if err != nil {
		Error(c, &errors.Error{Code: errors.BadRequest, Message: "error in query parameters", Err: err})
		return
	}

	results := []users.Follower{}
	filter := users.FollowerFilter{
		FollowerUserID: lo.ToPtr(id),
		Search:         query.Search,
	}
	err = users.GetManyFollowing(filter, &results)
	if err != nil {
		Error(c, err)
		return
	}

	pageInfo := pages.PageResult{
		Page:       1,
		PerPage:    len(results),
		TotalPages: 1,
	}

	ReturnMany(c, internal.Map(results, ToFollowerAPI), pageInfo)
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
		ItemID:      &item.ItemID,
		Description: &item.Description,
		ImageURL:    IncludeStaticRoot(item.ImageURL),
	}
}

func ToFollowerAPI(follower users.Follower) FollowerAPI {
	return FollowerAPI{
		DefaultModelAPI:  ToDefaultModelAPI(follower.DefaultModel),
		TargetUserID:     &follower.TargetUserID,
		TargetUsername:   &follower.TargetUsername,
		TargetFullName:   &follower.TargetFullName,
		FollowerUserID:   &follower.FollowerUserID,
		FollowerUsername: &follower.FollowerUsername,
		FollowerFullName: &follower.FollowerFullName,
	}
}
