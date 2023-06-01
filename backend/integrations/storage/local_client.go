package storage

import (
	"io/fs"
	"time"

	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/errors"
)

type LocalClient struct {
}

func NewLocalClient() *LocalClient {
	return &LocalClient{}
}

func (l *LocalClient) Upload(filename string, data []byte, options ...Options) (string, error) {
	localFile := LocalFile{
		Filename: filename,
		Data:     data,
	}
	err := db.Create(&localFile)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (l *LocalClient) Close() {}

/******************************************************
 * FS implementation used to retrieve files
 ******************************************************/

// LocalFS implements fs.FS
// It uses the db to store and retrieve files
type LocalFS struct{}

// Open opens the named file
// Source: https://pkg.go.dev/io/fs#FS
func (l *LocalFS) Open(name string) (fs.File, error) {
	localFile := LocalFile{}
	err := db.GetOne(db.M{"filename": name}, &localFile)
	if errors.Code(err) == errors.NotFound {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}
	return &localFile, nil
}

// LocalFS implements fs.File
type LocalFile struct {
	db.DefaultModel `db:",inline"`
	Filename        string `db:"filename"`
	Data            []byte `db:"data"`
	offset          int    // Tracks the current read position of the file
}

func (l *LocalFile) Stat() (fs.FileInfo, error) {
	return &LocalFileInfo{
		name:      l.Filename,
		size:      int64(len(l.Data)),
		updatedAt: l.UpdatedAt,
	}, nil
}

func (l *LocalFile) Read(output []byte) (int, error) {
	for i := range output {
		if len(l.Data) <= i+l.offset {
			return i, nil
		}
		output[i] = l.Data[i+l.offset]
	}
	l.offset += len(output)
	return len(output), nil
}

func (l *LocalFile) Close() error {
	return nil
}

// LocalFS implements fs.FileInfo
type LocalFileInfo struct {
	name      string
	size      int64
	updatedAt time.Time
}

func (l *LocalFileInfo) Name() string {
	return l.name
}

func (l *LocalFileInfo) Size() int64 {
	return l.size
}

func (l *LocalFileInfo) Mode() fs.FileMode {
	return fs.ModePerm // All permissions
}

func (l *LocalFileInfo) ModTime() time.Time {
	return l.updatedAt
}

func (l *LocalFileInfo) IsDir() bool {
	return false
}

func (l *LocalFileInfo) Sys() any {
	return nil
}
