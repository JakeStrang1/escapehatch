package flaggedip

import (
	"time"

	"github.com/JakeStrang1/saas-template/db"
)

const collection = "auth_flagged_ips"

func Create(flaggedIP *FlaggedIP) error {
	return db.Create(flaggedIP)
}

func GetUnexpiredByIP(ipAddress string, result *FlaggedIP) error {
	return db.GetOne(db.M{"address": ipAddress, "expires_at": db.M{"$gt": time.Now()}}, result)
}
