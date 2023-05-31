package storage

import (
	"context"
	"fmt"

	storage "cloud.google.com/go/storage"
)

type GCSClient struct {
	*storage.Client
}

func (g *GCSClient) Upload(filename string, data []byte) (string, error) {
	bkt := g.Bucket("escapehatch.appspot.com")
	obj := bkt.Object(filename)

	// Write something to obj.
	// w implements io.Writer.
	w := obj.NewWriter(context.Background())
	// Write some text to obj. This will either create the object or overwrite whatever is there already.
	_, err := w.Write(data)
	if err != nil {
		fmt.Println("error during file write")
		fmt.Println(err)
		// TODO: Handle error.
	}
	// Close, just like writing a file.
	if err := w.Close(); err != nil {
		fmt.Println("error during file close")
		fmt.Println(err)
	}

	return "test", nil
}

func NewGCSClient() *GCSClient {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		fmt.Println("error during new client")
		panic(err)
	}
	return &GCSClient{
		Client: client,
	}
}
