package users

import (
	"github.com/JakeStrang1/saas-template/db"
	"github.com/JakeStrang1/saas-template/internal/errors"
)

type User struct {
	db.DefaultModel `db:",inline"`
	Email           string `db:"email"`
}

func Exists(email string) (bool, error) {
	u := User{}
	err := GetByEmail(email, &u)
	if errors.Code(err) == errors.NotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func Create(document *User) error {
	err := db.EnsureUniqueIndex(&User{}, []string{"email"})
	if err != nil {
		return err
	}
	return db.Create(document)
}

func GetByID(id string, result *User) error {
	result.ID = id
	return db.GetByID(result)
}

func GetByEmail(email string, result *User) error {
	return db.GetOne(db.M{"email": email}, result)
}
