package flaggedemail

import (
	"time"

	"github.com/JakeStrang1/escapehatch/db"
)

type FlaggedEmail struct {
	db.DefaultModel `db:",inline"`
	Email           string    `db:"email"`
	ExpiresAt       time.Time `db:"expires_at"`
}

func (f *FlaggedEmail) CollectionName() string {
	return collection
}
