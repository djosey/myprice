// Package main runs the HTTP API server for receipt analysis.
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"myprice/server"
)

func main() {
	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get upload directory
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		// Default to uploads folder in current directory
		cwd, _ := os.Getwd()
		uploadDir = filepath.Join(cwd, "uploads")
	}

	// Create server
	srv := server.NewServer(uploadDir)

	// Create mux and register routes
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	// Add CORS middleware
	handler := corsMiddleware(mux)

	log.Printf("Starting MyPrice API server on :%s", port)
	log.Printf("Upload directory: %s", uploadDir)
	log.Printf("Endpoints:")
	log.Printf("  GET  /api/health       - Health check")
	log.Printf("  POST /api/upload       - Upload image")
	log.Printf("  POST /api/load-textract - Load Textract JSON")
	log.Printf("  POST /api/analyze      - Run full analysis")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// corsMiddleware adds CORS headers to all responses.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

