package users

import (
	"fmt"
	"regexp"

	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/errors"
)

type Filter struct {
	Page    *int    `db:"-"`
	PerPage *int    `db:"-"`
	Search  *string `db:"-"`
}

// MarshalBSON satisfies the db.Marshaler interface
func (f Filter) MarshalBSON() ([]byte, error) { // Non-pointer receiver so that filter can be passed as non-pointer
	doc := db.M{}
	if f.Search != nil {
		regexEscapedSearch := regexp.QuoteMeta(*f.Search)
		doc["$or"] = []db.M{
			{"username": db.M{"$regex": regexEscapedSearch, "$options": "i"}},
			{"full_name": db.M{"$regex": regexEscapedSearch, "$options": "i"}},
		}
	}
	return db.Marshal(doc)
}

func (f *Filter) Validate() error {
	if f.Page != nil {
		if *f.Page < 1 {
			return &errors.Error{Code: errors.Invalid, Message: "page must be 1 or greater"}
		}
	}

	if f.PerPage != nil {
		if *f.PerPage < 1 || *f.PerPage > maxPerPage {
			return &errors.Error{Code: errors.Invalid, Message: fmt.Sprintf("page size must be between 1 and %d, received '%d'", maxPerPage, *f.PerPage)}
		}
	}

	return nil
}
