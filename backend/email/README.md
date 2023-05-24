# email package

Used to send emails from the app.

## Highlights

- To send email, call `email.SendFromTemplate("Subject", "to.user@gmail.com", "path/to/template.txt", "path/to/template.html", templateValues)`

## Files

- `email.go` is the entry point to the package.
- `mock_mailer.go` defines a mock implementation of the mailer which is used for tests.
- `send_grid_mailer.go` defines a SendGrid mailer which is used for sending real emails.
- `template.go` has implementation for template-based emails.