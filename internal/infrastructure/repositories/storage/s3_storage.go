package storageimpl

import (
	"bytes"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/domain/enums"
)

type S3Storage struct {
	env      *bootstrap.Env
	storage  *bootstrap.S3
	clients  *s3.S3
	uploader *s3manager.Uploader
	buckets  map[enums.BucketType]string
}

func NewS3Storage(
	env *bootstrap.Env,
	storage *bootstrap.S3,
) *S3Storage {
	buckets := make(map[enums.BucketType]string)
	buckets[enums.ProfilePic] = storage.Buckets.ProfilePic
	return &S3Storage{
		env:     env,
		storage: storage,
		buckets: buckets,
	}
}

func (s3StorageS3Storage *S3Storage) setS3Client(bucketType enums.BucketType) {
	bucketTypes := enums.GetAllBucketTypes()
	if !slices.Contains(bucketTypes, bucketType) {
		panic(fmt.Errorf("bucket not exist"))
	}
	if s3StorageS3Storage.uploader != nil && s3StorageS3Storage.clients != nil {
		return
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(s3StorageS3Storage.storage.AccessKey, s3StorageS3Storage.storage.SecretKey, ""),
		Region:      aws.String(s3StorageS3Storage.storage.Region),
		Endpoint:    aws.String(s3StorageS3Storage.storage.Endpoint),
	})

	if err != nil {
		panic(fmt.Errorf("unable to create AWS session, %w", err))
	}

	s3StorageS3Storage.uploader = s3manager.NewUploader(sess)
	s3StorageS3Storage.clients = s3.New(sess)
}

func (s3StorageS3Storage *S3Storage) UploadObject(bucketType enums.BucketType, key string, data []byte) {
	s3StorageS3Storage.setS3Client(bucketType)
	bucket := s3StorageS3Storage.buckets[bucketType]

	_, err := s3StorageS3Storage.clients.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && (aerr.Code() == s3.ErrCodeNoSuchBucket || aerr.Code() == "NotFound") {
			_, err = s3StorageS3Storage.clients.CreateBucket(&s3.CreateBucketInput{
				Bucket: aws.String(bucket),
			})
			if err != nil {
				panic(fmt.Errorf("unable to create bucket %q, %w", bucket, err))
			}

			err = s3StorageS3Storage.clients.WaitUntilBucketExists(&s3.HeadBucketInput{
				Bucket: aws.String(bucket),
			})
			if err != nil {
				panic(fmt.Errorf("unable to confirm bucket %q exists, %w", bucket, err))
			}
		} else {
			panic(fmt.Errorf("unable to check bucket %q, %w", bucket, err))
		}
	}

	_, err = s3StorageS3Storage.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		panic(fmt.Errorf("unable to upload data to %q, %w", bucket, err))
	}
}

func (s3StorageS3Storage *S3Storage) DeleteObject(bucketType enums.BucketType, key string) error {
	s3StorageS3Storage.setS3Client(bucketType)
	bucket := s3StorageS3Storage.buckets[bucketType]

	_, err := s3StorageS3Storage.clients.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to delete %q from %q, %w", key, bucket, err)
	}

	err = s3StorageS3Storage.clients.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("error confirming deletion of %q, %w", key, err)
	}
	return nil
}

func (s3StorageS3Storage *S3Storage) GetPresignedURL(bucketType enums.BucketType, objectKey string, expiration time.Duration) string {
	s3StorageS3Storage.setS3Client(bucketType)
	bucket := s3StorageS3Storage.buckets[bucketType]

	req, _ := s3StorageS3Storage.clients.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		panic(fmt.Errorf("failed to generate presigned URL: %w", err))
	}

	return url
}
