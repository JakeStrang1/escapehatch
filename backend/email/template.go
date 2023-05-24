package email

import (
	"bytes"
	html "html/template"
	"path"
	fp "path/filepath"
	text "text/template"

	"github.com/JakeStrang1/escapehatch/internal/errors"
)

func sendFromTemplate(subject, toAddress, textTemplatePath, htmlTemplatePath string, templateValues interface{}) error {
	plainContent, err := executeTextTemplate(textTemplatePath, templateValues)
	if err != nil {
		return err
	}

	htmlContent, err := executeHTMLTemplate(htmlTemplatePath, templateValues)
	if err != nil {
		return err
	}

	return mailer.Send(SendParams{
		Subject:      subject,
		To:           toAddress,
		PlainContent: plainContent,
		HTMLContent:  htmlContent,
	})
}

func executeTextTemplate(filepath string, values interface{}) (string, error) {
	filename := path.Base(filepath)
	t, err := text.New(filename).ParseFiles(fp.FromSlash(filepath))
	if err != nil {
		return "", &errors.Error{Code: errors.Internal, Err: err}
	}

	b := new(bytes.Buffer)
	err = t.Execute(b, values)
	if err != nil {
		return "", &errors.Error{Code: errors.Internal, Err: err}
	}
	return b.String(), nil
}

func executeHTMLTemplate(filepath string, values interface{}) (string, error) {
	filename := path.Base(filepath)
	t, err := html.New(filename).ParseFiles(fp.FromSlash(filepath))
	if err != nil {
		return "", &errors.Error{Code: errors.Internal, Err: err}
	}

	b := new(bytes.Buffer)
	err = t.Execute(b, values)
	if err != nil {
		return "", &errors.Error{Code: errors.Internal, Err: err}
	}
	return b.String(), nil
}
