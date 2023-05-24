package session

import (
	"time"

	"github.com/JakeStrang1/saas-template/db"
)

type Session struct {
	db.DefaultModel `db:",inline"`
	UserID          string    `db:"user_id"`
	Token           string    `db:"token"`
	RememberMe      bool      `db:"remember_me"`
	ExpiresAt       time.Time `db:"expires_at"`
}

func (s *Session) CollectionName() string {
	return collection
}
