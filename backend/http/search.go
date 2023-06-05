package http

import (
	"github.com/JakeStrang1/escapehatch/internal/pages"
	"github.com/JakeStrang1/escapehatch/services/items"
	"github.com/JakeStrang1/escapehatch/services/users"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	userID := c.GetString(CtxKeyUserID)
	search := c.Query("search")

	results := []any{}
	err := items.GetSearch(search, &results)
	if err != nil {
		Error(c, err)
		return
	}

	pageInfo := pages.PageResult{
		Page:       1,
		PerPage:    len(results),
		TotalPages: 1,
	}

	ReturnMany(c, ToSearchableAPIs(userID, results), pageInfo)
}

func ToSearchableAPIs(selfID string, docs []any) []any {
	results := []any{}
	for _, doc := range docs {
		switch v := doc.(type) {
		case users.User:
			results = append(results, ToUserAPI(selfID, v))
		case items.Book:
			results = append(results, ToBookAPI(v))
		case items.Movie:
			results = append(results, ToMovieAPI(v))
		case items.TVSeries:
			results = append(results, ToTVSeriesAPI(v))
		default:
			panic("unknown type")
		}
	}
	return results
}
