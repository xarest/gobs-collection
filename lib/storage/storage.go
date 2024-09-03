package storage

import "context"

type IStorage interface {
	IsExisted(ctx context.Context, key string) (bool, error)
	FetchMetadata(c context.Context, key string) (map[string]string, error)
	CopyObject(c context.Context, sourceKey string, destKey string, metadata map[string]string) error
	DeleteObject(c context.Context, key string) error
	CommitFileUploaded(c context.Context, key string, parent string, target string, prefix FilePrefix) (string, func(*string) error, func() error, error)
}

func NewIStorage() IStorage {
	return &S3{}
}
