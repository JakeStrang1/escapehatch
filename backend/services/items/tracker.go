package items

import "github.com/JakeStrang1/escapehatch/db"

/***************************************************
 * Change
 ***************************************************/

type Change struct {
	db.DefaultModel `db:",inline"`
	ItemID          string         `db:"item_id"`
	ChangeType      ChangeType     `db:"change_type"`
	UpdatedBy       string         `db:"updated_by"`
	Data            any            `db:"data"`     // The data being changed
	Metadata        map[string]any `db:"metadata"` // Any additional data about the change itself
}

func (c *Change) CollectionName() string {
	return "item_changes"
}

/***************************************************
 * Tracker
 ***************************************************/

type Tracker struct {
	ItemID     string
	ChangeType ChangeType
	UpdatedBy  string
	Update     any
	Result     any
	Metadata   map[string]any
}

func (t *Tracker) Created(result any) *Tracker {
	t.ChangeType = ChangeTypeCreate
	t.Result = result
	return t
}

func (t *Tracker) Updated(update any, result any) *Tracker {
	t.ChangeType = ChangeTypeUpdate
	t.Update = update
	t.Result = result
	return t
}

func (t *Tracker) Deleted(reason string, userCount int) *Tracker {
	t.ChangeType = ChangeTypeDelete
	if t.Metadata == nil {
		t.Metadata = map[string]any{}
	}
	t.Metadata["reason"] = reason
	t.Metadata["users_impacted"] = userCount
	return t
}

func (t *Tracker) By(userID string) *Tracker {
	t.UpdatedBy = userID
	return t
}

func (t *Tracker) Save() error {
	if t.ItemID == "" {
		panic("missing item id")
	}

	if t.UpdatedBy == "" {
		panic("missing updated by")
	}

	var data any
	switch t.ChangeType {
	case ChangeTypeCreate:
		data = t.Result
	case ChangeTypeUpdate:
		data = t.getUpdateData()
	case ChangeTypeDelete:
		data = nil
	default:
		panic("missing change type")
	}

	change := Change{
		ItemID:     t.ItemID,
		ChangeType: t.ChangeType,
		UpdatedBy:  t.UpdatedBy,
		Data:       data,
		Metadata:   t.Metadata,
	}
	return db.Create(&change)
}

func (t *Tracker) getUpdateData() any {
	var update any
	switch v := t.Update.(type) {
	case BookUpdate:
		bookResult := t.Result.(*Book)
		if v.ImageURL != nil || len(v.ImageFileBody) > 0 {
			v.ImageURL = &bookResult.ImageURL
		}
		update = v
	case MovieUpdate:
		movieResult := t.Result.(*Movie)
		if v.ImageURL != nil || len(v.ImageFileBody) > 0 {
			v.ImageURL = &movieResult.ImageURL
		}
		update = v
	default:
		panic("unknown type")
	}
	return update
}

func Track(itemID string) *Tracker {
	return &Tracker{ItemID: itemID}
}
