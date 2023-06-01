package items

import (
	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/integrations/storage"
)

func CreateBook(userID string, result *Book) error {
	result.MediaType = MediaTypeBook
	result.CreatedBy = userID
	err := result.ValidateBookOnCreate()
	if err != nil {
		return err
	}

	var newImageURL string
	if len(result.ImageFileBody) != 0 {
		newImageURL, err = storage.Create(result.ImageFileName, result.ImageFileBody)
	} else {
		newImageURL, err = storage.UploadFromURL(result.ImageURL)
	}
	if err != nil {
		return err
	}
	result.ImageURL = newImageURL
	return db.Create(result)
}
