package flaggedip

import (
	"time"

	"github.com/JakeStrang1/escapehatch/db"
)

type FlaggedIP struct {
	db.DefaultModel `db:",inline"`
	Address         string    `db:"address"`
	ExpiresAt       time.Time `db:"expires_at"`
}

func (f *FlaggedIP) CollectionName() string {
	return collection
}
