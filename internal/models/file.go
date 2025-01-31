package models

import (
	"fmt"
	"pdf-microservice/internal/options"
	"strings"
)

type File struct {
	Filename string
	S3URL    string
	Bytes    []byte
}

func NewFile(ID int, FirstName string, LastName string, cfg *options.Config) *File {

	filename := fmt.Sprintf("%d-%s-%s.pdf", ID, FirstName, LastName)

	return &File{
		Filename: filename,
		S3URL:    CreateURL(cfg, filename),
	}

}

func CreateURL(cfg *options.Config, filename string) string {

	url := []string{"https:", cfg.S3.Endpoint, cfg.S3.BucketName, "tickets", filename}

	s3Url := strings.Join(url, "/")

	return s3Url
}
