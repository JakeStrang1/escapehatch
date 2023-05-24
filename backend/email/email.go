package email

var mailer Mailer

type Mailer interface {
	Send(SendParams) error
}

type SendParams struct {
	Subject      string
	To           string
	PlainContent string
	HTMLContent  string
}

func GetMailer() Mailer {
	return mailer
}

func SetupMockMailer() {
	mailer = NewMockMailer()
}

func SetupSendGridMailer(config SendGridConfig) {
	mailer = NewSendGridMailer(config)
}

// SendFromTemplate sends an email using the current Mailer
func SendFromTemplate(subject, toAddress, textTemplatePath, htmlTemplatePath string, templateValues interface{}) error {
	return sendFromTemplate(subject, toAddress, textTemplatePath, htmlTemplatePath, templateValues)
}
