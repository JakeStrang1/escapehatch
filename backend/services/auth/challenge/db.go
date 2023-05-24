package challenge

import (
	"time"

	"github.com/JakeStrang1/escapehatch/db"
)

const collection = "auth_challenges"

func Create(challenge *Challenge) error {
	return db.Create(challenge)
}

func GetUnverifiedByEmail(email string, result []Challenge) error {
	selector := db.M{
		"email":       email,
		"verified_at": time.Time{},
	}
	return db.GetMany(selector, &Challenge{}, &result)
}

func GetUnverifiedByEmailAndDate(email string, createdSince time.Time, result []Challenge) error {
	selector := db.M{
		"email":       email,
		"created_at":  db.M{"$gt": createdSince},
		"verified_at": time.Time{},
	}
	return db.GetMany(selector, &Challenge{}, &result)
}

func GetByEmailHash(emailHash string, result *Challenge) error {
	return db.GetOne(db.M{"email_hash": emailHash}, result)
}

func GetLatestByEmail(email string, result *Challenge) error {
	opts := db.GetOneOptions{
		Sort: [][]interface{}{{"created_at", -1}},
	}
	return db.GetOne(db.M{"email": email}, result, opts)
}

func VerifyByID(id string) error {
	challenge := Challenge{}
	challenge.ID = id
	err := db.GetByID(&challenge)
	if err != nil {
		return err
	}
	challenge.VerifiedAt = time.Now()
	return db.Update(&challenge)
}

func RejectByID(id string) error {
	challenge := Challenge{}
	challenge.ID = id
	err := db.GetByID(&challenge)
	if err != nil {
		return err
	}
	challenge.RejectedAt = time.Now()
	return db.Update(&challenge)
}

func IncrementAttempts(id string, result *Challenge) error {
	return db.IncrementOne(db.M{"_id": id}, "attempts", result)
}
