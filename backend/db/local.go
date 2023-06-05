package db

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

/****************************************************************************************
 * local.go
 *
 * This file is intended to:
 * - Contain db calls specific to local instances (i.e. an alternative exists for production)
 ****************************************************************************************/

func SearchLocal(search string, model mgm.Model, results any) error {
	selector := bson.M{"$text": bson.M{"$search": search}}
	return GetMany(selector, model, results)
}
