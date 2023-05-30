package users

import (
	"fmt"

	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/internal/pages"
)

const userCountCollection = "users_count"
const defaultPerPage = 25
const maxPerPage = 250

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

type UserCount struct {
	db.DefaultModel `db:",inline"`
	CurrentNumber   int `db:"current_number"`
}

func (u *UserCount) CollectionName() string {
	return userCountCollection
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
	err = db.EnsureUniqueIndex(&User{}, []string{"number"})
	if err != nil {
		return err
	}

	number, err := incrementUserCount()
	if err != nil {
		return err
	}
	document.Number = number

	err = db.Create(document)
	if err != nil {
		return err
	}

	return hydrate(document)
}

func GetPage(filter Filter, results *[]User) (*pages.PageResult, error) {
	err := filter.Validate()
	if err != nil {
		return nil, err
	}

	page := 1
	if filter.Page != nil {
		page = *filter.Page
	}

	pageSize := defaultPerPage
	if filter.PerPage != nil {
		pageSize = *filter.PerPage
	}

	total, err := db.GetPage(filter, &User{}, results, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &pages.PageResult{
		Page:       page,
		PerPage:    pageSize,
		TotalPages: total,
	}, nil
}

func GetByID(id string, result *User) error {
	result.ID = id
	err := db.GetByID(result)
	if err != nil {
		return err
	}

	return hydrate(result)
}

func GetByEmail(email string, result *User) error {
	err := db.GetOne(db.M{"email": email}, result)
	if err != nil {
		return err
	}

	return hydrate(result)
}

// incrementUserCount increments the users_count document and returns an atomically-reserved user number to be used for a new user
func incrementUserCount() (int, error) {
	userCount, err := db.GetCount(db.M{}, &UserCount{})
	if err != nil {
		return 0, err
	}

	if userCount == 0 {
		db.Create(&UserCount{CurrentNumber: 0})
	}

	if userCount > 1 {
		return 0, &errors.Error{Code: errors.Internal, Err: fmt.Errorf("cannot have multiple %s documents", userCountCollection)}
	}

	result := UserCount{}
	err = db.IncrementOne(db.M{}, "current_number", &result)
	if err != nil {
		return 0, err
	}
	return result.CurrentNumber, nil
}

func hydrate(user *User) error {
	followers := []Follower{}
	err := db.GetMany(db.M{"target_user_id": user.ID}, user, &followers)
	if err != nil {
		return err
	}
	user.Followers = followers

	following := []Follower{}
	err = db.GetMany(db.M{"follower_user_id": user.ID}, user, &following)
	if err != nil {
		return err
	}
	user.Following = following
	return nil
}
