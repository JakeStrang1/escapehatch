package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*********************
 * Override mgm DefaultModel (based on: https://github.com/Kamva/mgm/blob/81ef3e8616c782e2f7f025233e9fa5102c44df6b/model.go#L30-L33)
 *********************/

// DefaultModel struct contains a model's default fields.
type DefaultModel struct {
	IDField    `db:",inline"`
	DateFields `db:",inline"`
}

// Creating function calls the inner fields' defined hooks
// TODO: get context as param in the next version (4).
func (model *DefaultModel) Creating() error {
	model.IDField.Creating()
	model.DateFields.Creating()
	return nil
}

// Saving function calls the inner fields' defined hooks
// TODO: get context as param the next version(4).
func (model *DefaultModel) Saving() error {
	model.DateFields.Saving()
	return nil
}

// IDField struct contains a model's ID field.
type IDField struct {
	ID string `db:"_id,omitempty"`
}

// DateFields struct contains the `created_at` and `updated_at`
// fields that autofill when inserting or updating a model.
type DateFields struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// PrepareID method prepares the ID value to be used for filtering
// e.g convert bson.ObjectId value to hex-string
func (f *IDField) PrepareID(id interface{}) (interface{}, error) {
	if idObjectID, ok := id.(primitive.ObjectID); ok {
		return idObjectID.Hex(), nil
	}

	// Otherwise id must be string
	return id, nil
}

// GetID method returns a model's ID
func (f *IDField) GetID() interface{} {
	return f.ID
}

// SetID sets the value of a model's ID field.
func (f *IDField) SetID(id interface{}) {
	if idObjectID, ok := id.(primitive.ObjectID); ok {
		f.ID = idObjectID.Hex()
	} else if idString, ok := id.(string); ok {
		f.ID = idString
	}
}

//--------------------------------
// DateField methods
//--------------------------------

// Creating hook is used here to set the `created_at` field
// value when inserting a new model into the database.
// TODO: get context as param the next version(4).
func (f *DateFields) Creating() {
	f.CreatedAt = time.Now().UTC()
	f.UpdatedAt = time.Now().UTC()
}

// Saving hook is used here to set the `updated_at` field
// value when creating or updating a model.
// TODO: get context as param the next version(4).
func (f *DateFields) Saving() {
	f.UpdatedAt = time.Now().UTC()
}

//--------------------------------
// IDField methods
//--------------------------------

// Creating hook is used here to set the `_id` field
// value when inserting a new model into the database.
// TODO: get context as param the next version(4).
func (f *IDField) Creating() {
	f.ID = primitive.NewObjectID().Hex()
}
