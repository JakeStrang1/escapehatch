package session

import (
	"time"

	"github.com/JakeStrang1/escapehatch/db"
)

const collection = "auth_sessions"

func Create(session *Session) error {
	err := db.EnsureUniqueIndex(&Session{}, []string{"token"})
	if err != nil {
		return err
	}
	return db.Create(session)
}

func GetUnexpiredByUserID(userID string, result []Session) error {
	selector := db.M{
		"user_id":    userID,
		"expires_at": db.M{"$gt": time.Now()},
	}
	return db.GetMany(selector, &Session{}, &result)
}

func GetByToken(token string, result *Session) error {
	selector := db.M{
		"token": token,
	}
	return db.GetOne(selector, result)
}

func ExpireByID(id string) error {
	session := Session{}
	session.ID = id
	err := db.GetByID(&session)
	if err != nil {
		return err
	}
	session.ExpiresAt = time.Now()
	return db.Update(&session)
}
