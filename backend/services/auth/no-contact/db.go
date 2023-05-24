package nocontact

import (
	"github.com/JakeStrang1/escapehatch/db"
)

const collection = "auth_no_contacts"

func Create(document *NoContact) error {
	err := db.EnsureUniqueIndex(&NoContact{}, []string{"email"})
	if err != nil {
		return err
	}
	return db.Create(document)
}

func GetByEmail(email string, result *NoContact) error {
	return db.GetOne(db.M{"email": email}, result)
}
