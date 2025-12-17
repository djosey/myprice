// Package server provides HTTP API endpoints for the receipt analysis tools.
package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"myprice/tools"
)

// Server holds the HTTP server configuration.
type Server struct {
	uploadDir   string
	textractDir string
	projectRoot string
}

// NewServer creates a new HTTP API server.
func NewServer(uploadDir string) *Server {
	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Warning: could not create upload dir: %v", err)
	}

	// Determine project root (parent of uploads)
	projectRoot := filepath.Dir(uploadDir)

	// Textract cache directory
	textractDir := filepath.Join(projectRoot, "textract_cache")
	if err := os.MkdirAll(textractDir, 0755); err != nil {
		log.Printf("Warning: could not create textract cache dir: %v", err)
	}

	return &Server{
		uploadDir:   uploadDir,
		textractDir: textractDir,
		projectRoot: projectRoot,
	}
}

// RegisterRoutes registers all API endpoints.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/upload", s.handleUpload)
	mux.HandleFunc("/api/analyze", s.handleAnalyze)
}

// handleHealth returns server health status.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"service": "myprice-api",
		"version": "0.1.0",
	})
}

// UploadResponse is returned after successful image upload.
type UploadResponse struct {
	Success  bool   `json:"success"`
	FilePath string `json:"file_path"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

// handleUpload handles image file uploads.
func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		jsonError(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		jsonError(w, "No image file provided: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create destination file
	destPath := filepath.Join(s.uploadDir, header.Filename)
	dest, err := os.Create(destPath)
	if err != nil {
		jsonError(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	// Copy file contents
	size, err := io.Copy(dest, file)
	if err != nil {
		jsonError(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Determine MIME type
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	log.Printf("Uploaded image: %s (%d bytes)", destPath, size)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadResponse{
		Success:  true,
		FilePath: destPath,
		FileName: header.Filename,
		Size:     size,
		MimeType: mimeType,
	})
}

// AnalyzeRequest is the request body for the analyze endpoint.
type AnalyzeRequest struct {
	ImagePath string `json:"image_path"`
}

// AnalyzeResponse contains both textract and parsed output.
type AnalyzeResponse struct {
	Textract  tools.LoadTextractOutput `json:"textract"`
	LLMOutput map[string]any           `json:"llm_output"`
	Source    string                   `json:"source"` // Where the textract came from
}

// handleAnalyze runs the full analysis pipeline.
func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Find the actual image path
	imagePath := req.ImagePath
	if !filepath.IsAbs(imagePath) {
		// Check if it's in uploads folder
		uploadPath := filepath.Join(s.uploadDir, filepath.Base(imagePath))
		if _, err := os.Stat(uploadPath); err == nil {
			imagePath = uploadPath
		}
	}

	log.Printf("Analyzing image: %s", imagePath)

	// Find or generate Textract output
	textractPath, source, err := s.findOrRunTextract(imagePath)
	if err != nil {
		jsonError(w, "Textract failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Using Textract file: %s (source: %s)", textractPath, source)

	// Load textract data
	textractInput := tools.LoadTextractInput{Path: textractPath}
	_, textractOutput, err := tools.HandleLoadTextract(r.Context(), nil, textractInput)
	if err != nil {
		jsonError(w, "Failed to load textract: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build structured output from textract
	llmOutput := parseTextractToReceipt(textractOutput)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AnalyzeResponse{
		Textract:  textractOutput,
		LLMOutput: llmOutput,
		Source:    source,
	})
}

// findOrRunTextract finds an existing Textract result or runs Textract on the image.
func (s *Server) findOrRunTextract(imagePath string) (string, string, error) {
	// Get base name of image
	baseName := filepath.Base(imagePath)
	nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	// Check for cached textract output in cache folder
	cachedPath := filepath.Join(s.textractDir, nameWithoutExt+"_textract.json")
	if _, err := os.Stat(cachedPath); err == nil {
		log.Printf("Found cached Textract: %s", cachedPath)
		return cachedPath, "cached", nil
	}

	// Verify image exists before running Textract
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("image file not found: %s", imagePath)
	}

	// Run AWS Textract on the image
	log.Printf("Running AWS Textract on image: %s", imagePath)
	textractOutput, err := s.runTextract(imagePath, cachedPath)
	if err != nil {
		log.Printf("AWS Textract failed: %v", err)
		return "", "", fmt.Errorf("AWS Textract failed: %w. Please ensure AWS CLI is configured", err)
	}

	return textractOutput, "aws_textract", nil
}

// runTextract calls AWS Textract CLI to process an image.
func (s *Server) runTextract(imagePath, outputPath string) (string, error) {
	// Read image and base64 encode it
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	// Base64 encode the image
	base64Data := base64.StdEncoding.EncodeToString(imageData)

	log.Printf("Running AWS Textract (image size: %d bytes, base64 size: %d)", len(imageData), len(base64Data))

	// Call AWS Textract via CLI
	cmd := exec.Command("aws", "textract", "detect-document-text",
		"--region", "us-east-1",
		"--document", fmt.Sprintf(`{"Bytes":"%s"}`, base64Data),
	)

	output, err := cmd.Output()
	if err != nil {
		// Get stderr for better error messages
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("textract failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("textract command failed: %w", err)
	}

	// Save to cache
	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		return "", fmt.Errorf("failed to cache textract output: %w", err)
	}

	log.Printf("Cached Textract output: %s (%d bytes)", outputPath, len(output))
	return outputPath, nil
}

// parseTextractToReceipt converts textract lines to a structured receipt.
func parseTextractToReceipt(textract tools.LoadTextractOutput) map[string]any {
	receipt := map[string]any{
		"vendor":           "",
		"date":             "",
		"items":            []map[string]any{},
		"subtotal":         0.0,
		"tax":              0.0,
		"total":            0.0,
		"confidence_notes": "Parsed from Textract OCR output",
		"anomalies":        []string{},
	}

	items := []map[string]any{}
	var vendor string
	var date string
	var subtotal, tax, total float64

	for i, line := range textract.Lines {
		text := line.Text

		// First high-confidence line is often the vendor
		if i < 3 && line.Confidence > 90 && vendor == "" && len(text) > 3 {
			vendor = text
		}

		// Look for date patterns
		if containsDate(text) && date == "" {
			date = text
		}

		// Look for dollar amounts
		if containsPrice(text) {
			lowerText := strings.ToLower(text)
			price := extractPrice(text)

			if strings.Contains(lowerText, "subtotal") {
				subtotal = price
			} else if strings.Contains(lowerText, "tax") {
				tax = price
			} else if strings.Contains(lowerText, "total") && !strings.Contains(lowerText, "subtotal") {
				total = price
			} else if price > 0 {
				// Line item
				name := extractItemName(text)
				if name != "" && len(name) > 1 {
					items = append(items, map[string]any{
						"name":  name,
						"qty":   1,
						"price": price,
					})
				}
			}
		}
	}

	receipt["vendor"] = vendor
	receipt["date"] = date
	receipt["items"] = items
	receipt["subtotal"] = subtotal
	receipt["tax"] = tax
	receipt["total"] = total

	return receipt
}

// jsonError sends a JSON error response.
func jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"error":   true,
		"message": message,
	})
}
