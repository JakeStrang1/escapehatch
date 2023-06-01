package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/samber/lo"
)

type GCSClient struct {
	*storage.Client
	bucketName string
}

func (g *GCSClient) Upload(filename string, data []byte, options ...Options) (string, error) {
	opt := Options{}
	if len(options) > 1 {
		opt = options[0] // ignore additional options
	}

	ctx := context.Background()
	obj := g.Bucket(g.bucketName).Object(filename)
	w := obj.NewWriter(ctx)

	// Upload
	_, err := w.Write(data)
	if err != nil {
		return "", &errors.Error{Code: errors.Internal, Err: err}
	}
	if err := w.Close(); err != nil {
		return "", &errors.Error{Code: errors.Internal, Err: err}
	}

	// Public
	if lo.FromPtr(opt.Public) {
		if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return "", &errors.Error{Code: errors.Internal, Err: err}
		}
	}

	return filename, nil
}

func NewGCSClient(bucketName string) *GCSClient {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		fmt.Println("error during new client")
		panic(err)
	}
	return &GCSClient{
		Client:     client,
		bucketName: bucketName,
	}
}
