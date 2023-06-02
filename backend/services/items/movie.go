package items

import (
	"fmt"
	"strconv"

	"github.com/JakeStrang1/escapehatch/internal/errors"
)

type Movie struct {
	Item           `db:",inline"`
	Title          string   `db:"title"`
	PublishedYear  string   `db:"published_year"`
	IsSeries       bool     `db:"is_series"`
	SeriesTitle    string   `db:"series_title"`
	SequenceNumber string   `db:"sequence_number"`
	LeadActors     []string `db:"lead_actors"`
}

func (b *Movie) GetItem() Item {
	return b.Item
}

func (b *Movie) ValidateMovieOnCreate() error {
	err := b.Item.ValidateOnCreate()
	if err != nil {
		return err
	}

	if b.Title == "" {
		return &errors.Error{Code: errors.Invalid, Message: "title is required"}
	}

	if len(b.LeadActors) < 2 {
		return &errors.Error{Code: errors.Invalid, Message: "provide at least 2 lead actors"}
	}

	for _, actor := range b.LeadActors {
		if actor == "" {
			return &errors.Error{Code: errors.Invalid, Message: "actor name cannot be blank"}
		}
	}

	if b.PublishedYear == "" {
		return &errors.Error{Code: errors.Invalid, Message: "published year is required"}
	}

	if b.IsSeries && b.SeriesTitle == "" {
		return &errors.Error{Code: errors.Invalid, Message: "series title is required"}
	}

	return nil
}

func (b *Movie) SetDescription() {
	// Description: "The Lord of the Rings: The Fellowship of the Ring"
	if !b.IsSeries {
		b.Description = b.Title
		return
	}

	// Description: "The Lord of the Rings: The Fellowship of the Ring (The Lord of the Rings)"
	if b.SequenceNumber == "" {
		b.Description = fmt.Sprintf("%s (%s)", b.Title, b.SeriesTitle)
		return
	}

	// Description: "The Lord of the Rings: The Fellowship of the Ring (The Lord of the Rings #1)"
	if num, err := strconv.Atoi(b.SequenceNumber); err == nil {
		b.Description = fmt.Sprintf("%s (%s #%d)", b.Title, b.SeriesTitle, num)
		return
	}

	// Description: "The Lord of the Rings: The Fellowship of the Ring (The Lord of the Rings Volume 1)"
	b.Description = fmt.Sprintf("%s (%s %s)", b.Title, b.SeriesTitle, b.SequenceNumber)
}

func newMovie(id string) Movie {
	b := Movie{}
	b.ID = id
	return b
}
