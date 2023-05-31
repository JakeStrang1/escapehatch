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
	ImageFileName *string `json:"-"`
	ImageFileBody []byte  `json:"-"`
	Title         *string `json:"title"`
	CreatedBy     *string `json:"created_by"`
}

func (i *ItemAPI) BindMultipart(c *gin.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err}
	}

	if len(form.Value["image_url"]) > 0 {
		i.ImageURL = &form.Value["image_url"][0]
	}

	if len(form.Value["title"]) > 0 {
		i.Title = &form.Value["title"][0]
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

func ToItemAPI(item items.Item) ItemAPI {
	return ItemAPI{
		DefaultModelAPI: ToDefaultModelAPI(item.DefaultModel),
		MediaType:       (*string)(&item.MediaType),
		ImageURL:        &item.ImageURL,
		Title:           &item.Title,
		CreatedBy:       &item.CreatedBy,
	}
}

func ToItem(itemAPI ItemAPI) items.Item {
	return items.Item{
		ImageURL:      lo.FromPtr(itemAPI.ImageURL),
		ImageFileName: lo.FromPtr(itemAPI.ImageFileName),
		ImageFileBody: itemAPI.ImageFileBody,
		Title:         lo.FromPtr(itemAPI.Title),
	}
}
