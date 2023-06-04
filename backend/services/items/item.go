package items

import (
	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/errors"
)

type MediaType string

const (
	MediaTypeBook     MediaType = "book"
	MediaTypeMovie    MediaType = "movie"
	MediaTypeTVSeries MediaType = "tv_series"
)

type ChangeType string

const (
	ChangeTypeCreate ChangeType = "create"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
)

// ItemContainer describes a type that contains an item. All media type structs should implement this interface.
type ItemContainer interface {
	GetItem() *Item
}

type ItemUpdate struct {
	ImageURL      *string `db:"image_url,omitempty"`
	ImageFileName *string `db:"-"`
	ImageFileBody []byte  `db:"-"`
}

type Item struct {
	db.DefaultModel `db:",inline"`
	MediaType       MediaType `db:"media_type"`
	ImageURL        string    `db:"image_url"`
	ImageFileName   string    `db:"-"`
	ImageFileBody   []byte    `db:"-"`
	Description     string    `db:"title"`
	CreatedBy       string    `db:"created_by"`
	UserCount       int       `db:"-"` // How many users have added this to at least one shelf
}

func (i *Item) ValidateOnCreate() error {
	if i.ImageURL == "" && len(i.ImageFileBody) == 0 {
		return &errors.Error{Code: errors.Invalid, Message: "image is required"}
	}
	return nil
}

func (i *Item) Validate() error {
	if i.ImageURL == "" && len(i.ImageFileBody) == 0 {
		return &errors.Error{Code: errors.Invalid, Message: "image is required"}
	}
	return nil
}

func (i *Item) ApplyUpdate(userID string, update ItemUpdate) error {
	if update.ImageURL != nil {
		i.ImageURL = *update.ImageURL
	}

	if update.ImageFileName != nil {
		i.ImageFileName = *update.ImageFileName
	}

	if update.ImageFileBody != nil {
		i.ImageFileBody = update.ImageFileBody
	}

	return nil
}
