package http

import (
	"io"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/services/items"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type ItemAPI struct {
	DefaultModelAPI
	MediaType     *string `json:"media_type"`
	ImageURL      *string `json:"image_url"`
	ImageFileName *string `json:"-"`           // Only used by multipart POST request
	ImageFileBody []byte  `json:"-"`           // Only used by multipart POST request
	Description   *string `json:"description"` // Read-only
	CreatedBy     *string `json:"created_by"`
	UserCount     *int    `json:"user_count"`
}

func (i *ItemAPI) BindMultipart(c *gin.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err}
	}

	if len(form.Value["image_url"]) > 0 {
		i.ImageURL = &form.Value["image_url"][0]
	}

	if len(form.File["image_file"]) > 0 {
		i.ImageFileName = &form.File["image_file"][0].Filename
		file, err := form.File["image_file"][0].Open()
		if err != nil {
			return &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err}
		}
		i.ImageFileBody, err = io.ReadAll(file)
		if err != nil {
			return &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err}
		}
	}

	return nil
}

func AddItem(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(CtxKeyUserID)

	item, err := items.Add(userID, id)
	if err != nil {
		Error(c, err)
		return
	}

	var resultAPI any
	switch v := item.(type) {
	case *items.Book:
		resultAPI = ToBookAPI(*v)
	case *items.Movie:
		resultAPI = ToMovieAPI(*v)
	case *items.TVSeries:
		resultAPI = ToTVSeriesAPI(*v)
	default:
	}

	ReturnOne(c, resultAPI)
}

func RemoveItem(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(CtxKeyUserID)

	item, err := items.Remove(userID, id)
	if err != nil {
		Error(c, err)
		return
	}

	var resultAPI any
	switch v := item.(type) {
	case *items.Book:
		resultAPI = ToBookAPI(*v)
	case *items.Movie:
		resultAPI = ToMovieAPI(*v)
	case *items.TVSeries:
		resultAPI = ToTVSeriesAPI(*v)
	default:
	}

	ReturnOne(c, resultAPI)
}

func ToItemAPI(item items.Item) ItemAPI {
	return ItemAPI{
		DefaultModelAPI: ToDefaultModelAPI(item.DefaultModel),
		MediaType:       (*string)(&item.MediaType),
		ImageURL:        IncludeStaticRoot(item.ImageURL),
		Description:     &item.Description,
		CreatedBy:       &item.CreatedBy,
		UserCount:       &item.UserCount,
	}
}

func ToItem(itemAPI ItemAPI) items.Item {
	return items.Item{
		ImageURL:      lo.FromPtr(itemAPI.ImageURL),
		ImageFileName: lo.FromPtr(itemAPI.ImageFileName),
		ImageFileBody: itemAPI.ImageFileBody,
	}
}
