package http

import (
	"strings"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/services/items"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type MovieAPI struct {
	ItemAPI        `json:",inline"`
	Title          *string  `json:"title"`
	PublishedYear  *string  `json:"published_year"`
	IsSeries       *bool    `json:"is_series"`
	SeriesTitle    *string  `json:"series_title"`
	SequenceNumber *string  `json:"sequence_number"`
	LeadActors     []string `json:"lead_actors"`
}

func (b *MovieAPI) BindMultipart(c *gin.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err}
	}

	err = b.ItemAPI.BindMultipart(c)
	if err != nil {
		return err
	}

	if len(form.Value["title"]) > 0 {
		b.SeriesTitle = &form.Value["title"][0]
	}

	if len(form.Value["published_year"]) > 0 {
		b.PublishedYear = &form.Value["published_year"][0]
	}

	if len(form.Value["series_title"]) > 0 {
		b.SeriesTitle = &form.Value["series_title"][0]
	}

	if len(form.Value["is_series"]) > 0 {
		b.IsSeries = lo.ToPtr(strings.ToLower(form.Value["is_series"][0]) == "true")
	}

	if len(form.Value["sequence_number"]) > 0 {
		b.SequenceNumber = &form.Value["sequence_number"][0]
	}

	if len(form.Value["lead_actors"]) > 0 {
		b.LeadActors = form.Value["lead_actors"]
	}

	return nil
}

func CreateMovie(c *gin.Context) {
	switch c.ContentType() {
	case "multipart/form-data":
		CreateMovieMultipart(c)
	default:
		CreateMovieJSON(c)
	}
}

func CreateMovieMultipart(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)

	body := MovieAPI{}
	err := body.BindMultipart(c)
	if err != nil {
		Error(c, &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err})
		return
	}

	result := ToMovie(body)
	err = items.CreateMovie(userID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToMovieAPI(result))
}

func CreateMovieJSON(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)

	body := MovieAPI{}
	err := Body(c, &body)
	if err != nil {
		Error(c, err)
		return
	}

	result := ToMovie(body)
	err = items.CreateMovie(userID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToMovieAPI(result))
}

func ToMovieAPI(movie items.Movie) MovieAPI {
	return MovieAPI{
		ItemAPI:        ToItemAPI(movie.Item),
		Title:          &movie.Title,
		PublishedYear:  &movie.PublishedYear,
		IsSeries:       &movie.IsSeries,
		SeriesTitle:    &movie.SeriesTitle,
		SequenceNumber: &movie.SequenceNumber,
		LeadActors:     movie.LeadActors,
	}
}

func ToMovie(movieAPI MovieAPI) items.Movie {
	return items.Movie{
		Item:           ToItem(movieAPI.ItemAPI),
		Title:          lo.FromPtr(movieAPI.Title),
		PublishedYear:  lo.FromPtr(movieAPI.PublishedYear),
		IsSeries:       lo.FromPtr(movieAPI.IsSeries),
		SeriesTitle:    lo.FromPtr(movieAPI.SeriesTitle),
		SequenceNumber: lo.FromPtr(movieAPI.SequenceNumber),
		LeadActors:     movieAPI.LeadActors,
	}
}
