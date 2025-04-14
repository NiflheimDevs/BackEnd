package storageimpl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/domain/enums"
)

type S3Storage struct {
	Storage  *bootstrap.S3
	Client   *s3.Client
	Uploader *manager.Uploader
	Buckets  map[enums.BucketType]string
}

func NewS3Storage(constants *bootstrap.Constants, storage *bootstrap.S3) *S3Storage {
	buckets := make(map[enums.BucketType]string)
	buckets[enums.ProfilePic] = storage.Buckets.ProfilePic
	return &S3Storage{
		Storage: storage,
		Buckets: buckets,
	}
}

func (s *S3Storage) setS3Client(ctx context.Context, bucketType enums.BucketType) {
	bucketTypes := enums.GetAllBucketTypes()
	if !slices.Contains(bucketTypes, bucketType) {
		panic(fmt.Errorf("bucket does not exist"))
	}
	if s.Client != nil && s.Uploader != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(s.Storage.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(s.Storage.AccessKey, s.Storage.SecretKey, ""),
		),
	)
	if err != nil {
		panic(fmt.Errorf("unable to load AWS config: %w", err))
	}

	s.Client = s3.New(s3.Options{
		Credentials:  cfg.Credentials,
		Region:       s.Storage.Region,
		BaseEndpoint: aws.String(s.Storage.Endpoint),
		HTTPClient:   cfg.HTTPClient,
		Retryer:      cfg.Retryer(),
		Logger:       cfg.Logger,
	})

	s.Uploader = manager.NewUploader(s.Client)
}

func (s *S3Storage) UploadObject(ctx context.Context, bucketType enums.BucketType, key string, file *multipart.FileHeader) error {
	// Set up client and get bucket name
	s.setS3Client(ctx, bucketType)
	bucket := s.Buckets[bucketType]

	// Open the file
	fileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("unable to open file %q: %w", file.Filename, err)
	}
	defer fileReader.Close()

	// Check if bucket exists
	_, err = s.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NotFound" {
			// Create bucket with proper region configuration
			_, err = s.Client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket: aws.String(bucket),
				CreateBucketConfiguration: &types.CreateBucketConfiguration{
					LocationConstraint: types.BucketLocationConstraint(s.Storage.Region),
				},
			})
			if err != nil {
				return fmt.Errorf("unable to create bucket %q: %w", bucket, err)
			}

			// Wait for bucket creation with proper error handling
			waiter := s3.NewBucketExistsWaiter(s.Client)
			err = waiter.Wait(ctx, &s3.HeadBucketInput{
				Bucket: aws.String(bucket),
			}, 5*time.Minute) // More reasonable timeout
			if err != nil {
				return fmt.Errorf("timed out waiting for bucket %q creation: %w", bucket, err)
			}
		} else {
			return fmt.Errorf("unable to check bucket %q: %w", bucket, err)
		}
	}

	// Upload the file
	_, err = s.Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   fileReader,
	})
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q: %w", file.Filename, bucket, err)
	}

	return nil
}

func (s *S3Storage) DeleteObject(ctx context.Context, bucketType enums.BucketType, key string) error {
	// Initialize client with error handling
	s.setS3Client(ctx, bucketType)

	bucket := s.Buckets[bucketType]

	// First check if bucket exists
	_, err := s.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("bucket %q does not exist or is inaccessible: %w", bucket, err)
	}

	// Delete the object
	_, err = s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to delete %q from %q: %w", key, bucket, err)
	}

	// Wait for deletion to complete (with more reasonable timeout)
	waiter := s3.NewObjectNotExistsWaiter(s.Client)
	waitTimeout := 30 * time.Second // More generous timeout
	err = waiter.Wait(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, waitTimeout)

	if err != nil {
		log.Printf("warning: timeout waiting for deletion confirmation of %q: %v", key, err)
	}

	return nil
}

func (s *S3Storage) GetPresignedURL(ctx context.Context, bucketType enums.BucketType, objectKey string, expiration time.Duration) string {
	s.setS3Client(ctx, bucketType)
	bucket := s.Buckets[bucketType]

	presignClient := s3.NewPresignClient(s.Client)

	url, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		panic(fmt.Errorf("failed to generate presigned URL: %w", err))
	}

	return url.URL
}
