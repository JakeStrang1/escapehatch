package auth

import (
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/JakeStrang1/escapehatch/email"
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/internal/secret"
	"github.com/JakeStrang1/escapehatch/services/auth/challenge"
	flaggedemail "github.com/JakeStrang1/escapehatch/services/auth/flagged-email"
	flaggedip "github.com/JakeStrang1/escapehatch/services/auth/flagged-ip"
	nocontact "github.com/JakeStrang1/escapehatch/services/auth/no-contact"
	"github.com/JakeStrang1/escapehatch/services/auth/templates"
	"github.com/JakeStrang1/escapehatch/services/users"
	"golang.org/x/crypto/bcrypt"
)

var verifyTimeLimit = 10 * time.Minute

const bcryptCost = 12

func SignUp(email string, metaData map[string]string, ipAddress string) error {
	err := approveSignUp(email, ipAddress)
	if err != nil {
		return err
	}

	emailHash, secret, err := saveChallenge(email, metaData, false, ipAddress)
	if err != nil {
		return err
	}

	err = sendSignUpChallenge(email, emailHash, secret)
	if err != nil {
		return err
	}

	return nil
}

// GetSecretFromEmail is a helper that can be used to extract the secret from the plaintext email body
// of a challenge email. This is only likely to be used by tests.
func GetSecretFromEmail(plainContent string) (string, error) {
	ss := strings.Split(plainContent, "type this phrase into your browser: ")
	if len(ss) != 2 {
		return "", errors.New(errors.Internal, "could not get secret from email body")
	}

	ss = strings.Split(ss[1], "\n")
	return ss[0], nil
}

// approveSignUp performs business-level validation for the Sign Up action.
func approveSignUp(email string, ipAddress string) error {
	err := userMustNotExist(email)
	if err != nil {
		return err
	}

	err = emailMustNotBeNoContact(email)
	if err != nil {
		return err
	}

	err = emailMustNotBeFlagged(email)
	if err != nil {
		return err
	}

	err = ipMustNotBeFlagged(ipAddress)
	if err != nil {
		return err
	}

	err = maxChallengesMustNotBeExceeded(email)
	if err != nil {
		return err
	}

	return nil
}

func saveChallenge(email string, metadata map[string]string, rememberMe bool, ipAddress string) (string, string, error) {
	s := secret.New()
	secretHash, err := bcrypt.GenerateFromPassword([]byte(s), bcryptCost)
	if err != nil {
		return "", "", &errors.Error{Code: errors.Internal, Err: err}
	}
	emailHash, err := bcrypt.GenerateFromPassword([]byte(email), bcryptCost)
	if err != nil {
		return "", "", &errors.Error{Code: errors.Internal, Err: err}
	}

	err = challenge.Create(&challenge.Challenge{
		Email:      email,
		EmailHash:  string(emailHash),
		SecretHash: string(secretHash),
		Metadata:   metadata,
		RememberMe: rememberMe,
		IPAddress:  ipAddress,
		ExpiresAt:  time.Now().Add(verifyTimeLimit),
	})
	if err != nil {
		return "", "", err
	}

	return string(emailHash), s, nil
}

func sendSignUpChallenge(emailAddress, emailHash, secret string) error {
	values := templates.SignUp{
		VerifyLink: verifyLink(emailHash, secret),
		Secret:     secret,
		NotYouLink: notYouLink(emailHash, secret),
	}

	return email.SendFromTemplate("Verify your Escapehatch account", emailAddress, "services/auth/templates/signUp.txt", "services/auth/templates/signUp.html", values)
}

func userMustNotExist(email string) error {
	exists, err := users.Exists(email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New(errors.EmailUnavailable, "email address is taken")
	}

	return nil
}

func emailMustNotBeNoContact(email string) error {
	n := nocontact.NoContact{}
	err := nocontact.GetByEmail(email, &n)
	if errors.Code(err) == errors.NotFound {
		return nil
	}
	if err != nil {
		return err
	}

	return errors.New(errors.DoNotContact, "email address is on do-not-contact list")
}

func emailMustNotBeFlagged(email string) error {
	f := flaggedemail.FlaggedEmail{}
	err := flaggedemail.GetUnexpiredByEmail(email, &f)
	if errors.Code(err) == errors.NotFound {
		return nil
	}
	if err != nil {
		return err
	}

	return errors.New(errors.EmailFlagged, "email address is temporarily flagged, please try again later")
}

func ipMustNotBeFlagged(ipAddress string) error {
	f := flaggedip.FlaggedIP{}
	err := flaggedip.GetUnexpiredByIP(ipAddress, &f)
	if errors.Code(err) == errors.NotFound {
		return nil
	}
	if err != nil {
		return err
	}

	return errors.New(errors.IPFlagged, "ip address is temporarily flagged, please try again later")
}

func maxChallengesMustNotBeExceeded(email string) error {
	challenges := []challenge.Challenge{}
	err := challenge.GetUnverifiedByEmailAndDate(email, time.Now().Add(-30*time.Minute), challenges)
	if err != nil {
		return err
	}
	if len(challenges) >= 3 {
		err := flaggedemail.Create(&flaggedemail.FlaggedEmail{
			Email:     email,
			ExpiresAt: time.Now().Add(30 * time.Minute),
		})
		if err != nil {
			return err
		}
		return errors.New(errors.EmailFlagged, "too many attempts, please try again later")
	}

	return nil
}

func verifyLink(emailHash, secret string) string {
	u, err := url.Parse(os.Getenv("FRONTEND_HOST"))
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "verify")
	q := u.Query()
	q.Set("emailHash", emailHash)
	q.Set("secret", secret)
	u.RawQuery = q.Encode()
	return u.String()
}

func notYouLink(emailHash, secret string) string {
	u, err := url.Parse(os.Getenv("FRONTEND_HOST"))
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "not-you")
	q := u.Query()
	q.Set("emailHash", emailHash)
	q.Set("secret", secret)
	u.RawQuery = q.Encode()
	return u.String()
}
