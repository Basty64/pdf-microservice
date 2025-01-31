package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"pdf-microservice/internal/handlers"
	"pdf-microservice/internal/options"
	"pdf-microservice/internal/save/s3-storage"
	"time"
)

func main() {

	configPath := "pdf-microservice-config-dev.toml"

	cfg, err := options.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	s3Client, err := s3_storage.NewS3Client(cfg)
	if err != nil {
		log.Fatalf("Error creating S3 client: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Method(http.MethodPost, "/generate", handlers.GeneratePDFHandler(cfg, s3Client))

	log.Println("Server starting on port " + cfg.Api.Port)
	if err := http.ListenAndServe(":"+cfg.Api.Port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
