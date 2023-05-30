package db

import "go.mongodb.org/mongo-driver/bson"

type M map[string]interface{}

type Marshaler bson.Marshaler

var Marshal func(interface{}) ([]byte, error) = bson.Marshal
