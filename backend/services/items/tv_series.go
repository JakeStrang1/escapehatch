package items

import (
	"fmt"

	"github.com/JakeStrang1/escapehatch/internal/errors"
)

type TVSeries struct {
	Item              `db:",inline"`
	Title             string   `db:"title"`
	LeadActors        []string `db:"lead_actors"`
	TVSeriesStartYear string   `db:"tv_series_start_year"`
	TVSeriesEndYear   string   `db:"tv_series_end_year"`
}

func (b *TVSeries) GetItem() Item {
	return b.Item
}

func (b *TVSeries) ValidateTVSeriesOnCreate() error {
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

	if b.TVSeriesStartYear == "" {
		return &errors.Error{Code: errors.Invalid, Message: "series start year is required"}
	}

	if b.TVSeriesEndYear == "" {
		return &errors.Error{Code: errors.Invalid, Message: "series end year is required (can be \"present\" if still ongoing)"}
	}

	return nil
}

func (b *TVSeries) SetDescription() {
	b.Description = fmt.Sprintf("%s (%s - %s)", b.Title, b.TVSeriesStartYear, b.TVSeriesEndYear)
}

func newTVSeries(id string) TVSeries {
	b := TVSeries{}
	b.ID = id
	return b
}
