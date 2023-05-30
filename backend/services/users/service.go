package users

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/internal/pages"
)

const userCountCollection = "users_count"
const defaultPerPage = 25
const maxPerPage = 250
const defaultUsernamePrefix = "_" // All default usernames will start with '_' which is not an allowed username char otherwise

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
	err := document.ValidateOnCreate()
	if err != nil {
		return err
	}

	err = db.EnsureUniqueIndex(&User{}, []string{"email"})
	if err != nil {
		return err
	}
	err = db.EnsureUniqueIndex(&User{}, []string{"number"})
	if err != nil {
		return err
	}
	err = db.EnsureUniqueIndex(&User{}, []string{"username"})
	if err != nil {
		return err
	}

	number, err := incrementUserCount()
	if err != nil {
		return err
	}
	document.Number = number

	if document.Username == "" {
		document.Username = GenerateDefaultUsername()
	}

	err = db.Create(document)
	if err != nil {
		return err
	}

	return hydrate(document)
}

func Update(id string, update UserUpdate, result *User) error {
	result.ID = id
	err := db.GetByID(result)
	if err != nil {
		return err
	}

	err = result.ApplyUpdate(update)
	if err != nil {
		return err
	}

	return db.Update(result)
}

func Follow(targetUserID string, followerUserID string, result *User) error {
	if targetUserID == followerUserID {
		return &errors.Error{Code: errors.Invalid, Message: "cannot follow yourself"}
	}

	err := db.EnsureUniqueIndex(&Follower{}, []string{"target_user_id", "follower_user_id"})
	if err != nil {
		return err
	}

	follower := Follower{
		TargetUserID:   targetUserID,
		FollowerUserID: followerUserID,
	}
	err = db.Create(&follower)
	if err != nil {
		return err
	}

	result.ID = targetUserID
	err = GetByID(targetUserID, result)
	if err != nil {
		return err
	}

	return nil
}

func Unfollow(targetUserID string, followerUserID string, result *User) error {
	if targetUserID == followerUserID {
		return &errors.Error{Code: errors.Invalid, Message: "cannot unfollow yourself"}
	}

	selector := db.M{
		"target_user_id":   targetUserID,
		"follower_user_id": followerUserID,
	}
	follower := Follower{}
	err := db.GetOne(selector, &follower)
	if err != nil {
		return err
	}

	err = db.DeleteByID(&follower)
	if err != nil {
		return err
	}

	result.ID = targetUserID
	err = GetByID(targetUserID, result)
	if err != nil {
		return err
	}

	return nil
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

func GetManyFollowers(filter FollowerFilter, results *[]Follower) error {
	// Get user
	user := User{}
	err := GetByID(*filter.TargetUserID, &user)
	if err != nil {
		return err
	}

	// Filter followers
	followers := user.Followers
	if filter.Search != nil {
		followers = user.SearchFollowers(*filter.Search)
	}

	*results = followers
	return nil
}

func GetManyFollowing(filter FollowerFilter, results *[]Follower) error {
	// Get user
	user := User{}
	err := GetByID(*filter.FollowerUserID, &user)
	if err != nil {
		return err
	}

	// Filter following
	following := user.Following
	if filter.Search != nil {
		following = user.SearchFollowing(*filter.Search)
	}

	*results = following
	return nil
}

// GenerateDefaultUsername returns a new username that can be used as a placeholder until a user selects their own
// It has a recognizable prefix so the app can recognize that it needs to be changed.
func GenerateDefaultUsername() string {
	username := defaultUsernamePrefix
	letters := "abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < 8; i++ {
		username = username + string(letters[rand.Intn(len(letters))])
	}
	return username
}

func ValidateUsername(username string) error {
	if len(username) > 20 {
		return &errors.Error{Code: errors.Invalid, Message: "username cannot be greater than 20 characters"}
	}

	if len(username) < 3 {
		return &errors.Error{Code: errors.Invalid, Message: "username cannot be less than 3 characters"}
	}

	if !usernameCharRegex.MatchString(username) {
		return &errors.Error{Code: errors.Invalid, Message: "username must contain only lowercase letters, numbers, and periods"}
	}

	if string(username[0]) == "." || string(username[len(username)-1]) == "." {
		return &errors.Error{Code: errors.Invalid, Message: "username cannot start or end with a period"}
	}

	if strings.Contains(username, "..") {
		return &errors.Error{Code: errors.Invalid, Message: "username cannot have 2 consecutive periods"}
	}

	user := User{}
	err := db.GetOne(db.M{"username": username}, &user)
	if err == nil {
		return &errors.Error{Code: errors.Invalid, Message: "username is already taken"}
	}
	if errors.Code(err) == errors.NotFound {
		// Do nothing, this is good
	} else if err != nil {
		return err
	}

	return nil
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
	// User's followers
	followers := []Follower{}
	err := db.GetMany(db.M{"target_user_id": user.ID}, &Follower{}, &followers)
	if err != nil {
		return err
	}
	for i := range followers {
		err = hydrateFollower(&followers[i])
		if err != nil {
			return err
		}
	}
	user.Followers = followers

	// User's following
	following := []Follower{}
	err = db.GetMany(db.M{"follower_user_id": user.ID}, &Follower{}, &following)
	if err != nil {
		return err
	}
	for i := range following {
		err = hydrateFollower(&following[i])
		if err != nil {
			return err
		}
	}
	user.Following = following
	return nil
}

func hydrateFollower(follower *Follower) error {
	// Target user
	targetUser := User{}
	targetUser.ID = follower.TargetUserID
	err := db.GetByID(&targetUser)
	if err != nil {
		return err
	}
	follower.TargetUsername = targetUser.Username
	follower.TargetFullName = targetUser.FullName

	// Follower user
	followerUser := User{}
	followerUser.ID = follower.FollowerUserID
	err = db.GetByID(&followerUser)
	if err != nil {
		return err
	}
	follower.FollowerUsername = followerUser.Username
	follower.FollowerFullName = followerUser.FullName

	return nil
}
