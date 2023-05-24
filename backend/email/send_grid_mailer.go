package email

import (
	"github.com/JakeStrang1/saas-template/internal/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	Client *sendgrid.Client
	From   *mail.Email
}

type SendGridConfig struct {
	APIKey      string
	FromAddress string
}

func NewSendGridMailer(config SendGridConfig) *SendGridMailer {
	return &SendGridMailer{
		Client: sendgrid.NewSendClient(config.APIKey),
		From:   mail.NewEmail("SaaS Template", config.FromAddress),
	}
}

func (s *SendGridMailer) Send(params SendParams) error {
	to := mail.NewEmail("", params.To)
	message := mail.NewSingleEmail(s.From, params.Subject, to, params.PlainContent, params.HTMLContent)
	_, err := s.Client.Send(message)
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}
