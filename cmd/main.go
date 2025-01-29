package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"pdf-microservice/internal/models"
	"pdf-microservice/internal/options"
	"pdf-microservice/internal/pdf"
	"pdf-microservice/internal/save/local"
	"pdf-microservice/internal/save/s3"
	"time"
)

var Cfg *options.Config

func main() {

	configPath := "pdf-microservice-config-dev.toml"

	cfg, err := options.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	Cfg = cfg

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

	var requestData []models.RequestDataNew
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(r.Body)

	pdfBytes, err := pdf.GeneratePDF(requestData[0])
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s.pdf", uuid.New().String())

	var response map[string]string

	if Cfg.Api.LocalSave {
		err = local.SaveLocalPDF(filename, pdfBytes)
		if err != nil {
			log.Printf("Failed to save PDF locally: %v", err)

			http.Error(w, "Failed to save PDF locally", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response["save-pdf"] = "Successfully generated and saved PDF locally: " + filename

		log.Printf("Successfully generated and saved PDF locally: %s", filename)
		return
	}

	ctx := context.Background()

	s3Url, err := s3.UploadToS3(ctx, Cfg, filename, pdfBytes)
	if err != nil {
		log.Printf("Failed to upload to S3: %v", err)
		http.Error(w, "Failed to upload to S3", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response["s3_url"] = s3Url

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode json response: %v", err)
		http.Error(w, "Failed to encode json response", http.StatusInternalServerError)
	}

	log.Printf("Successfully generated and uploaded PDF to %s: %s", Cfg.Minio.BucketName, s3Url)
}
