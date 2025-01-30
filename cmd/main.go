package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/minio/minio-go/v7"
	"log"
	"net/http"
	"pdf-microservice/internal/models"
	"pdf-microservice/internal/options"
	"pdf-microservice/internal/pdf"
	"pdf-microservice/internal/save/local"
	"pdf-microservice/internal/save/s3-storage"
	"time"
)

var Cfg *options.Config
var s3Client *minio.Client

func main() {

	configPath := "pdf-microservice-config-dev.toml"

	cfg, err := options.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	Cfg = cfg

	s3Client, err = s3_storage.NewS3Client(Cfg)
	if err != nil {
		log.Fatalf("Error creating S3 client: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/generate", generatePDFHandler)

	log.Println("Server starting on port " + cfg.Api.Port)
	if err := http.ListenAndServe(":"+cfg.Api.Port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func generatePDFHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	var requestData []models.RequestData
	if err = json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	resultPDFSBytes := make(map[string][]byte)

	for _, adult := range requestData[0].User.Adults {

		resultPDFSBytes[adult.FirstName], err = pdf.GeneratePDF(requestData[0].Ticket, adult)
		if err != nil {
			log.Printf("Error generating PDF: %v", err)
			http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		}
	}

	filename := fmt.Sprintf("%d.pdf", requestData[0].Ticket.ID)

	response := make(map[string]string)
	var pdfKey string

	if Cfg.Api.LocalSave {
		for _, adult := range requestData[0].User.Adults {
			err = local.SaveLocalPDF(Cfg, filename, resultPDFSBytes[adult.FirstName])
			if err != nil {
				log.Printf("Failed to save PDF locally: %v", err)
				http.Error(w, "Failed to save PDF locally", http.StatusBadRequest)
			}
			pdfKey = fmt.Sprint(adult.FirstName + "-local-pdf")
			response[pdfKey] = filename
		}

		log.Printf("Successfully generated and saved PDF locally: %s", filename)
	}

	var s3Url string
	var s3Key string

	for _, adult := range requestData[0].User.Adults {
		s3Url, err = s3_storage.UploadFile(Cfg, s3Client, filename, resultPDFSBytes[adult.FirstName])
		if err != nil {
			log.Printf("Failed to upload to S3: %v", err)
			http.Error(w, "Failed to upload to S3", http.StatusBadRequest)
		}

		s3Key = fmt.Sprint(adult.FirstName + "-s3-storage-url")

		response[s3Key] = s3Url

	}

	if err == nil {
		log.Printf("Successfully generated and uploaded PDF to %s: %s", Cfg.S3.BucketName, response)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode json response: %v", err)
		http.Error(w, "Failed to encode json response", http.StatusBadRequest)

	}
}

func SaveFile() {

}
