package local

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"pdf-microservice/internal/options"
)

func SaveLocalPDF(cfg *options.Config, filename string, pdfBytes []byte) error {

	if _, err := os.Stat(cfg.S3.FilePath); os.IsNotExist(err) {
		err = os.Mkdir(cfg.S3.FilePath, os.ModeDir|0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", cfg.Api.DirName, err)
		}
	}

	filePath := filepath.Join(cfg.Api.DirName, filename)
	err := ioutil.WriteFile(filePath, pdfBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to save pdf to file: %w", err)
	}
	return nil
}
