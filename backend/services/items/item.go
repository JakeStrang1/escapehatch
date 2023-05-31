package items

import (
	"time"

	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/errors"
)

type MediaType string

const (
	MediaTypeBook     MediaType = "book"
	MediaTypeMovie    MediaType = "movie"
	MediaTypeTVSeries MediaType = "tv_series"
)

type Item struct {
	db.DefaultModel `db:",inline"`
	MediaType       MediaType `db:"media_type"`
	ImageURL        string    `db:"image_url"`
	ImageFileName   string    `db:"-"`
	ImageFileBody   []byte    `db:"-"`
	Title           string    `db:"title"`
	CreatedBy       string    `db:"created_by"`
	ChangeLog       []Change  `db:"change_log"`
}

func (i *Item) ValidateOnCreate() error {
	if i.ImageURL == "" && len(i.ImageFileBody) == 0 {
		return &errors.Error{Code: errors.Invalid, Message: "image is required"}
	}

	if i.Title == "" {
		return &errors.Error{Code: errors.Invalid, Message: "title is required"}
	}

	return nil
}

type Book struct {
	Item           `db:",inline"`
	Author         string `db:"author"`
	PublishedYear  string `db:"published_year"`
	IsSeries       bool   `db:"is_series"`
	SeriesTitle    string `db:"series_title"`
	SequenceNumber string `db:"sequence_number"`
}

func (b *Book) ValidateBookOnCreate() error {
	err := b.Item.ValidateOnCreate()
	if err != nil {
		return err
	}

	if b.Author == "" {
		return &errors.Error{Code: errors.Invalid, Message: "author is required"}
	}

	if b.PublishedYear == "" {
		return &errors.Error{Code: errors.Invalid, Message: "published year is required"}
	}

	if b.IsSeries && b.SeriesTitle == "" {
		return &errors.Error{Code: errors.Invalid, Message: "series title is required"}
	}

	return nil
}

type Movie struct {
	Item           `db:",inline"`
	PublishedYear  string   `db:"published_year"`
	IsSeries       bool     `db:"is_series"`
	SeriesTitle    string   `db:"series_title"`
	SequenceNumber string   `db:"sequence_number"`
	LeadActors     []string `db:"lead_actors"`
}

type TVSeries struct {
	Item              `db:",inline"`
	LeadActors        []string `db:"lead_actors"`
	TVSeriesStartYear string   `db:"tv_series_start_year"`
	TVSeriesEndYear   string   `db:"tv_series_end_year"`
}

type Change struct {
	UpdatedAt time.Time              `db:"updated_at"`
	UpdatedBy string                 `db:"updated_by"`
	Old       map[string]interface{} `db:"old"`
	New       map[string]interface{} `db:"new"`
}
