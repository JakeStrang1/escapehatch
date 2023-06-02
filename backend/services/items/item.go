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

// ItemContainer describes a type that contains an item. All media type structs should implement this interface.
type ItemContainer interface {
	GetItem() Item
}

type Item struct {
	db.DefaultModel `db:",inline"`
	MediaType       MediaType `db:"media_type"`
	ImageURL        string    `db:"image_url"`
	ImageFileName   string    `db:"-"`
	ImageFileBody   []byte    `db:"-"`
	Description     string    `db:"title"`
	CreatedBy       string    `db:"created_by"`
	ChangeLog       []Change  `db:"change_log"`
	UserCount       int       `db:"-"` // How many users have added this to at least one shelf
}

func (i *Item) ValidateOnCreate() error {
	if i.ImageURL == "" && len(i.ImageFileBody) == 0 {
		return &errors.Error{Code: errors.Invalid, Message: "image is required"}
	}
	return nil
}

type Change struct {
	UpdatedAt time.Time              `db:"updated_at"`
	UpdatedBy string                 `db:"updated_by"`
	Old       map[string]interface{} `db:"old"`
	New       map[string]interface{} `db:"new"`
}
