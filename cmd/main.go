package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pdf-microservice/internal/models"
	"pdf-microservice/internal/pdf"
	"time"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/generate", generatePDFHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port " + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func generatePDFHandler(w http.ResponseWriter, r *http.Request) {

	var requestData models.RequestData
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

	pdfBytes, err := pdf.GeneratePDF(requestData)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s.pdf", uuid.New().String())

	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		log.Println("S3_BUCKET_NAME environment variable is not set. Using default value.")
		bucketName = "default-pdf-bucket"
	}

	err = saveLocalPDF(filename, pdfBytes)
	if err != nil {
		log.Printf("Failed to save PDF locally: %v", err)
		http.Error(w, "Failed to save PDF locally", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Successfully generated and saved PDF locally: " + filename}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode json response: %v", err)
		http.Error(w, "Failed to encode json response", http.StatusInternalServerError)
	}
	log.Printf("Successfully generated and saved PDF locally: %s", filename)
	return

	//s3Url, err := s3.UploadToS3(ctx, bucketName, filename, pdfBytes)
	//if err != nil {
	//	log.Printf("Failed to upload to S3: %v", err)
	//	http.Error(w, "Failed to upload to S3", http.StatusInternalServerError)
	//	return
	//}

	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	////response := map[string]string{"s3_url": s3Url}
	//
	//response := "success"
	//
	//if err := json.NewEncoder(w).Encode(response); err != nil {
	//	log.Printf("Failed to encode json response: %v", err)
	//	http.Error(w, "Failed to encode json response", http.StatusInternalServerError)
	//}
	//
	//log.Printf("Successfully generated and uploaded PDF to %s: %s", bucketName, s3Url)
}

func saveLocalPDF(filename string, pdfBytes []byte) error {
	// Create a "local-pdfs" directory if it doesn't exist.
	if _, err := os.Stat("local-pdfs"); os.IsNotExist(err) {
		err := os.Mkdir("local-pdfs", os.ModeDir|0755)
		if err != nil {
			return fmt.Errorf("failed to create directory 'local-pdfs': %w", err)
		}
	}

	filePath := filepath.Join("local-pdfs", filename)
	err := ioutil.WriteFile(filePath, pdfBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to save pdf to file: %w", err)
	}
	return nil
}
