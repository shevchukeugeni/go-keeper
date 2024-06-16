package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"go.uber.org/zap"

	"keeper-project/internal/store/file"
	"keeper-project/pkg/minio"
	"keeper-project/types"
)

type minioStorage struct {
	client *minio.Client
}

func NewStorage(logger *zap.Logger, endpoint, accessKeyID, secretAccessKey string) (file.Storage, error) {
	client, err := minio.NewClient(endpoint, accessKeyID, secretAccessKey, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}
	return &minioStorage{
		client: client,
	}, nil
}

func (m *minioStorage) GetFile(ctx context.Context, bucketName, fileID string) (*types.File, error) {
	obj, err := m.client.GetFile(ctx, bucketName, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file. err: %w", err)
	}
	defer obj.Close()
	objectInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file. err: %w", err)
	}
	buffer := make([]byte, objectInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to get objects. err: %w", err)
	}
	f := types.File{
		ID:       objectInfo.Key,
		Name:     objectInfo.UserMetadata["Name"],
		Size:     objectInfo.Size,
		Bytes:    buffer,
		Metadata: objectInfo.UserMetadata["Metadata"],
	}

	return &f, nil
}

func (m *minioStorage) GetFilesList(ctx context.Context, bucketName string) ([]*types.Key, error) {
	objs, err := m.client.GetFilesList(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get file. err: %w", err)
	}

	var keys []*types.Key

	for _, obj := range objs {
		keys = append(keys, &types.Key{
			Id:  obj.ID,
			Key: obj.Name,
		})
	}

	return keys, nil
}

func (m *minioStorage) CreateFile(ctx context.Context, bucketName string, file *types.File) error {
	err := m.client.UploadFile(ctx, file.ID, file.Name, bucketName, file.Metadata, file.Size, bytes.NewBuffer(file.Bytes))
	if err != nil {
		return err
	}
	return nil
}

func (m *minioStorage) DeleteFile(ctx context.Context, bucketName, fileId string) error {
	err := m.client.DeleteFile(ctx, bucketName, fileId)
	if err != nil {
		return err
	}
	return nil
}
