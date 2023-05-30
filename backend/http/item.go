package http

import (
	"github.com/JakeStrang1/escapehatch/services/items"
	"github.com/samber/lo"
)

type ItemAPI struct {
	DefaultModelAPI
	MediaType *string `json:"media_type"`
	Image     *string `json:"image"`
	Title     *string `json:"title"`
	CreatedBy *string `json:"created_by"`
}

func ToItemAPI(item items.Item) ItemAPI {
	return ItemAPI{
		DefaultModelAPI: ToDefaultModelAPI(item.DefaultModel),
		MediaType:       (*string)(&item.MediaType),
		Image:           &item.Image,
		Title:           &item.Title,
		CreatedBy:       &item.CreatedBy,
	}
}

func ToItem(itemAPI ItemAPI) items.Item {
	return items.Item{
		Image: lo.FromPtr(itemAPI.Image),
		Title: lo.FromPtr(itemAPI.Title),
	}
}
