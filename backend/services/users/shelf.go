package users

import (
	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/errors"
)

type Shelf struct {
	db.DefaultModel `db:",inline"`
	Name            string      `db:"name"`
	Items           []ShelfItem `db:"items"`
}

func (s *Shelf) HasItem(itemID string) bool {
	for _, item := range s.Items {
		if item.ItemID == itemID {
			return true
		}
	}
	return false
}

func (s *Shelf) AddItem(item ShelfItem) error {
	err := s.Saving() // Call hook manually to set UpdatedAt
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	s.Items = append(s.Items, item)
	return nil
}

type ShelfItem struct {
	ItemID      string `db:"item_id"`     // FK: item.id
	Description string `db:"description"` // Cache: item.Description (calculated from item fields)
	ImageURL    string `db:"image_url"`   // Cache: item.image_url
}
