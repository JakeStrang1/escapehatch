package db

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/kamva/mgm/v3"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/****************************************************************************************
 * db.go
 *
 * This file is intended to:
 * - Setup the DB connection
 * - Define helpers for interacting with the DB
 ****************************************************************************************/

func Setup(host string, dbName string) error {
	err := mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI(host), &options.ClientOptions{
		Registry: Registry(), // Override the default registry to use "db" struct tag instead of "bson"
	})
	if err != nil {
		return &errors.Error{Code: errors.Internal, Message: "couldn't connect to database", Err: err}
	}
	return nil
}

func Close() {
	// Nothing to do
}

func Create(document mgm.Model) error {
	err := mgm.Coll(document).Create(document)
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}

func Delete(document mgm.Model) error {
	err := mgm.Coll(document).Delete(document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &errors.Error{Code: errors.NotFound, Err: err}
	}
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}

// GetByID returns the document with matching ID
func GetByID(document mgm.Model) error {
	err := mgm.Coll(document).FindByID(document.GetID(), document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &errors.Error{Code: errors.NotFound, Err: err}
	}
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}

type GetOneOptions struct {
	Sort [][]any
}

func (o *GetOneOptions) ToFindOneOptions() *options.FindOneOptions {
	result := &options.FindOneOptions{}
	if o.Sort != nil {
		d := bson.D{}
		for _, key := range o.Sort {
			d = append(d, bson.E{Key: key[0].(string), Value: key[1]})
		}
		result.Sort = d
	}
	return result
}

// GetOne returns the first result of the selector query
func GetOne(selector map[string]interface{}, result mgm.Model, opts ...GetOneOptions) error {
	if len(opts) > 1 {
		return &errors.Error{Code: errors.Internal, Err: fmt.Errorf("just don't")}
	}

	findOpts := &options.FindOneOptions{}
	if len(opts) == 1 {
		findOpts = opts[0].ToFindOneOptions()
	}

	err := mgm.Coll(result).First(selector, result, findOpts)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &errors.Error{Code: errors.NotFound, Err: err}
	}
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}

type GetManyOptions struct {
	Limit *int64
	Skip  *int64
}

func (o *GetManyOptions) ToFindOptions() *options.FindOptions {
	result := &options.FindOptions{}
	if o.Limit != nil {
		result.Limit = o.Limit
	}
	if o.Skip != nil {
		result.Skip = o.Skip
	}
	return result
}

// GetMany returns all results of the selector query
func GetMany(selector interface{}, model mgm.Model, results interface{}, opts ...GetManyOptions) error {
	if len(opts) > 1 {
		return &errors.Error{Code: errors.Internal, Err: fmt.Errorf("just don't")}
	}
	findOpts := &options.FindOptions{}
	if len(opts) == 1 {
		findOpts = opts[0].ToFindOptions()
	}

	err := mgm.Coll(model).SimpleFind(results, selector, findOpts)
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}

func GetPage(selector interface{}, model mgm.Model, results interface{}, page int, perPage int, opts ...GetManyOptions) (int, error) {
	if len(opts) == 0 {
		opts = append(opts, GetManyOptions{})
	}
	opts[0].Limit = lo.ToPtr(int64(perPage))
	opts[0].Skip = lo.ToPtr(int64((page - 1) * perPage))

	err := GetMany(selector, model, results, opts...)
	if err != nil {
		return 0, err
	}

	// Count total pages
	count, err := mgm.Coll(model).CountDocuments(mgm.Ctx(), selector)
	if err != nil {
		return 0, &errors.Error{Code: errors.Internal, Err: err}
	}
	totalPages := int(math.Ceil(float64(count) / float64(perPage)))

	return totalPages, nil
}

func GetCount(selector interface{}, model mgm.Model) (int, error) {
	count, err := mgm.Coll(model).CountDocuments(mgm.Ctx(), selector)
	if err != nil {
		return 0, &errors.Error{Code: errors.Internal, Err: err}
	}
	return int(count), nil
}

func EnsureUniqueIndex(model mgm.Model, keys []string) error {
	d := bson.D{}
	for _, key := range keys {
		d = append(d, bson.E{Key: key, Value: 1})
	}
	_, err := mgm.Coll(model).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: d,
		Options: &options.IndexOptions{
			Unique: lo.ToPtr(true),
		},
	})
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}

func EnsureOptionalUniqueIndex(model mgm.Model, key string) error {
	d := bson.D{{Key: key, Value: 1}}
	_, err := mgm.Coll(model).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: d,
		Options: &options.IndexOptions{
			Unique:                  lo.ToPtr(true),
			PartialFilterExpression: bson.M{key: bson.M{"$exists": true, "$gt": ""}},
		},
	})
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}

// EnsureTextIndex creates a text index used for a later search query
// Example search query: { $text: { $search: "\"coffee shop\"" } }
// Note this can only do full word watches, not partial matches!
// More info: https://www.mongodb.com/docs/manual/core/link-text-indexes
func EnsureTextIndex(model mgm.Model, keys []string) error {
	d := bson.D{}
	for _, key := range keys {
		d = append(d, bson.E{Key: key, Value: "text"})
	}
	_, err := mgm.Coll(model).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: d,
	})
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}

func IncrementOne(selector map[string]interface{}, fieldName string, result mgm.Model) error {
	update := bson.M{"$inc": bson.M{fieldName: 1}, "$set": bson.M{"updated_at": time.Now()}}
	err := mgm.Coll(result).FindOneAndUpdate(context.Background(), selector, update, &options.FindOneAndUpdateOptions{
		ReturnDocument: lo.ToPtr(options.After),
	}).Decode(result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &errors.Error{Code: errors.NotFound, Err: err}
	}
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}

func Update(result mgm.Model) error {
	err := mgm.Coll(result).Update(result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &errors.Error{Code: errors.NotFound, Err: err}
	}
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}
	return nil
}
