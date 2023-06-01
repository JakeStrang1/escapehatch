package http

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/internal/pages"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

var StaticURLRoot string

type Response struct {
	Data  interface{} `json:"data"`
	Pages *Pages      `json:"pages,omitempty"`
}

type Pages struct {
	Next       string `json:"next"`
	Previous   string `json:"previous"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	TotalPages int    `json:"total_pages"`
}

type PageQuery struct {
	Page    *int `form:"page"`
	PerPage *int `form:"per_page"`
}

func GetUserID(c *gin.Context) string {
	return c.GetString(CtxKeyUserID)
}

func Body(c *gin.Context, result interface{}) error {
	err := c.ShouldBindJSON(result)
	if err != nil {
		BindError(c, err)
		return err
	}

	v := validator.New()
	err = v.Struct(result)
	if err != nil {
		ValidationError(c, err)
		return err
	}

	return nil
}

func Error(c *gin.Context, err error) {
	fmt.Println(errors.Details(err))
	c.AbortWithStatusJSON(errors.Status(err), errors.Ensure(err))
}

func BindError(c *gin.Context, err error) {
	e := &errors.Error{Code: errors.BadRequest, Message: err.Error(), Err: err}
	fmt.Println(errors.Details(e))
	Error(c, e)
}

func ValidationError(c *gin.Context, err error) {
	e := &errors.Error{Code: errors.Invalid, Message: err.Error(), Err: err}
	fmt.Println(errors.Details(e))
	Error(c, e)
}

func ReturnOne(c *gin.Context, result interface{}) {
	response := Response{Data: result}
	c.JSON(200, response)
}

func ReturnMany(c *gin.Context, results interface{}, pageInfo pages.PageResult) {
	response := Response{
		Data: results,
		Pages: &Pages{
			Next:       next(*c.Request.URL, pageInfo),
			Previous:   previous(*c.Request.URL, pageInfo),
			Page:       pageInfo.Page,
			PerPage:    pageInfo.PerPage,
			TotalPages: pageInfo.TotalPages,
		},
	}
	c.JSON(200, response)
}

// IncludeStaticRoot prepends the static url to the given path to form a fully qualified URL
// Returns a pointer for convenience, is never nil
func IncludeStaticRoot(urlpath string) *string {
	s, err := url.JoinPath(StaticURLRoot, urlpath)
	if err != nil {
		panic(err) // error can only occur if base is malformed
	}
	return &s
}

func next(u url.URL, pageInfo pages.PageResult) string {
	if !pageInfo.HasMore() {
		return ""
	}
	q := u.Query()
	q.Set("page", strconv.Itoa(pageInfo.Page+1))
	q.Set("per_page", strconv.Itoa(pageInfo.PerPage))
	u.RawQuery = q.Encode()
	return u.RequestURI()
}

func previous(u url.URL, pageInfo pages.PageResult) string {
	if pageInfo.Page <= 1 {
		return ""
	}
	q := u.Query()
	q.Set("page", strconv.Itoa(pageInfo.Page-1))
	q.Set("per_page", strconv.Itoa(pageInfo.PerPage))
	u.RawQuery = q.Encode()
	return u.RequestURI()
}
