package users

import (
	"github.com/JakeStrang1/escapehatch/db"
)

type User struct {
	db.DefaultModel `db:",inline"`
	Email           string     `db:"email"`    // Unique
	Username        string     `db:"username"` // Unique, not yet implemented
	Number          int        `db:"number"`   // Unique, number indicating how early you were to the platform (#1, #2, #3, etc..)
	FullName        string     `db:"full_name"`
	Shelves         []Shelf    `db:"shelves"`
	Followers       []Follower `db:"-"` // Hydrated
	Following       []Follower `db:"-"` // Hydrated
}

// Follows returns true if the user follows the given user
func (u *User) Follows(userID string) bool {
	for _, following := range u.Following {
		if following.TargetUserID == userID {
			return true
		}
	}
	return false
}

// FollowedBy returns true if the user is followed by the given user
func (u *User) FollowedBy(userID string) bool {
	for _, follower := range u.Followers {
		if follower.FollowerID == userID {
			return true
		}
	}
	return false
}
