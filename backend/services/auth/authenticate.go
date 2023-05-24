package auth

import (
	"time"

	"github.com/JakeStrang1/saas-template/internal/errors"
	"github.com/JakeStrang1/saas-template/services/auth/session"
	"github.com/JakeStrang1/saas-template/services/users"
)

func Authenticate(sessionToken string) (*users.User, error) {
	if sessionToken == "" {
		return nil, errors.New(errors.SessionInvalid, "session token cannot be blank, please sign in")
	}

	s := session.Session{}
	err := session.GetByToken(sessionToken, &s)
	if errors.Code(err) == errors.NotFound {
		return nil, &errors.Error{Code: errors.SessionInvalid, Message: "unknown session token, please sign in", Err: err}
	}
	if err != nil {
		return nil, err
	}

	if s.ExpiresAt.Before(time.Now()) {
		return nil, errors.New(errors.SessionExpired, "session has expired, please sign in again")
	}

	u := users.User{}
	err = users.GetByID(s.UserID, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
