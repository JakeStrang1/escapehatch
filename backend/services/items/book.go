package items

import (
	"fmt"
	"strconv"

	"github.com/JakeStrang1/escapehatch/internal/errors"
)

type BookUpdate struct {
	ItemUpdate
	Title          *string
	Author         *string
	PublishedYear  *string
	IsSeries       *bool
	SeriesTitle    *string
	SequenceNumber *string
}

type Book struct {
	Item           `db:",inline"`
	Title          string `db:"title"`
	Author         string `db:"author"`
	PublishedYear  string `db:"published_year"`
	IsSeries       bool   `db:"is_series"`
	SeriesTitle    string `db:"series_title"`
	SequenceNumber string `db:"sequence_number"`
}

func (b *Book) GetItem() *Item {
	return &b.Item
}

func (b *Book) ValidateOnCreate() error {
	err := b.Item.ValidateOnCreate()
	if err != nil {
		return err
	}

	if b.Title == "" {
		return &errors.Error{Code: errors.Invalid, Message: "title is required"}
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

func (b *Book) Validate() error {
	err := b.Item.Validate()
	if err != nil {
		return err
	}

	if b.Title == "" {
		return &errors.Error{Code: errors.Invalid, Message: "title is required"}
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

func (b *Book) ApplyUpdate(userID string, update BookUpdate) error {
	old := map[string]any{}
	new := map[string]any{}

	if update.Title != nil {
		old["title"] = b.Title
		new["title"] = *update.Title
		b.Title = *update.Title
	}

	if update.Author != nil {
		old["author"] = b.Author
		new["author"] = *update.Author
		b.Author = *update.Author
	}

	if update.PublishedYear != nil {
		old["published_year"] = b.PublishedYear
		new["published_year"] = *update.PublishedYear
		b.PublishedYear = *update.PublishedYear
	}

	if update.IsSeries != nil {
		old["is_series"] = b.IsSeries
		new["is_series"] = *update.IsSeries
		b.IsSeries = *update.IsSeries
	}

	if update.SeriesTitle != nil {
		old["series_title"] = b.SeriesTitle
		new["series_title"] = *update.SeriesTitle
		b.SeriesTitle = *update.SeriesTitle
	}

	if update.SequenceNumber != nil {
		old["sequence_number"] = b.SequenceNumber
		new["sequence_number"] = *update.SequenceNumber
		b.SequenceNumber = *update.SequenceNumber
	}

	err := b.Item.ApplyUpdate(userID, update.ItemUpdate, old, new, nil)
	if err != nil {
		return err
	}

	return nil
}

func (b *Book) SetDescription() {
	// Description: "The Fellowship of the Ring"
	if !b.IsSeries {
		b.Description = b.Title
		return
	}

	// Description: "The Fellowship of the Ring (The Lord of the Rings)"
	if b.SequenceNumber == "" {
		b.Description = fmt.Sprintf("%s (%s)", b.Title, b.SeriesTitle)
		return
	}

	// Description: "The Fellowship of the Ring (The Lord of the Rings #1)"
	if num, err := strconv.Atoi(b.SequenceNumber); err == nil {
		b.Description = fmt.Sprintf("%s (%s #%d)", b.Title, b.SeriesTitle, num)
		return
	}

	// Description: "The Fellowship of the Ring (The Lord of the Rings Volume 1)"
	b.Description = fmt.Sprintf("%s (%s %s)", b.Title, b.SeriesTitle, b.SequenceNumber)
}

func newBook(id string) Book {
	b := Book{}
	b.ID = id
	return b
}
