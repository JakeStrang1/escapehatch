package users

import "github.com/JakeStrang1/escapehatch/db"

// Follower represents a follower/followee relationship between two users
type Follower struct {
	db.DefaultModel  `db:",inline"`
	TargetUserID     string `db:"target_user_id"` // FK: user.id
	TargetUsername   string `db:"-"`
	TargetFullName   string `db:"-"`
	FollowerUserID   string `db:"follower_user_id"` // FK: user.id
	FollowerUsername string `db:"-"`
	FollowerFullName string `db:"-"`
}
