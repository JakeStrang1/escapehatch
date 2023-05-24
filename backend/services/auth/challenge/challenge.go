package challenge

import (
	"time"

	"github.com/JakeStrang1/saas-template/db"
)

type Challenge struct {
	db.DefaultModel `db:",inline"`
	Email           string            `db:"email"`
	EmailHash       string            `db:"email_hash"`
	SecretHash      string            `db:"secret_hash"`
	Metadata        map[string]string `db:"metadata"`
	RememberMe      bool              `db:"remember_me"`
	IPAddress       string            `db:"ip_address"`
	Attempts        int               `db:"attempts"`
	VerifiedAt      time.Time         `db:"verified_at"`
	ExpiresAt       time.Time         `db:"expires_at"`
	RejectedAt      time.Time         `db:"rejected_at"`
}

func (c *Challenge) CollectionName() string {
	return collection
}
