package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"log"
	"net/url"
	"os"
	"time"
)

func UploadToS3(ctx context.Context, bucketName, filename string, fileBytes []byte) (string, error) {

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "eu-central-1"
	}

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsAccessKeyID == "" || awsSecretAccessKey == "" {
		log.Println("AWS credentials are not set via environment variables. Using default AWS config.")
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config: %v", err)
	}

	if awsAccessKeyID != "" && awsSecretAccessKey != "" {
		creds := credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")
		cfg.Credentials = aws.NewCredentialsCache(creds)
	}

	s3Client := s3.NewFromConfig(cfg)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	s3Url, err := generatePresignedUrl(ctx, s3Client, bucketName, filename)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url to S3: %v", err)
	}

	return s3Url, nil
}

func generatePresignedUrl(ctx context.Context, s3Client *s3.Client, bucketName, filename string) (string, error) {

	presignClient := s3.NewPresignClient(s3Client)
	presignedGetObject, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	}, func(options *s3.PresignOptions) {
		options.Expires = 60 * 60 * time.Second // URL expires in 1 hour
	})

	if err != nil {
		return "", fmt.Errorf("failed to presign get object to S3: %v", err)
	}
	u, err := url.Parse(presignedGetObject.URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse presigned url: %w", err)
	}
	return u.String(), nil
}
