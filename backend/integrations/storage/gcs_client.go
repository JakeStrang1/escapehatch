package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/JakeStrang1/escapehatch/internal/errors"
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

	// // Public - this code didn't work, leaving it here just as a reminder that it didn't work (no error, permissions didn't change).
	// // I ended up just making the entire bucket public instead.
	// if lo.FromPtr(opt.Public) {
	// 	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
	// 		return "", &errors.Error{Code: errors.Internal, Err: err}
	// 	}
	// }

	// Upload
	w := obj.NewWriter(ctx)
	_, err := w.Write(data)
	if err != nil {
		return "", &errors.Error{Code: errors.Internal, Err: err}
	}
	if err := w.Close(); err != nil {
		return "", &errors.Error{Code: errors.Internal, Err: err}
	}

	return filename, nil
}

func (g *GCSClient) Close() {
	g.Client.Close()
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
