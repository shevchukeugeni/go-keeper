package file

import (
	"context"

	"keeper-project/types"
)

type Storage interface {
	GetFile(ctx context.Context, bucketName, fileName string) (*types.File, error)
	GetFilesList(ctx context.Context, bucketName string) ([]*types.Key, error)
	CreateFile(ctx context.Context, bucketName string, file *types.File) error
	DeleteFile(ctx context.Context, bucketName, fileName string) error
}
