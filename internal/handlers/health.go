package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthHandlers struct{}

func NewHealthHandlers() *HealthHandlers {
	return &HealthHandlers{}
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Health handles GET /health
func (h *HealthHandlers) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Version: "0.1.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Ready handles GET /ready
func (h *HealthHandlers) Ready(w http.ResponseWriter, r *http.Request) {
	// TODO: Add actual readiness checks (DB, K8s)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}
