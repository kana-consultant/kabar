package minio

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	Client *minio.Client
	Bucket string
}

func NewMinioService(
	endpoint,
	accessKey,
	secretKey,
	bucket string,
) (*MinioService, error) {

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(
			ctx,
			bucket,
			minio.MakeBucketOptions{},
		)
		if err != nil {
			return nil, err
		}
	}

	return &MinioService{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (s *MinioService) Upload(
	ctx context.Context,
	objectName string,
	file io.Reader,
	size int64,
	contentType string,
) (string, error) {

	exists, err := s.Client.BucketExists(ctx, s.Bucket)
	if err != nil {
		return "", err
	}

	if !exists {
		err = s.Client.MakeBucket(ctx, s.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}
	}

	_, err = s.Client.PutObject(
		ctx,
		s.Bucket,
		objectName,
		file,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (s *MinioService) Delete(ctx context.Context, objectName string) error {
	return s.Client.RemoveObject(ctx, s.Bucket, objectName, minio.RemoveObjectOptions{})
}

func (s *MinioService) GetURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := s.Client.PresignedGetObject(ctx, s.Bucket, objectName, expiry, url.Values{})
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *MinioService) List(ctx context.Context, prefix string) ([]string, error) {
	var files []string

	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	for object := range s.Client.ListObjects(ctx, s.Bucket, opts) {
		if object.Err != nil {
			return nil, object.Err
		}
		files = append(files, object.Key)
	}

	return files, nil
}
