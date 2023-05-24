package http

import (
	"time"

	"github.com/JakeStrang1/saas-template/db"
	"github.com/samber/lo"
)

type DefaultModelAPI struct {
	ID        *string    `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func ToDefaultModelAPI(model db.DefaultModel) DefaultModelAPI {
	return DefaultModelAPI{
		ID:        &model.ID,
		CreatedAt: &model.CreatedAt,
		UpdatedAt: &model.UpdatedAt,
	}
}

func FromDefaultModelAPI(model DefaultModelAPI) db.DefaultModel {
	return db.DefaultModel{
		IDField: db.IDField{ID: lo.FromPtr(model.ID)},
		DateFields: db.DateFields{
			CreatedAt: lo.FromPtr(model.CreatedAt),
			UpdatedAt: lo.FromPtr(model.UpdatedAt),
		},
	}
}
