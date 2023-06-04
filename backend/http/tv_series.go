package http

import (
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/services/items"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type TVSeriesAPI struct {
	ItemAPI           `json:",inline"`
	Title             *string  `json:"title"`
	TVSeriesStartYear *string  `json:"tv_series_start_year"`
	TVSeriesEndYear   *string  `json:"tv_series_end_year"`
	LeadActors        []string `json:"lead_actors"`
}

func (b *TVSeriesAPI) BindMultipart(c *gin.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err}
	}

	err = b.ItemAPI.BindMultipart(c)
	if err != nil {
		return err
	}

	if len(form.Value["title"]) > 0 {
		b.TVSeriesStartYear = &form.Value["title"][0]
	}

	if len(form.Value["tv_series_start_year"]) > 0 {
		b.TVSeriesStartYear = &form.Value["tv_series_start_year"][0]
	}

	if len(form.Value["tv_series_end_year"]) > 0 {
		b.TVSeriesEndYear = &form.Value["tv_series_end_year"][0]
	}

	if len(form.Value["lead_actors"]) > 0 {
		b.LeadActors = form.Value["lead_actors"]
	}

	return nil
}

func CreateTVSeries(c *gin.Context) {
	switch c.ContentType() {
	case "multipart/form-data":
		CreateTVSeriesMultipart(c)
	default:
		CreateTVSeriesJSON(c)
	}
}

func CreateTVSeriesMultipart(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)

	body := TVSeriesAPI{}
	err := body.BindMultipart(c)
	if err != nil {
		Error(c, &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err})
		return
	}

	result := ToTVSeries(body)
	err = items.CreateTVSeries(userID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToTVSeriesAPI(result))
}

func CreateTVSeriesJSON(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)

	body := TVSeriesAPI{}
	err := Body(c, &body)
	if err != nil {
		Error(c, err)
		return
	}

	result := ToTVSeries(body)
	err = items.CreateTVSeries(userID, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToTVSeriesAPI(result))
}

func UpdateTVSeries(c *gin.Context) {
	switch c.ContentType() {
	case "multipart/form-data":
		UpdateTVSeriesMultipart(c)
	default:
		UpdateTVSeriesJSON(c)
	}
}

func UpdateTVSeriesMultipart(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)
	id := c.Param("id")

	body := TVSeriesAPI{}
	err := body.BindMultipart(c)
	if err != nil {
		Error(c, &errors.Error{Code: errors.BadRequest, Message: "malformed request", Err: err})
		return
	}

	result := items.TVSeries{}
	update := ToTVSeriesUpdate(body)
	err = items.UpdateTVSeries(userID, id, update, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToTVSeriesAPI(result))
}

func UpdateTVSeriesJSON(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)
	id := c.Param("id")

	body := TVSeriesAPI{}
	err := Body(c, &body)
	if err != nil {
		Error(c, err)
		return
	}

	result := items.TVSeries{}
	update := ToTVSeriesUpdate(body)
	err = items.UpdateTVSeries(userID, id, update, &result)
	if err != nil {
		Error(c, err)
		return
	}

	ReturnOne(c, ToTVSeriesAPI(result))
}

func ToTVSeriesAPI(tvSeries items.TVSeries) TVSeriesAPI {
	return TVSeriesAPI{
		ItemAPI:           ToItemAPI(tvSeries.Item),
		Title:             &tvSeries.Title,
		TVSeriesStartYear: &tvSeries.TVSeriesStartYear,
		TVSeriesEndYear:   &tvSeries.TVSeriesEndYear,
		LeadActors:        tvSeries.LeadActors,
	}
}

func ToTVSeries(tvSeriesAPI TVSeriesAPI) items.TVSeries {
	return items.TVSeries{
		Item:              ToItem(tvSeriesAPI.ItemAPI),
		Title:             lo.FromPtr(tvSeriesAPI.Title),
		TVSeriesStartYear: lo.FromPtr(tvSeriesAPI.TVSeriesStartYear),
		TVSeriesEndYear:   lo.FromPtr(tvSeriesAPI.TVSeriesEndYear),
		LeadActors:        tvSeriesAPI.LeadActors,
	}
}

func ToTVSeriesUpdate(tvSeriesAPI TVSeriesAPI) items.TVSeriesUpdate {
	return items.TVSeriesUpdate{
		ItemUpdate:        ToItemUpdate(tvSeriesAPI.ItemAPI),
		Title:             tvSeriesAPI.Title,
		TVSeriesStartYear: tvSeriesAPI.TVSeriesStartYear,
		TVSeriesEndYear:   tvSeriesAPI.TVSeriesEndYear,
		LeadActors:        tvSeriesAPI.LeadActors,
	}
}
