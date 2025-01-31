package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"net/http"
	"pdf-microservice/internal/models"
	"pdf-microservice/internal/options"
	"pdf-microservice/internal/pdf"
	"pdf-microservice/internal/save/local"
	"pdf-microservice/internal/save/s3-storage"
)

func GeneratePDFHandler(cfg *options.Config, s3Client *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var err error
		var requestData []models.RequestData
		if err = json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		resultFiles := make(map[string]*models.File)
		var pdfKey string
		var s3Key string
		response := make(map[string]string)

		for _, adult := range requestData[0].User.Adults {

			file := models.NewFile(requestData[0].Ticket.ID, adult.FirstName, adult.LastName, cfg)
			file.Bytes, err = pdf.GeneratePDF(requestData[0].Ticket, adult, file.S3URL)
			if err != nil {
				log.Printf("Error generating PDF: %v", err)
				http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
			}
			resultFiles[file.Filename] = file

			if cfg.Api.LocalSave {
				err = local.SaveLocalPDF(cfg, file.Filename, file.Bytes)
				if err != nil {
					log.Printf("Failed to save PDF locally: %v", err)
					http.Error(w, "Failed to save PDF locally", http.StatusBadRequest)
				}
				pdfKey = fmt.Sprint(adult.FirstName + "-local-pdf")
				response[pdfKey] = file.Filename
				log.Printf("Successfully generated and saved PDF locally: %s", file.Filename)
			}

			err = s3_storage.UploadFile(cfg, s3Client, file.Filename, file.Bytes)
			if err != nil {
				log.Printf("Failed to upload to S3: %v", err)
				http.Error(w, "Failed to upload to S3", http.StatusBadRequest)
			}

			s3Key = fmt.Sprint(adult.FirstName + "-" + adult.LastName + "-s3-storage-url")

			response[s3Key] = file.S3URL
		}

		if err == nil {
			log.Printf("Successfully generated and uploaded PDF to %s: %s", cfg.S3.BucketName, response)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode json response: %v", err)
			http.Error(w, "Failed to encode json response", http.StatusBadRequest)

		}
	}
}
