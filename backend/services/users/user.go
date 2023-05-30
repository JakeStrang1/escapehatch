package users

import (
	"regexp"
	"strings"

	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/errors"
)

var usernameCharRegex = regexp.MustCompile(`^[a-z0-9.]+$`) // Only lowercase letters, numbers, and periods. https://regex101.com/r/1s0zQz/1

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

func (u *User) ValidateOnCreate() error {
	if u.Email == "" {
		return &errors.Error{Code: errors.Invalid, Message: "email is missing"}
	}

	if u.Username != "" {
		// Blank username is allowed, we will provide a default

		if len(u.Username) > 20 {
			return &errors.Error{Code: errors.Invalid, Message: "username cannot be greater than 20 characters"}
		}

		if len(u.Username) < 3 {
			return &errors.Error{Code: errors.Invalid, Message: "username cannot be less than 3 characters"}
		}

		if !usernameCharRegex.MatchString(u.Username) {
			return &errors.Error{Code: errors.Invalid, Message: "username must contain only lowercase letters, numbers, and periods"}
		}

		if string(u.Username[0]) == "." || string(u.Username[len(u.Username)-1]) == "." {
			return &errors.Error{Code: errors.Invalid, Message: "username cannot start or end with a period"}
		}

		if strings.Contains(u.Username, "..") {
			return &errors.Error{Code: errors.Invalid, Message: "username cannot have 2 consecutive periods"}
		}
	}

	return nil
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
