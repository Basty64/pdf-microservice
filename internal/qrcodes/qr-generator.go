package qrcodes

import (
	"bytes"
	"fmt"
	"github.com/skip2/go-qrcode"
	"image"
	"image/png"
	"pdf-microservice/internal/options"
	"strings"
)

func CreateURL(cfg *options.Config, filename string) string {

	url := []string{"https:", cfg.S3.Endpoint, cfg.S3.BucketName, "tickets", filename}

	s3Url := strings.Join(url, "/")

	return s3Url
}

func GenerateQRCode(data string) ([]byte, error) {
	qrCode, err := qrcode.Encode(data, qrcode.Medium, 512)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	img, _, err := image.Decode(bytes.NewReader(qrCode))
	if err != nil {
		return nil, fmt.Errorf("failed to decode qr code: %w", err)
	}

	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode qr code to png: %w", err)
	}

	return buf.Bytes(), nil
}
