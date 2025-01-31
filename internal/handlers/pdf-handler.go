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
	"sync"
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
		var mu sync.Mutex
		response := make(map[string]string)

		var wg sync.WaitGroup
		sem := make(chan struct{}, 10) // Ограничение на 10 горутин

		for _, adult := range requestData[0].User.Adults {
			wg.Add(1)
			go func(adult models.Adult) {
				defer wg.Done()
				sem <- struct{}{} // Семафор
				defer func() { <-sem }()

				file := models.NewFile(requestData[0].Ticket.ID, adult.FirstName, adult.LastName, cfg)
				file.Bytes, err = pdf.GeneratePDF(requestData[0].Ticket, adult, file.S3URL)
				if err != nil {
					log.Printf("Error generating PDF for %s %s: %v", adult.FirstName, adult.LastName, err)
					return
				}

				mu.Lock()
				resultFiles[file.Filename] = file
				mu.Unlock()

				if cfg.Api.LocalSave {
					err = local.SaveLocalPDF(cfg, file.Filename, file.Bytes)
					if err != nil {
						log.Printf("Failed to save PDF locally for %s %s: %v", adult.FirstName, adult.LastName, err)
						return
					}
					pdfKey := fmt.Sprint(adult.FirstName + "-" + adult.LastName + "-local-pdf")
					mu.Lock()
					response[pdfKey] = file.Filename
					mu.Unlock()
				}

				err = s3_storage.UploadFile(cfg, s3Client, file.Filename, file.Bytes)
				if err != nil {
					log.Printf("Failed to upload to S3 for %s %s: %v", adult.FirstName, adult.LastName, err)
					return
				}
				s3Key := fmt.Sprint(adult.FirstName + "-" + adult.LastName + "-s3-storage-url")
				mu.Lock()
				response[s3Key] = file.S3URL
				mu.Unlock()
			}(adult)
		}

		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode JSON response: %v", err)
			http.Error(w, "Failed to encode JSON response", http.StatusBadRequest)
			return
		}

		log.Printf("Successfully generated and uploaded PDFs")
	}
}
