package storage

import (
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
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
// oldFilename is discarded except for the file extension
func Create(oldFilename string, data []byte, options ...Options) (string, error) {
	ext := filepath.Ext(oldFilename)
	u := uuid.New()
	newFilename := u.String() + ext
	return Upload(newFilename, data, options...)
}

// CreateFromURL will download the file at the url, and reupload it with a new filename to avoid any name conflicts
func CreateFromURL(fileURL string, options ...Options) (string, error) {
	// Validate URL
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", &errors.Error{Code: errors.Invalid, Message: "url is invalid", Err: err}
	}

	// Download file
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", &errors.Error{Code: errors.Invalid, Message: "error fetching file from url", Err: err}
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", &errors.Error{Code: errors.Internal, Err: err}
	}

	// File extension
	ext := path.Ext(path.Base(u.Path))
	if ext == "" {
		ext = extensionFromResponse(resp)
	}

	return Create(ext, data)
}

// Upload will upload the given file with the given filename. It will overwrite an existing file with that name.
// Recommended to use Create instead to avoid overwriting.
func Upload(filename string, data []byte, options ...Options) (string, error) {
	return client.Upload(filename, data, options...)
}

func extensionFromResponse(resp *http.Response) string {
	mimeTypes := resp.Header.Values("Content-Type")
	for _, mimeType := range mimeTypes {
		exts, _ := mime.ExtensionsByType(mimeType)
		if len(exts) > 0 {
			return exts[0]
		}
	}
	return ""
}
