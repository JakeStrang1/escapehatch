package storage

import (
	"fmt"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	// "google.golang.org/appengine"
	// "google.golang.org/appengine/file"
)

var client Client

type Client interface {
	Upload(filename string, data []byte) (string, error)
}

func SetupGCSClient() {
	client = NewGCSClient()
}

func SetupMockClient() {
	client = NewMockClient()
}

func Upload(filename string, data []byte) (string, error) {
	return client.Upload(filename, data)
}

func UploadFromURL(url string) (string, error) {
	return "", &errors.Error{Code: errors.Internal, Err: fmt.Errorf("not implemented")}
}
