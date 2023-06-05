package storage

import (
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/internal/images"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

var client Client

type Client interface {
	Upload(filename string, data []byte, options ...Options) (string, error)
	FileExists(filename string) (bool, error)
	Close()
}

type Options struct {
	Public         *bool
	ImageCompress  *bool // If true, must have a non-zero value for ImageMaxWidth, ImageMaxHeight, and ImageMaxKB
	ImageMaxWidth  *int
	ImageMaxHeight *int
	ImageMaxKB     *int
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

	if len(options) > 0 && lo.FromPtr(options[0].ImageCompress) {
		var err error
		data, err = images.CompressedJPEG(data, *options[0].ImageMaxWidth, *options[0].ImageMaxHeight, *options[0].ImageMaxKB)
		if err != nil {
			return "", err
		}
		ext = ".jpg"
	}

	newFilename := u.String() + ext
	return Upload(newFilename, data, options...)
}

// CreateFromURL will download the file at the url, and reupload it with a new filename to avoid any name conflicts
// If the url matches the name of an existing file, no changes will be made and the file url will be returned with no error.
func CreateFromURL(fileURL string, options ...Options) (string, error) {
	if exists, err := FileExists(fileURL); err != nil {
		return "", err
	} else if exists {
		return fileURL, nil
	}

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

	return Create(ext, data, options...)
}

// Upload will upload the given file with the given filename. It will overwrite an existing file with that name.
// Recommended to use Create instead to avoid overwriting.
func Upload(filename string, data []byte, options ...Options) (string, error) {
	return client.Upload(filename, data, options...)
}

func FileExists(fileURL string) (bool, error) {
	return client.FileExists(fileURL)
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
