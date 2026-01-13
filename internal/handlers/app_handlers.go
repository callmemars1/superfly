package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/superfly/superfly/internal/service"
)

type AppHandlers struct {
	appService *service.AppService
}

func NewAppHandlers(appService *service.AppService) *AppHandlers {
	return &AppHandlers{
		appService: appService,
	}
}

// CreateAppRequest represents the request body for creating an app
type CreateAppRequest struct {
	Name            string `json:"name"`
	Slug            string `json:"slug,omitempty"`
	Image           string `json:"image"`
	Port            int32  `json:"port,omitempty"`
	Replicas        int32  `json:"replicas,omitempty"`
	CPULimit        string `json:"cpu_limit,omitempty"`
	MemoryLimit     string `json:"memory_limit,omitempty"`
	Domain          string `json:"domain,omitempty"`
	HealthCheckPath string `json:"health_check_path,omitempty"`
}

// UpdateAppRequest represents the request body for updating an app
type UpdateAppRequest struct {
	Name            *string `json:"name,omitempty"`
	Image           *string `json:"image,omitempty"`
	Port            *int32  `json:"port,omitempty"`
	Replicas        *int32  `json:"replicas,omitempty"`
	CPULimit        *string `json:"cpu_limit,omitempty"`
	MemoryLimit     *string `json:"memory_limit,omitempty"`
	Domain          *string `json:"domain,omitempty"`
	HealthCheckPath *string `json:"health_check_path,omitempty"`
}

// CreateApp handles POST /api/apps
func (h *AppHandlers) CreateApp(w http.ResponseWriter, r *http.Request) {
	var req CreateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Image == "" {
		respondError(w, http.StatusBadRequest, "image is required")
		return
	}

	// Create app
	app, err := h.appService.CreateApp(r.Context(), service.CreateAppInput{
		Name:            req.Name,
		Slug:            req.Slug,
		Image:           req.Image,
		Port:            req.Port,
		Replicas:        req.Replicas,
		CPULimit:        req.CPULimit,
		MemoryLimit:     req.MemoryLimit,
		Domain:          req.Domain,
		HealthCheckPath: req.HealthCheckPath,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, app)
}

// ListApps handles GET /api/apps
func (h *AppHandlers) ListApps(w http.ResponseWriter, r *http.Request) {
	apps, err := h.appService.ListApps(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, apps)
}

// GetApp handles GET /api/apps/:id
func (h *AppHandlers) GetApp(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid app ID")
		return
	}

	app, err := h.appService.GetApp(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "App not found")
		return
	}

	respondJSON(w, http.StatusOK, app)
}

// UpdateApp handles PATCH /api/apps/:id
func (h *AppHandlers) UpdateApp(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid app ID")
		return
	}

	var req UpdateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	app, err := h.appService.UpdateApp(r.Context(), id, service.UpdateAppInput{
		Name:            req.Name,
		Image:           req.Image,
		Port:            req.Port,
		Replicas:        req.Replicas,
		CPULimit:        req.CPULimit,
		MemoryLimit:     req.MemoryLimit,
		Domain:          req.Domain,
		HealthCheckPath: req.HealthCheckPath,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, app)
}

// DeleteApp handles DELETE /api/apps/:id
func (h *AppHandlers) DeleteApp(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid app ID")
		return
	}

	if err := h.appService.DeleteApp(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RestartApp handles POST /api/apps/:id/restart
func (h *AppHandlers) RestartApp(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid app ID")
		return
	}

	if err := h.appService.RestartApp(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "App restart initiated",
	})
}

// Helper functions

type errorResponse struct {
	Error string `json:"error"`
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
