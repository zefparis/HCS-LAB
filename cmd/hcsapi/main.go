package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/corehuman/hcs-lab-api/internal/hcs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const version = "1.0.0-hcs-lab"

var (
	startTime = time.Now()
	generator *hcs.Generator
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
	Secure  bool   `json:"secure"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// GenerateRequest wraps the input profile to support both flat and nested ("hcs") payloads
type GenerateRequest struct {
	// HCS allows payloads of the form { "hcs": { ...InputProfile... } }
	HCS *hcs.InputProfile `json:"hcs,omitempty"`
	// Embedded InputProfile allows flat payloads { ...InputProfile... }
	hcs.InputProfile
}

func main() {
	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize HCS generator
	var err error
	generator, err = hcs.NewGenerator()
	if err != nil {
		log.Fatalf("Failed to initialize HCS generator: %v", err)
	}

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:*",
			"https://localhost:*",
			"http://127.0.0.1:*",
			"https://127.0.0.1:*",
			"https://*.vercel.app",
			"https://vercel.app",
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Routes
	r.Get("/", handleRoot)
	r.Get("/health", handleHealth)
	r.Post("/api/generate", handleGenerate)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("HCS Lab API v%s starting on %s", version, addr)
	log.Printf("Environment: PORT=%s", port)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"service": "HCS Lab API",
		"version": version,
		"status":  "running",
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)

	response := HealthResponse{
		Status:  "healthy",
		Version: version,
		Uptime:  formatDuration(uptime),
		Secure:  true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	// Parse request body, accepting both flat and nested ("hcs") profiles
	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	// Select the effective input profile
	var input hcs.InputProfile
	if req.HCS != nil {
		input = *req.HCS
	} else {
		input = req.InputProfile
	}

	// Generate HCS codes
	output, err := generator.Generate(&input)
	if err != nil {
		// Determine if it's a validation error or internal error
		if contains(err.Error(), "invalid") || contains(err.Error(), "must be") {
			sendError(w, http.StatusBadRequest, "Validation error", err.Error())
		} else {
			sendError(w, http.StatusInternalServerError, "Generation failed", err.Error())
		}
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func sendError(w http.ResponseWriter, code int, error string, message string) {
	response := ErrorResponse{
		Error:   error,
		Message: message,
		Code:    code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || len(s) > len(substr) && contains(s[1:], substr)
}
