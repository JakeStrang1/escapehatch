package auth

import (
	"time"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/services/auth/challenge"
	flaggedemail "github.com/JakeStrang1/escapehatch/services/auth/flagged-email"
	flaggedip "github.com/JakeStrang1/escapehatch/services/auth/flagged-ip"
	nocontact "github.com/JakeStrang1/escapehatch/services/auth/no-contact"
	"github.com/JakeStrang1/escapehatch/services/auth/session"
	"github.com/JakeStrang1/escapehatch/services/users"
	"golang.org/x/crypto/bcrypt"
)

func NotYou(secret, emailHash string, doNotContact bool) error {
	ch, err := validateChallenge(secret, emailHash)
	if err != nil {
		return err
	}

	// The following are all independent measures, error checking at the end

	var doNotContactErr error
	if doNotContact {
		doNotContactErr = createDoNoContact(ch.Email)
	}

	rejectErr := rejectChallenges(ch.Email)

	expireErr := expireSessions(ch.Email)

	flagIPErr := flaggedip.Create(&flaggedip.FlaggedIP{
		Address:   ch.IPAddress,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})

	flagEmailErr := flaggedemail.Create(&flaggedemail.FlaggedEmail{
		Email:     ch.Email,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	})

	return firstError(doNotContactErr, rejectErr, expireErr, flagIPErr, flagEmailErr)
}

func validateChallenge(secret, emailHash string) (*challenge.Challenge, error) {
	ch := challenge.Challenge{}
	err := challenge.GetByEmailHash(emailHash, &ch)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(ch.SecretHash), []byte(secret))
	if err != nil {
		return nil, &errors.Error{Code: errors.ChallengeFailed, Message: "secret is incorrect", Err: err}
	}

	return &ch, nil
}

func createDoNoContact(email string) error {
	err := nocontact.Create(&nocontact.NoContact{
		Email: email,
	})
	return err
}

func rejectChallenges(email string) error {

	challenges := []challenge.Challenge{}
	err := challenge.GetUnverifiedByEmail(email, challenges)
	if err != nil {
		return err
	}

	err = nil
	for _, ch := range challenges {
		rejectErr := challenge.RejectByID(ch.ID)
		if err != nil {
			err = rejectErr
		}
	}
	if errors.Code(err) == errors.NotFound {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return err
}

func expireSessions(email string) error {
	exists, err := users.Exists(email)
	if err != nil {
		return err
	}
	if !exists {
		return nil // No action required
	}

	u := users.User{}
	err = users.GetByEmail(email, &u)
	if errors.Code(err) == errors.NotFound {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	if err != nil {
		return err
	}

	sessions := []session.Session{}
	err = session.GetUnexpiredByUserID(u.ID, sessions)
	if err != nil {
		return err
	}

	err = nil
	for _, s := range sessions {
		expireErr := session.ExpireByID(s.ID)
		if err != nil {
			err = expireErr
		}
	}
	if errors.Code(err) == errors.NotFound {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return err
}

func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
