package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type Object struct {
	ID   string
	Size int64
	Name string
}

type Client struct {
	logger      *zap.Logger
	minioClient *minio.Client
}

func NewClient(endpoint, accessKeyID, secretAccessKey string, logger *zap.Logger) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}

	return &Client{
		logger:      logger,
		minioClient: minioClient,
	}, nil
}

func (c *Client) GetFile(ctx context.Context, bucketName, fileId string) (*minio.Object, error) {
	obj, err := c.minioClient.GetObject(ctx, bucketName, fileId, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file with id: %s from minio bucket %s. err: %w", fileId, bucketName, err)
	}

	return obj, nil
}

func (c *Client) GetFilesList(ctx context.Context, bucketName string) ([]*Object, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var list []*Object

	for lobj := range c.minioClient.ListObjects(reqCtx, bucketName, minio.ListObjectsOptions{WithMetadata: true}) {
		if lobj.Err != nil {
			c.logger.Error("failed to list object from minio bucket",
				zap.String("bucket name", bucketName), zap.Error(lobj.Err))
			continue
		}
		obj := new(Object)
		obj.ID = lobj.Key
		obj.Size = lobj.Size
		obj.Name = lobj.UserMetadata["X-Amz-Meta-Name"]
		list = append(list, obj)
	}

	return list, nil
}

func (c *Client) UploadFile(ctx context.Context, fileId, fileName, bucketName, metadata string, fileSize int64, reader io.Reader) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	exists, errBucketExists := c.minioClient.BucketExists(ctx, bucketName)
	if errBucketExists != nil || !exists {
		c.logger.Warn("no bucket with this key. creating new one...", zap.String("key", bucketName))
		err := c.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create new bucket. err: %w", err)
		}
	}

	c.logger.Debug("put new object to bucket", zap.String("file", fileName), zap.String("bucket", bucketName))
	_, err := c.minioClient.PutObject(reqCtx, bucketName, fileId, reader, fileSize,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name":     fileName,
				"Metadata": metadata,
			},
			ContentType: "application/octet-stream",
		})
	if err != nil {
		return fmt.Errorf("failed to upload file. err: %w", err)
	}
	return nil
}

func (c *Client) DeleteFile(ctx context.Context, bucketName, fileName string) error {
	err := c.minioClient.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file. err: %w", err)
	}
	return nil
}
