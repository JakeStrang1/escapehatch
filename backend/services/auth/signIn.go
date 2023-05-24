package auth

import (
	"github.com/JakeStrang1/saas-template/email"
	"github.com/JakeStrang1/saas-template/internal/errors"
	"github.com/JakeStrang1/saas-template/services/auth/templates"
	"github.com/JakeStrang1/saas-template/services/users"
)

func SignIn(email string, metaData map[string]string, rememberMe bool, ipAddress string) error {
	err := approveSignIn(email, ipAddress)
	if err != nil {
		return err
	}

	emailHash, secret, err := saveChallenge(email, metaData, rememberMe, ipAddress)
	if err != nil {
		return err
	}

	err = sendSignInChallenge(email, emailHash, secret)
	if err != nil {
		return err
	}

	return nil
}

// approveSignIn performs business-level validation for the Sign In action.
func approveSignIn(email string, ipAddress string) error {
	err := userMustExist(email)
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

func sendSignInChallenge(emailAddress, emailHash, secret string) error {
	values := templates.SignIn{
		VerifyLink: verifyLink(emailHash, secret),
		Secret:     secret,
		NotYouLink: notYouLink(emailHash, secret),
	}

	return email.SendFromTemplate("Sign In to SaaS Template", emailAddress, "services/auth/templates/signIn.txt", "services/auth/templates/signIn.html", values)
}

func userMustExist(email string) error {
	exists, err := users.Exists(email)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(errors.EmailInvalid, "email address not registered")
	}

	return nil
}
