package nocontact

import "github.com/JakeStrang1/saas-template/db"

type NoContact struct {
	db.DefaultModel `db:",inline"`
	Email           string `db:"email"`
}

func (n *NoContact) CollectionName() string {
	return collection
}
