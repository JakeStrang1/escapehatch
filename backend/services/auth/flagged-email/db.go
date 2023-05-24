package flaggedemail

import (
	"time"

	"github.com/JakeStrang1/escapehatch/db"
)

const collection = "auth_flagged_emails"

func Create(flaggedEmail *FlaggedEmail) error {
	return db.Create(flaggedEmail)
}

func GetUnexpiredByEmail(email string, flaggedEmail *FlaggedEmail) error {
	return db.GetOne(db.M{"email": email, "expires_at": db.M{"$gt": time.Now()}}, flaggedEmail)
}
