package minio

import (
	"context"
	"io"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	Client         *minio.Client
	Bucket         string
	PublicEndpoint string
}

func NewMinioService(
	endpoint string,
	PublicEndpoint string,
	accessKey string,
	secretKey string,
	bucket string,
) (*MinioService, error) {

	client, err := minio.New(PublicEndpoint, &minio.Options{
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
		Client:         client,
		Bucket:         bucket,
		PublicEndpoint: PublicEndpoint,
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
	u, err := s.Client.PresignedGetObject(ctx, s.Bucket, objectName, expiry, url.Values{})
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (s *MinioService) RefreshArticleImages(ctx context.Context, article string) string {
	re := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)
	count := 0

	return re.ReplaceAllStringFunc(article, func(tag string) string {
		reSrc := regexp.MustCompile(`src="([^"]+)"`)
		match := reSrc.FindStringSubmatch(tag)
		if len(match) < 2 {
			return tag
		}

		oldURL := match[1]
		newURL, err := s.GetURL(ctx, oldURL, 7*24*time.Hour)
		if err != nil {
			log.Printf("[ERROR] Failed to refresh image #%d: %v", count+1, err)
			return tag
		}

		count++
		log.Printf("[SUCCESS] Image #%d refreshed", count)
		return strings.Replace(tag, oldURL, newURL, 1)
	})
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
