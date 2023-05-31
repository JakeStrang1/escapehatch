package items

import "github.com/JakeStrang1/escapehatch/db"

func CreateBook(userID string, result *Book) error {
	result.MediaType = MediaTypeBook
	result.CreatedBy = userID
	err := result.ValidateBookOnCreate()
	if err != nil {
		return err
	}
	return db.Create(result)
}
