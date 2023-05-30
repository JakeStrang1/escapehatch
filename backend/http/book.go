package http

import (
	"github.com/JakeStrang1/escapehatch/services/items"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type BookAPI struct {
	ItemAPI        `json:",inline"`
	Author         *string `json:"author"`
	PublishedYear  *string `json:"published_year"`
	IsSeries       *bool   `json:"is_series"`
	SeriesTitle    *string `json:"series_title"`
	SequenceNumber *string `json:"sequence_number"`
}

func CreateBook(c *gin.Context) {
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
		Author:         lo.FromPtr(bookAPI.Author),
		PublishedYear:  lo.FromPtr(bookAPI.PublishedYear),
		IsSeries:       lo.FromPtr(bookAPI.IsSeries),
		SeriesTitle:    lo.FromPtr(bookAPI.SeriesTitle),
		SequenceNumber: lo.FromPtr(bookAPI.SequenceNumber),
	}
}
