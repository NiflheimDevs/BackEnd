package storage

import (
	"time"

	"github.com/niflheimdevs/backend/internal/domain/enums"
)

type S3Storage interface {
	UploadObject(bucketType enums.BucketType, key string, file []byte)
	DeleteObject(bucketType enums.BucketType, key string) error
	GetPresignedURL(bucketType enums.BucketType, objectKey string, expiration time.Duration) string
}
