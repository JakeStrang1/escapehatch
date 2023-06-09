package items

import (
	"fmt"

	"github.com/JakeStrang1/escapehatch/internal/errors"
)

type Filter struct {
	Page    *int `db:"-"`
	PerPage *int `db:"-"`
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
