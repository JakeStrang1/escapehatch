package http

import (
	"github.com/JakeStrang1/saas-template/internal"
	"github.com/JakeStrang1/saas-template/internal/pages"
	"github.com/JakeStrang1/saas-template/services/posts"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type CommentAPI struct {
	DefaultModelAPI
	OwnerID *string `json:"owner_id"`
	PostID  *string `json:"post_id"`
	Body    *string `json:"body" validate:"required"`
	ReplyTo *string `json:"reply_to"`
}

func CreateComment(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)
	postID := c.Param("id")

	body := CommentAPI{}
	err := Body(c, &body)
	if err != nil {
		Error(c, err)
		return
	}

	result := ToComment(body)
	err = posts.CreateComment(userID, postID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToCommentAPI(result))
}

func GetComments(c *gin.Context) {
	postID := c.Param("id")

	results := []posts.Comment{}
	err := posts.GetComments(postID, &results)
	if err != nil {
		Error(c, err)
		return
	}

	pageInfo := pages.PageResult{
		Page:       1,
		PerPage:    len(results),
		TotalPages: 1,
	}

	ReturnMany(c, internal.Map(results, ToCommentAPI), pageInfo)
}

func GetComment(c *gin.Context) {
	postID := c.Param("id")
	commentID := c.Param("commentID")

	result := posts.Comment{}
	err := posts.GetCommentByID(postID, commentID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToCommentAPI(result))
}

func UpdateComment(c *gin.Context) {
	postID := c.Param("id")
	commentID := c.Param("commentID")

	body := CommentAPI{}
	err := Body(c, &body)
	if err != nil {
		Error(c, err)
		return
	}

	result := posts.Comment{}
	update := ToCommentUpdate(body)
	err = posts.UpdateComment(postID, commentID, update, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToCommentAPI(result))
}

func DeleteComment(c *gin.Context) {
	postID := c.Param("id")
	commentID := c.Param("commentID")

	err := posts.DeleteComment(postID, commentID)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, struct{}{})
}

func ToComment(commentAPI CommentAPI) posts.Comment {
	return posts.Comment{
		DefaultModel: FromDefaultModelAPI(commentAPI.DefaultModelAPI),
		OwnerID:      lo.FromPtr(commentAPI.OwnerID),
		Body:         lo.FromPtr(commentAPI.Body),
	}
}

func ToCommentAPI(comment posts.Comment) CommentAPI {
	return CommentAPI{
		DefaultModelAPI: ToDefaultModelAPI(comment.DefaultModel),
		OwnerID:         &comment.OwnerID,
		PostID:          &comment.PostID,
		Body:            &comment.Body,
		ReplyTo:         &comment.ReplyTo,
	}
}

func ToCommentUpdate(commentAPI CommentAPI) posts.CommentUpdate {
	return posts.CommentUpdate{
		Body: commentAPI.Body,
	}
}
