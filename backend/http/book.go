package http

import (
	"strings"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/services/items"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type BookAPI struct {
	ItemAPI        `json:",inline"`
	Title          *string `json:"title"`
	Author         *string `json:"author"`
	PublishedYear  *string `json:"published_year"`
	IsSeries       *bool   `json:"is_series"`
	SeriesTitle    *string `json:"series_title"`
	SequenceNumber *string `json:"sequence_number"`
}

func (b *BookAPI) BindMultipart(c *gin.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err}
	}

	err = b.ItemAPI.BindMultipart(c)
	if err != nil {
		return err
	}

	if len(form.Value["title"]) > 0 {
		b.Author = &form.Value["title"][0]
	}

	if len(form.Value["author"]) > 0 {
		b.Author = &form.Value["author"][0]
	}

	if len(form.Value["published_year"]) > 0 {
		b.PublishedYear = &form.Value["published_year"][0]
	}

	if len(form.Value["is_series"]) > 0 {
		b.IsSeries = lo.ToPtr(strings.ToLower(form.Value["is_series"][0]) == "true")
	}

	if len(form.Value["series_title"]) > 0 {
		b.SeriesTitle = &form.Value["series_title"][0]
	}

	if len(form.Value["sequence_number"]) > 0 {
		b.SequenceNumber = &form.Value["sequence_number"][0]
	}

	return nil
}

func CreateBook(c *gin.Context) {
	switch c.ContentType() {
	case "multipart/form-data":
		CreateBookMultipart(c)
	default:
		CreateBookJSON(c)
	}
}

func CreateBookMultipart(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)

	body := BookAPI{}
	err := body.BindMultipart(c)
	if err != nil {
		Error(c, &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err})
		return
	}

	result := ToBook(body)
	err = items.CreateBook(userID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToBookAPI(result))
}

func CreateBookJSON(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)

	body := BookAPI{}
	err := Body(c, &body)
	if err != nil {
		Error(c, err)
		return
	}

	result := ToBook(body)
	err = items.CreateBook(userID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToBookAPI(result))
}

func ToBookAPI(book items.Book) BookAPI {
	return BookAPI{
		ItemAPI:        ToItemAPI(book.Item),
		Title:          &book.Title,
		Author:         &book.Author,
		PublishedYear:  &book.PublishedYear,
		IsSeries:       &book.IsSeries,
		SeriesTitle:    &book.SeriesTitle,
		SequenceNumber: &book.SequenceNumber,
	}
}

func ToBook(bookAPI BookAPI) items.Book {
	return items.Book{
		Item:           ToItem(bookAPI.ItemAPI),
		Title:          lo.FromPtr(bookAPI.Title),
		Author:         lo.FromPtr(bookAPI.Author),
		PublishedYear:  lo.FromPtr(bookAPI.PublishedYear),
		IsSeries:       lo.FromPtr(bookAPI.IsSeries),
		SeriesTitle:    lo.FromPtr(bookAPI.SeriesTitle),
		SequenceNumber: lo.FromPtr(bookAPI.SequenceNumber),
	}
}
