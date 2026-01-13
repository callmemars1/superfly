package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/superfly/superfly/internal/db"
	"github.com/superfly/superfly/internal/k8s"
)

type AppService struct {
	pool      *pgxpool.Pool
	queries   *db.Queries
	k8sClient *k8s.Client
}

func NewAppService(pool *pgxpool.Pool, k8sClient *k8s.Client) *AppService {
	return &AppService{
		pool:      pool,
		queries:   db.New(pool),
		k8sClient: k8sClient,
	}
}

type CreateAppInput struct {
	Name            string
	Slug            string
	Image           string
	Port            int32
	Replicas        int32
	CPULimit        string
	MemoryLimit     string
	Domain          string
	HealthCheckPath string
}

type UpdateAppInput struct {
	Name            *string
	Image           *string
	Port            *int32
	Replicas        *int32
	CPULimit        *string
	MemoryLimit     *string
	Domain          *string
	HealthCheckPath *string
}

// CreateApp creates a new app and deploys it to Kubernetes
func (s *AppService) CreateApp(ctx context.Context, input CreateAppInput) (*db.App, error) {
	// Generate slug if not provided
	if input.Slug == "" {
		input.Slug = slugify(input.Name)
	}

	// Validate slug
	if err := validateSlug(input.Slug); err != nil {
		return nil, err
	}

	// Check if slug already exists
	exists, err := s.queries.CheckSlugExists(ctx, input.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("app with slug '%s' already exists", input.Slug)
	}

	// Check if domain already exists
	if input.Domain != "" {
		exists, err := s.queries.CheckDomainExists(ctx, input.Domain)
		if err != nil {
			return nil, fmt.Errorf("failed to check domain: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("domain '%s' already in use", input.Domain)
		}
	}

	// Set defaults
	if input.Port == 0 {
		input.Port = 8080
	}
	if input.Replicas == 0 {
		input.Replicas = 1
	}
	if input.CPULimit == "" {
		input.CPULimit = "500m"
	}
	if input.MemoryLimit == "" {
		input.MemoryLimit = "256Mi"
	}
	if input.HealthCheckPath == "" {
		input.HealthCheckPath = "/"
	}

	// Create app in database
	app, err := s.queries.CreateApp(ctx, db.CreateAppParams{
		Slug:            input.Slug,
		Name:            input.Name,
		Image:           input.Image,
		Port:            input.Port,
		Replicas:        input.Replicas,
		CpuLimit:        input.CPULimit,
		MemoryLimit:     input.MemoryLimit,
		Domain:          input.Domain,
		HealthCheckPath: input.HealthCheckPath,
		Status:          "pending",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create app: %w", err)
	}

	// Deploy to Kubernetes
	go func() {
		deployCtx := context.Background()
		if err := s.deployApp(deployCtx, &app); err != nil {
			// Update status to failed
			_, _ = s.queries.UpdateAppStatus(deployCtx, db.UpdateAppStatusParams{
				ID:     app.ID,
				Status: "failed",
			})
		}
	}()

	return &app, nil
}

// deployApp deploys an app to Kubernetes
func (s *AppService) deployApp(ctx context.Context, app *db.App) error {
	// Update status to deploying
	_, err := s.queries.UpdateAppStatus(ctx, db.UpdateAppStatusParams{
		ID:     app.ID,
		Status: "deploying",
	})
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Ensure namespace exists
	if err := s.k8sClient.EnsureNamespace(ctx); err != nil {
		return fmt.Errorf("failed to ensure namespace: %w", err)
	}

	spec := k8s.AppSpec{
		Name:            app.Name,
		Slug:            app.Slug,
		Image:           app.Image,
		Port:            app.Port,
		Replicas:        app.Replicas,
		CPULimit:        app.CpuLimit,
		MemoryLimit:     app.MemoryLimit,
		Domain:          app.Domain,
		HealthCheckPath: app.HealthCheckPath,
	}

	// Create Deployment
	deployment := k8s.BuildDeployment(spec)
	if err := s.k8sClient.ApplyDeployment(ctx, deployment); err != nil {
		return fmt.Errorf("failed to apply deployment: %w", err)
	}

	// Create Service
	service := k8s.BuildService(spec)
	if err := s.k8sClient.ApplyService(ctx, service); err != nil {
		return fmt.Errorf("failed to apply service: %w", err)
	}

	// Create Ingress if domain is provided
	if app.Domain != "" {
		ingress := k8s.BuildIngress(spec)
		if err := s.k8sClient.ApplyIngress(ctx, ingress); err != nil {
			return fmt.Errorf("failed to apply ingress: %w", err)
		}
	}

	// Wait for deployment to be ready (with timeout)
	waitCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	if err := s.k8sClient.WaitForDeployment(waitCtx, app.Slug, 5*time.Minute); err != nil {
		// Don't fail completely, just log warning
		// Deployment might still succeed after we return
	}

	// Update status to running
	_, err = s.queries.UpdateAppStatus(ctx, db.UpdateAppStatusParams{
		ID:     app.ID,
		Status: "running",
	})
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// GetApp gets an app by ID
func (s *AppService) GetApp(ctx context.Context, id uuid.UUID) (*db.App, error) {
	app, err := s.queries.GetApp(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get app: %w", err)
	}
	return &app, nil
}

// ListApps lists all apps
func (s *AppService) ListApps(ctx context.Context) ([]db.App, error) {
	apps, err := s.queries.ListApps(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}
	return apps, nil
}

// UpdateApp updates an app and redeploys if necessary
func (s *AppService) UpdateApp(ctx context.Context, id uuid.UUID, input UpdateAppInput) (*db.App, error) {
	// Get current app
	currentApp, err := s.queries.GetApp(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get app: %w", err)
	}

	// Check if domain changed and if new domain is available
	if input.Domain != nil && *input.Domain != currentApp.Domain {
		exists, err := s.queries.CheckDomainExists(ctx, *input.Domain)
		if err != nil {
			return nil, fmt.Errorf("failed to check domain: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("domain '%s' already in use", *input.Domain)
		}
	}

	// Update app in database
	app, err := s.queries.UpdateApp(ctx, db.UpdateAppParams{
		ID:              id,
		Name:            input.Name,
		Image:           input.Image,
		Port:            input.Port,
		Replicas:        input.Replicas,
		CpuLimit:        input.CPULimit,
		MemoryLimit:     input.MemoryLimit,
		Domain:          input.Domain,
		HealthCheckPath: input.HealthCheckPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update app: %w", err)
	}

	// Redeploy if certain fields changed
	needsRedeploy := input.Image != nil || input.Port != nil || input.Replicas != nil ||
		input.CPULimit != nil || input.MemoryLimit != nil || input.HealthCheckPath != nil ||
		input.Domain != nil

	if needsRedeploy {
		go func() {
			deployCtx := context.Background()
			_ = s.deployApp(deployCtx, &app)
		}()
	}

	return &app, nil
}

// DeleteApp deletes an app and its Kubernetes resources
func (s *AppService) DeleteApp(ctx context.Context, id uuid.UUID) error {
	// Get app
	app, err := s.queries.GetApp(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err)
	}

	// Delete Kubernetes resources
	_ = s.k8sClient.DeleteIngress(ctx, app.Slug)
	_ = s.k8sClient.DeleteService(ctx, app.Slug)
	_ = s.k8sClient.DeleteDeployment(ctx, app.Slug)

	// Delete from database
	if err := s.queries.DeleteApp(ctx, id); err != nil {
		return fmt.Errorf("failed to delete app: %w", err)
	}

	return nil
}

// RestartApp restarts an app by triggering a rolling restart
func (s *AppService) RestartApp(ctx context.Context, id uuid.UUID) error {
	// Get app
	app, err := s.queries.GetApp(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err)
	}

	// Restart deployment
	if err := s.k8sClient.RestartDeployment(ctx, app.Slug); err != nil {
		return fmt.Errorf("failed to restart deployment: %w", err)
	}

	return nil
}

// slugify converts a name to a valid Kubernetes resource name (slug)
func slugify(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove invalid characters (keep only alphanumeric and hyphens)
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// Limit length to 63 characters (K8s limit)
	if len(slug) > 63 {
		slug = slug[:63]
	}

	return slug
}

// validateSlug validates a slug according to Kubernetes naming rules
func validateSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	if len(slug) > 63 {
		return fmt.Errorf("slug cannot be longer than 63 characters")
	}

	// Must start and end with alphanumeric
	reg := regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
	if !reg.MatchString(slug) {
		return fmt.Errorf("slug must consist of lowercase alphanumeric characters or '-', and must start and end with an alphanumeric character")
	}

	return nil
}
