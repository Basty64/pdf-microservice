package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"net/url"
	"os"
	"pdf-microservice/internal/options"
	"time"
)

func UploadToS3(ctx context.Context, cfg *options.Config, filename string, fileBytes []byte) (string, error) {

	staticCreds := credentials.NewStaticCredentialsProvider(cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, "")

	// Создание резолвера конечной точки с опциями
	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && region == region {
			return aws.Endpoint{
				URL:               cfg.Minio.Endpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		}
		// Если не совпадает сервис или регион, возвращаем ошибку
		return aws.Endpoint{}, fmt.Errorf("endpoint for %s not found", service)
	})

	// Создание конфигурации для клиента S3
	cfgNew, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Minio.Region),
		config.WithCredentialsProvider(staticCreds),
		config.WithEndpointResolverWithOptions(endpointResolver),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Создание клиента S3
	s3Client := s3.NewFromConfig(cfgNew)

	// Чтение файла PDF
	file, err := os.Open(cfg.Minio.FilePath)
	if err != nil {
		log.Fatalf("unable to open file %q, %v", cfg.Minio.FilePath, err)
	}
	defer file.Close()

	// Загрузка файла в S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(cfg.Minio.BucketName),
		Key:    aws.String(cfg.Minio.ObjectKey),
		Body:   file,
	})
	if err != nil {
		log.Fatalf("unable to upload file, %v", err)
	}

	log.Println("File uploaded successfully")

	s3Url, err := generatePresignedUrl(ctx, s3Client, cfg.Minio.BucketName, filename)
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
