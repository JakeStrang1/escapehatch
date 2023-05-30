package users

import "github.com/JakeStrang1/escapehatch/db"

type Shelf struct {
	db.DefaultModel `db:",inline"`
	Name            string      `db:"name"`
	Items           []ShelfItem `db:"items"`
}

type ShelfItem struct {
	ItemID string `db:"item_id"` // FK: item.id
	Image  string `db:"image"`   // Cache: item.image
}
