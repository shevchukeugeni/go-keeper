package file

import (
	"context"

	"go.uber.org/zap"

	"keeper-project/internal/store"
	"keeper-project/types"
)

var _ store.FileService = &service{}

type service struct {
	storage Storage
	logger  *zap.Logger
}

func NewService(fileStorage Storage, logger *zap.Logger) (store.FileService, error) {
	return &service{
		storage: fileStorage,
		logger:  logger,
	}, nil
}

func (s *service) GetFile(ctx context.Context, bucketName, fileId string) (*types.File, error) {
	ret, err := s.storage.GetFile(ctx, bucketName, fileId)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *service) GetFilesList(ctx context.Context, bucketName string) ([]*types.Key, error) {
	list, err := s.storage.GetFilesList(ctx, bucketName)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (s *service) Create(ctx context.Context, bucketName string, dto types.CreateFileDTO) error {
	dto.NormalizeName()
	file, err := types.NewFile(dto)
	if err != nil {
		return err
	}
	err = s.storage.CreateFile(ctx, bucketName, file)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Delete(ctx context.Context, bucketName, fileName string) error {
	err := s.storage.DeleteFile(ctx, bucketName, fileName)
	if err != nil {
		return err
	}
	return nil
}
