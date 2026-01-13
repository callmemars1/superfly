package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/superfly/superfly/internal/config"
	"github.com/superfly/superfly/internal/handlers"
	"github.com/superfly/superfly/internal/k8s"
	"github.com/superfly/superfly/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	logger := log.New(os.Stdout, "[superfly] ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting Superfly API server...")

	// Connect to database
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Test database connection
	if err := dbpool.Ping(context.Background()); err != nil {
		logger.Fatalf("Failed to ping database: %v", err)
	}
	logger.Println("✓ Connected to database")

	// Initialize Kubernetes client
	k8sClient, err := k8s.NewClient(cfg.KubernetesInCluster, cfg.Kubeconfig)
	if err != nil {
		logger.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	logger.Println("✓ Connected to Kubernetes cluster")

	// Ensure namespace exists
	ctx := context.Background()
	if err := k8sClient.EnsureNamespace(ctx); err != nil {
		logger.Printf("Warning: Failed to ensure namespace: %v", err)
	}

	// Initialize services
	appService := service.NewAppService(dbpool, k8sClient)

	// Initialize handlers
	appHandlers := handlers.NewAppHandlers(appService)
	healthHandlers := handlers.NewHealthHandlers()

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check routes
	r.Get("/health", healthHandlers.Health)
	r.Get("/ready", healthHandlers.Ready)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/apps", func(r chi.Router) {
			r.Get("/", appHandlers.ListApps)
			r.Post("/", appHandlers.CreateApp)
			r.Get("/{id}", appHandlers.GetApp)
			r.Patch("/{id}", appHandlers.UpdateApp)
			r.Delete("/{id}", appHandlers.DeleteApp)
			r.Post("/{id}/restart", appHandlers.RestartApp)
		})
	})

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.APIHost, cfg.APIPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Printf("🚀 Server listening on %s", addr)
		logger.Printf("📝 API documentation: http://%s/api", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server stopped gracefully")
}
