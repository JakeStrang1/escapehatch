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

type ChangeType string

const (
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
)

// ItemContainer describes a type that contains an item. All media type structs should implement this interface.
type ItemContainer interface {
	GetItem() *Item
}

type ItemUpdate struct {
	ImageURL      *string
	ImageFileName *string
	ImageFileBody []byte
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

func (i *Item) Validate() error {
	if i.ImageURL == "" && len(i.ImageFileBody) == 0 {
		return &errors.Error{Code: errors.Invalid, Message: "image is required"}
	}
	return nil
}

func (i *Item) ApplyUpdate(userID string, update ItemUpdate, old map[string]any, new map[string]any, metadata map[string]any) error {
	if update.ImageURL != nil {
		old["image_url"] = i.ImageURL
		new["image_url"] = *update.ImageURL
		i.ImageURL = *update.ImageURL
	}

	if update.ImageFileName != nil {
		old["image_url"] = i.ImageURL
		new["image_url"] = "from file"
		i.ImageFileName = *update.ImageFileName
	}

	if update.ImageFileBody != nil {
		old["image_url"] = i.ImageURL
		new["image_url"] = "from file"
		i.ImageFileBody = update.ImageFileBody
	}

	i.ChangeLog = append(i.ChangeLog, Change{
		ChangeType: ChangeTypeUpdate,
		UpdatedAt:  time.Now(),
		UpdatedBy:  userID,
		Old:        old,
		New:        new,
		Metadata:   metadata,
	})

	return nil
}

func (i *Item) MarkDeleted(reason string, userID string, userCount int) {
	i.ChangeLog = append(i.ChangeLog, Change{
		ChangeType: ChangeTypeDelete,
		UpdatedAt:  time.Now(),
		UpdatedBy:  userID,
		Metadata: map[string]any{
			"reason":         reason,
			"users_impacted": userCount,
		},
	})
}

type Change struct {
	ChangeType ChangeType     `db:"change_type"`
	UpdatedAt  time.Time      `db:"updated_at"`
	UpdatedBy  string         `db:"updated_by"`
	Old        map[string]any `db:"old"`
	New        map[string]any `db:"new"`
	Metadata   map[string]any `db:"metadata"`
}
