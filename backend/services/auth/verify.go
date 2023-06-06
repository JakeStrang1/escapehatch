package auth

import (
	"strings"
	"time"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/services/auth/challenge"
	"github.com/JakeStrang1/escapehatch/services/auth/session"
	"github.com/JakeStrang1/escapehatch/services/users"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const maxAttempts = 5 // Maximum number of attempts for verifying challenge

// - Endpoint to receive codes and return a session cookie to the client
// 	- Checks cache for a valid entry for email/passphrase that isn't expired (also storing bcrypted email in cache for searching)
// 	- If not found, return error.
// 	- If found, but "Not you" is set, flag the IP of the original sender somewhere. If "block permanently" then save the email to the do-not-email list in DB (bcrypted). Return.
// 	- If found, but expired, return expired error.
// 	- Check if user exists. Create if not.
// 	- Saves a session token for the user in db/cache
// 	- Returns the cookie to the client (with expiration date if remember me was selected)
// 	- Returns the original url being targeted (metadata)

func Verify(email, secret, emailHash string) (*session.Session, map[string]string, error) {
	ch, err := verifyChallenge(email, secret, emailHash)
	if err != nil {
		return nil, nil, err
	}
	email = ch.Email

	u, err := ensureUserExists(email)
	if err != nil {
		return nil, nil, err
	}

	session, err := createSession(u.ID, ch.RememberMe)
	if err != nil {
		return nil, nil, err
	}

	return session, ch.Metadata, nil
}

func verifyChallenge(email, secret, emailHash string) (*challenge.Challenge, error) {
	if email == "" {
		ch := challenge.Challenge{}
		err := challenge.GetByEmailHash(emailHash, &ch)
		if err != nil {
			return nil, err
		}
		email = ch.Email
	}

	ch := challenge.Challenge{}
	err := challenge.GetLatestByEmail(email, &ch)
	if err != nil {
		return nil, err
	}

	err = challenge.IncrementAttempts(ch.ID, &ch)
	if errors.Code(err) == errors.NotFound {
		// Would be strange for this to occur
		return nil, &errors.Error{Code: errors.Unexpected, Err: err}
	}
	if err != nil {
		return nil, err
	}
	if ch.Attempts > maxAttempts {
		return nil, errors.New(errors.ChallengeMaxAttempts, "too many attempts")
	}

	err = bcrypt.CompareHashAndPassword([]byte(ch.SecretHash), []byte(strings.ToLower(secret)))
	if err != nil {
		return nil, &errors.Error{Code: errors.ChallengeFailed, Message: "secret is incorrect", Err: err}
	}

	if time.Now().After(ch.ExpiresAt) {
		return nil, errors.New(errors.ChallengeExpired, "challenge has already expired")
	}

	zeroTime := time.Time{}
	if ch.RejectedAt != zeroTime {
		return nil, errors.New(errors.ChallengeInvalidated, "challenge has been flagged as suspicious")
	}

	if ch.VerifiedAt != zeroTime {
		return nil, errors.New(errors.ChallengeAlreadyUsed, "challenge has already been used")
	}

	err = challenge.VerifyByID(ch.ID)
	if errors.Code(err) == errors.NotFound {
		// Would be strange for this to occur
		return nil, &errors.Error{Code: errors.Unexpected, Err: err}
	}
	if err != nil {
		return nil, err
	}

	return &ch, nil
}

func ensureUserExists(email string) (*users.User, error) {
	exists, err := users.Exists(email)
	if err != nil {
		return nil, err
	}

	u := users.User{}
	if !exists {
		u.Email = email
		err = users.Create(&u)
	} else {
		err = users.GetByEmail(email, &u)
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func createSession(userID string, rememberMe bool) (*session.Session, error) {
	token := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)
	if rememberMe {
		expiresAt = time.Now().AddDate(0, 0, 30)
	}
	s := session.Session{
		UserID:     userID,
		Token:      token.String(),
		RememberMe: rememberMe,
		ExpiresAt:  expiresAt,
	}
	err := session.Create(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
