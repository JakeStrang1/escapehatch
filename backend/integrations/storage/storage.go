package storage

import (
	"fmt"
	"path/filepath"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/google/uuid"
)

var client Client

type Client interface {
	Upload(filename string, data []byte, options ...Options) (string, error)
	Close()
}

type Options struct {
	Public *bool
}

func SetupGCSClient(bucketName string) {
	client = NewGCSClient(bucketName)
}

func SetupLocalClient() {
	client = NewLocalClient()
}

func Close() {
	client.Close()
}

// Create will upload a file with a new filename to avoid any name conflicts
func Create(oldFilename string, data []byte, options ...Options) (string, error) {
	ext := filepath.Ext(oldFilename)
	u := uuid.New()
	newFilename := u.String() + ext
	return Upload(newFilename, data, options...)
}

// Upload will upload the given file with the given filename. It will overwrite an existing file with that name.
// Recommended to use Create instead to avoid overwriting.
func Upload(filename string, data []byte, options ...Options) (string, error) {
	return client.Upload(filename, data, options...)
}

func UploadFromURL(url string, options ...Options) (string, error) {
	return "", &errors.Error{Code: errors.Internal, Err: fmt.Errorf("not implemented")}
}
