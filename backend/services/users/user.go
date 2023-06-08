package users

import (
	"regexp"
	"strings"

	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/errors"
)

var UserSearchPaths = []string{"username", "full_name"}
var usernameCharRegex = regexp.MustCompile(`^[a-z0-9.]+$`) // Only lowercase letters, numbers, and periods. https://regex101.com/r/1s0zQz/1

type UserUpdate struct {
	Username *string `db:"username"`
	FullName *string `db:"full_name"`
}

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

	if u.Username != "" { // Blank username is allowed on create, we will provide a default
		u.Username = strings.ToLower(u.Username) // Force lowercase
		err := ValidateUsername(u.Username)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *User) ApplyUpdate(update UserUpdate) error {
	if update.Username != nil {
		err := ValidateUsername(*update.Username)
		if err != nil {
			return err
		}
		u.Username = *update.Username
	}
	if update.FullName != nil {
		if *update.FullName == "" {
			return &errors.Error{Code: errors.Invalid, Message: "full name cannot be blank"}
		}
		u.FullName = *update.FullName
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
		if follower.FollowerUserID == userID {
			return true
		}
	}
	return false
}

// SearchFollowers returns all the user's followers that contain the given search string in their username or full name
func (u *User) SearchFollowers(search string) []Follower {
	search = strings.ToLower(search) // case-insensitive search
	results := []Follower{}
	for _, follower := range u.Followers {
		if strings.Contains(follower.FollowerUsername, search) {
			results = append(results, follower)
		} else if strings.Contains(strings.ToLower(follower.FollowerFullName), search) {
			results = append(results, follower)
		}
	}
	return results
}

// SearchFollowing returns everyone that the user is following that contains the given search string in their username or full name
func (u *User) SearchFollowing(search string) []Follower {
	search = strings.ToLower(search) // case-insensitive search
	results := []Follower{}
	for _, following := range u.Following {
		if strings.Contains(following.TargetUsername, search) {
			results = append(results, following)
		} else if strings.Contains(strings.ToLower(following.TargetFullName), search) {
			results = append(results, following)
		}
	}
	return results
}

func (u *User) GetShelfByName(result *Shelf) bool {
	for _, shelf := range u.Shelves {
		if shelf.Name == result.Name {
			*result = shelf
			return true
		}
	}
	return false
}

func (u *User) AddShelf(result *Shelf) error {
	err := result.Creating() // Call hook manually to set ID and CreatedAt
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	u.Shelves = append(u.Shelves, *result)
	return nil
}

func (u *User) HasItem(itemID string) bool {
	for _, shelf := range u.Shelves {
		for _, item := range shelf.Items {
			if item.ItemID == itemID {
				return true
			}
		}
	}
	return false
}

func (u *User) UpdateShelf(shelf *Shelf) {
	for i := range u.Shelves {
		if u.Shelves[i].ID == shelf.ID {
			u.Shelves[i] = *shelf
		}
	}
}

func (u *User) RemoveItemFromAllShelves(itemID string) error {
	for i := range u.Shelves {
		if u.Shelves[i].HasItem(itemID) {
			err := u.Shelves[i].RemoveItem(itemID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
