package users

import "github.com/JakeStrang1/escapehatch/db"

type UserCount struct {
	db.DefaultModel `db:",inline"`
	CurrentNumber   int `db:"current_number"`
}

func (u *UserCount) CollectionName() string {
	return userCountCollection
}
