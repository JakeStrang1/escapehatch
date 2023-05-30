package users

import "github.com/JakeStrang1/escapehatch/db"

// Follower represents a followee/follower relationship between two users
type Follower struct {
	db.DefaultModel `db:",inline"`
	TargetUserID    string `db:"target_user_id"`   // FK: user.id
	FollowerID      string `db:"follower_user_id"` // FK: user.id
}
