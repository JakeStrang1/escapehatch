package http

import (
	"github.com/JakeStrang1/escapehatch/internal"
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/services/posts"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type PostAPI struct {
	DefaultModelAPI
	OwnerID  *string      `json:"owner_id"`
	Body     *string      `json:"body" validate:"required"`
	Comments []CommentAPI `json:"comments"`
}

func CreatePost(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)

	body := PostAPI{}
	err := Body(c, &body)
	if err != nil {
		Error(c, err)
		return
	}

	result := ToPost(body)
	err = posts.Create(userID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToPostAPI(result))
}

func GetPosts(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)

	query := PageQuery{}
	err := c.BindQuery(&query)
	if err != nil {
		Error(c, &errors.Error{Code: errors.BadRequest, Message: "error in query parameters", Err: err})
		return
	}

	results := []posts.Post{}
	filter := posts.GetManyFilter{
		OwnerID: &userID,
		Page:    query.Page,
		PerPage: query.PerPage,
	}
	pageInfo, err := posts.GetPage(filter, &results)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnMany(c, internal.Map(results, ToPostAPI), *pageInfo)
}

func GetPost(c *gin.Context) {
	id := c.Param("id")

	result := posts.Post{}
	err := posts.GetByID(id, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToPostAPI(result))
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")

	body := PostAPI{}
	err := Body(c, &body)
	if err != nil {
		Error(c, err)
		return
	}

	result := posts.Post{}
	update := ToPostUpdate(body)
	err = posts.Update(id, update, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToPostAPI(result))
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")

	err := posts.Delete(id)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, struct{}{})
}

func ToPost(postAPI PostAPI) posts.Post {
	return posts.Post{
		DefaultModel: FromDefaultModelAPI(postAPI.DefaultModelAPI),
		OwnerID:      lo.FromPtr(postAPI.OwnerID),
		Body:         lo.FromPtr(postAPI.Body),
	}
}

func ToPostAPI(post posts.Post) PostAPI {
	return PostAPI{
		DefaultModelAPI: ToDefaultModelAPI(post.DefaultModel),
		OwnerID:         &post.OwnerID,
		Body:            &post.Body,
		Comments:        internal.Map(post.Comments, ToCommentAPI),
	}
}

func ToPostUpdate(postAPI PostAPI) posts.PostUpdate {
	return posts.PostUpdate{
		Body: postAPI.Body,
	}
}

func AccessPolicyPost(c *gin.Context) {
	// Allow system/super user
	if GetUserID(c) == "" {
		c.Next()
		return
	}

	// Get post
	result := posts.Post{}
	err := posts.GetByID(c.Param("id"), &result)
	if err != nil {
		Error(c, err)
		return
	}

	// Check owner
	if !result.BelongsTo(GetUserID(c)) {
		Error(c, errors.NewNotFound())
		return
	}

	c.Next()
}
