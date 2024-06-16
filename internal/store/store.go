package store

import (
	"context"

	"keeper-project/types"
)

type User interface {
	CreateUser(ctx context.Context, user *types.User) error
	GetByLogin(ctx context.Context, login string) (*types.User, error)
}

type Secrets[T any] interface {
	Create(context.Context, string, string, *T) error
	Get(context.Context, string, string) (*T, error)
	GetKeysList(context.Context, string) ([]types.Key, error)
	Update(context.Context, string, string, *T) error
	Delete(context.Context, string, string) error
}

type FileService interface {
	GetFile(ctx context.Context, bucketName, fileName string) (f *types.File, err error)
	GetFilesList(ctx context.Context, bucketName string) ([]*types.Key, error)
	Create(ctx context.Context, bucketName string, dto types.CreateFileDTO) error
	Delete(ctx context.Context, bucketName, fileName string) error
}
