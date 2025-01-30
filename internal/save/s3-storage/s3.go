package s3_storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"pdf-microservice/internal/options"
)

func NewS3Client(cfg *options.Config) (*minio.Client, error) {

	client, err := minio.New(cfg.S3.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.S3.AccessKeyID, cfg.S3.SecretAccessKey, ""),
	})

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return client, nil
}

func UploadFile(cfg *options.Config, client *minio.Client, filename string, fileBytes []byte) error {

	key := fmt.Sprintf("tickets/" + filename)

	// Загрузка файла в S3
	_, err := client.PutObject(context.Background(), cfg.S3.BucketName, key, bytes.NewReader(fileBytes), int64(len(fileBytes)), minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		return err
	}

	return nil

}
