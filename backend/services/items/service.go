package items

import (
	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/integrations/storage"
	"github.com/samber/lo"
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
		newImageURL, err = storage.Create(result.ImageFileName, result.ImageFileBody, storage.Options{Public: lo.ToPtr(true)})
	} else {
		newImageURL, err = storage.CreateFromURL(result.ImageURL)
	}
	if err != nil {
		return err
	}
	result.ImageURL = newImageURL
	return db.Create(result)
}

func CreateMovie(userID string, result *Movie) error {
	result.MediaType = MediaTypeMovie
	result.CreatedBy = userID
	err := result.ValidateMovieOnCreate()
	if err != nil {
		return err
	}

	var newImageURL string
	if len(result.ImageFileBody) != 0 {
		newImageURL, err = storage.Create(result.ImageFileName, result.ImageFileBody, storage.Options{Public: lo.ToPtr(true)})
	} else {
		newImageURL, err = storage.CreateFromURL(result.ImageURL)
	}
	if err != nil {
		return err
	}
	result.ImageURL = newImageURL
	return db.Create(result)
}
