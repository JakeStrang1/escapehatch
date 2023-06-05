package db

import (
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

/****************************************************************************************
 * atlas.go
 *
 * This file is intended to:
 * - Contain db calls specific to MongoDB Atlas (i.e. won't work on local MongoDB instances)
 ****************************************************************************************/

func SearchAtlas(search string, paths []string, model mgm.Model, results any) error {
	pipeline := []bson.M{
		{"$search": bson.M{
			"text": bson.M{
				"query": search,
				"path":  []string{"username", "full_name"},
			},
		}},
		{"$addFields": bson.M{"search_score": bson.M{"$meta": "searchScore"}}},
	}
	cursor, err := mgm.Coll(model).Aggregate(mgm.Ctx(), pipeline)
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}

	err = cursor.All(mgm.Ctx(), results)
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}
