# Superfly Project Structure

Complete overview of the codebase organization.

## Directory Layout

```
superfly/
├── cmd/                          # Application entrypoints
│   └── api/                      # API server
│       └── main.go              # Server initialization & routing
│
├── internal/                     # Private application code
│   ├── config/                  # Configuration management
│   │   └── config.go            # Environment variable loading
│   │
│   ├── db/                      # Generated database code (sqlc)
│   │   ├── db.go                # Database connection helpers
│   │   ├── models.go            # Go structs for DB tables
│   │   └── apps.sql.go          # Generated query functions
│   │
│   ├── handlers/                # HTTP request handlers
│   │   ├── app_handlers.go     # App CRUD endpoints
│   │   └── health.go            # Health check endpoints
│   │
│   ├── k8s/                     # Kubernetes integration
│   │   ├── client.go            # K8s client & operations
│   │   └── resources.go         # K8s resource templates
│   │
│   └── service/                 # Business logic
│       └── app_service.go       # App deployment logic
│
├── db/                          # Database files
│   ├── migrations/              # SQL migration files (goose)
│   │   └── 00001_create_apps_table.sql
│   │
│   └── queries/                 # SQL queries (sqlc)
│       └── apps.sql             # App-related queries
│
├── web/                         # Frontend (future)
│   └── (Svelte app will go here)
│
├── manifests/                   # K8s manifests (future)
│   └── (System K8s configs)
│
├── .air.toml                    # Live reload config
├── .env                         # Environment variables (git-ignored)
├── .env.example                 # Example environment file
├── .gitignore                   # Git ignore patterns
├── go.mod                       # Go module definition
├── go.sum                       # Go dependency checksums
├── Makefile                     # Development commands
├── sqlc.yaml                    # sqlc configuration
│
├── README.md                    # Main documentation
├── QUICKSTART.md                # 10-minute setup guide
├── DEVELOPMENT.md               # Development guide
├── API.md                       # API documentation
├── EXAMPLES.md                  # Usage examples
├── PROJECT_STRUCTURE.md         # This file
│
├── dev-setup.sh                 # Development environment setup
├── verify-setup.sh              # Verify installation
└── test-api.sh                  # API testing script
```

## Key Files Explained

### `cmd/api/main.go`
**Purpose**: Application entry point  
**Responsibilities**:
- Load configuration
- Connect to database
- Initialize Kubernetes client
- Setup HTTP router
- Start server
- Handle graceful shutdown

**Key components**:
```go
- config.Load()              // Load env vars
- pgxpool.New()              // Connect to PostgreSQL
- k8s.NewClient()            // Connect to K8s cluster
- service.NewAppService()    // Initialize business logic
- handlers.NewAppHandlers()  // Initialize HTTP handlers
- chi.NewRouter()            // Setup HTTP router
- server.ListenAndServe()    // Start HTTP server
```

---

### `internal/config/config.go`
**Purpose**: Configuration management  
**Loads**:
- Database connection string
- Kubernetes config path
- API server settings
- Registry URL

---

### `internal/k8s/client.go`
**Purpose**: Kubernetes API client  
**Operations**:
- Create/update/delete Deployments
- Create/update/delete Services
- Create/update/delete Ingresses
- Check deployment status
- Restart deployments

---

### `internal/k8s/resources.go`
**Purpose**: Generate Kubernetes manifests  
**Functions**:
- `BuildDeployment()` - Creates Deployment YAML
- `BuildService()` - Creates Service YAML
- `BuildIngress()` - Creates Ingress YAML

**Features**:
- Health checks (liveness + readiness)
- Resource limits
- Rolling update strategy
- Prometheus annotations
- TLS configuration

---

### `internal/service/app_service.go`
**Purpose**: Core business logic  
**Methods**:
- `CreateApp()` - Create and deploy app
- `GetApp()` - Retrieve app details
- `ListApps()` - List all apps
- `UpdateApp()` - Update and redeploy app
- `DeleteApp()` - Delete app and K8s resources
- `RestartApp()` - Rolling restart

**Key logic**:
- Slug generation and validation
- Domain uniqueness checks
- Asynchronous deployment
- Status tracking

---

### `internal/handlers/app_handlers.go`
**Purpose**: HTTP API handlers  
**Endpoints**:
- `POST /api/apps` - Create app
- `GET /api/apps` - List apps
- `GET /api/apps/:id` - Get app
- `PATCH /api/apps/:id` - Update app
- `DELETE /api/apps/:id` - Delete app
- `POST /api/apps/:id/restart` - Restart app

---

### `db/migrations/00001_create_apps_table.sql`
**Purpose**: Database schema  
**Table**: `apps`  
**Columns**:
- `id` (UUID) - Primary key
- `slug` (VARCHAR) - K8s resource name
- `name` (VARCHAR) - Display name
- `image` (TEXT) - Docker image
- `port`, `replicas`, `cpu_limit`, `memory_limit`
- `domain` - Public domain
- `status` - Deployment status
- `created_at`, `updated_at`, `last_deployed_at`

---

### `db/queries/apps.sql`
**Purpose**: SQL queries for sqlc  
**Queries**:
- `CreateApp` - Insert new app
- `GetApp` - Get by ID
- `GetAppBySlug` - Get by slug
- `ListApps` - Get all apps
- `UpdateApp` - Update app
- `UpdateAppStatus` - Update status
- `DeleteApp` - Delete app
- `CheckSlugExists` - Validate uniqueness
- `CheckDomainExists` - Validate uniqueness

---

### `sqlc.yaml`
**Purpose**: Configure sqlc code generation  
**Settings**:
- SQL engine: PostgreSQL
- Queries directory: `db/queries/`
- Schema directory: `db/migrations/`
- Output directory: `internal/db/`
- Package name: `db`
- SQL package: `pgx/v5`

---

## Data Flow

### Creating an App

```
1. HTTP Request
   POST /api/apps {"name": "My App", "image": "nginx:alpine", ...}
   ↓
2. app_handlers.go
   Validate input, parse JSON
   ↓
3. app_service.go
   - Validate slug uniqueness
   - Validate domain uniqueness
   - Insert into database
   - Start async deployment
   ↓
4. k8s/resources.go
   Generate K8s manifests (Deployment, Service, Ingress)
   ↓
5. k8s/client.go
   Apply manifests to K8s cluster
   ↓
6. Kubernetes
   - Create pods
   - Start containers
   - Setup networking
   - Request TLS cert from Let's Encrypt
   ↓
7. app_service.go
   Update status to "running"
   ↓
8. HTTP Response
   Return app object with ID and status
```

### Querying an App

```
1. HTTP Request
   GET /api/apps/{id}
   ↓
2. app_handlers.go
   Parse ID from URL
   ↓
3. app_service.go
   Query database via sqlc
   ↓
4. PostgreSQL
   Return app row
   ↓
5. HTTP Response
   Return app JSON
```

### Deleting an App

```
1. HTTP Request
   DELETE /api/apps/{id}
   ↓
2. app_handlers.go
   Parse ID
   ↓
3. app_service.go
   - Delete K8s Ingress
   - Delete K8s Service
   - Delete K8s Deployment
   - Delete from database
   ↓
4. Kubernetes
   Clean up all resources
   ↓
5. HTTP Response
   204 No Content
```

---

## Technology Stack

### Backend
- **Language**: Go 1.22
- **HTTP Router**: chi v5
- **Database**: PostgreSQL 16
- **DB Driver**: pgx v5
- **Query Builder**: sqlc
- **Migrations**: goose
- **K8s Client**: client-go
- **Config**: godotenv

### Infrastructure
- **Orchestration**: K3S
- **Ingress**: Traefik
- **Certificates**: cert-manager
- **Registry**: Docker Registry v2
- **Storage**: Local Path Provisioner

### Development Tools
- **Live Reload**: Air
- **Testing**: curl + kubectl
- **Build**: Make

---

## Design Decisions

### Why sqlc?
- Type-safe queries
- No ORM overhead
- SQL-first approach
- Great for PostgreSQL

### Why pgx over database/sql?
- Better PostgreSQL support
- Native support for PostgreSQL types
- Better performance
- Connection pooling

### Why chi router?
- Lightweight
- Standard library compatibility
- Great middleware support
- Simple API

### Why K3S over K8S?
- Lightweight (single binary)
- Perfect for homelab
- Full K8s compatibility
- Easy to install

### Why client-go?
- Official K8s Go client
- Well-maintained
- Complete API coverage
- Type-safe

---

## Database Schema

```sql
CREATE TABLE apps (
    id UUID PRIMARY KEY,
    slug VARCHAR(63) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    image TEXT NOT NULL,
    port INTEGER NOT NULL DEFAULT 8080,
    replicas INTEGER NOT NULL DEFAULT 1,
    cpu_limit VARCHAR(10) NOT NULL DEFAULT '500m',
    memory_limit VARCHAR(10) NOT NULL DEFAULT '256Mi',
    domain VARCHAR(255) UNIQUE,
    health_check_path VARCHAR(255) NOT NULL DEFAULT '/',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_deployed_at TIMESTAMPTZ
);

CREATE INDEX idx_apps_slug ON apps(slug);
CREATE INDEX idx_apps_status ON apps(status);
CREATE INDEX idx_apps_domain ON apps(domain) WHERE domain IS NOT NULL;
```

---

## API Routes

```go
GET    /health                 // Health check
GET    /ready                  // Readiness check

GET    /api/apps               // List all apps
POST   /api/apps               // Create app
GET    /api/apps/:id           // Get app details
PATCH  /api/apps/:id           // Update app
DELETE /api/apps/:id           // Delete app
POST   /api/apps/:id/restart   // Restart app
```

---

## Environment Variables

```bash
DATABASE_URL              # PostgreSQL connection string
KUBERNETES_IN_CLUSTER     # true if running in K8s, false for local dev
KUBECONFIG               # Path to kubeconfig (local dev only)
REGISTRY_URL             # Container registry URL
API_PORT                 # API server port (default: 8080)
API_HOST                 # API server host (default: 0.0.0.0)
ENV                      # Environment (development/production)
LOG_LEVEL                # Log level (debug/info/warn/error)
```

---

## Kubernetes Resources

### Namespace: `superfly-system`
System components:
- PostgreSQL database
- Container registry
- Control plane API (future)
- Web UI (future)

### Namespace: `superfly-apps`
User applications:
- All deployed apps
- Each app gets: Deployment + Service + Ingress

---

## Future Additions

### Feature 2: Build from GitHub
```
internal/
└── build/
    ├── kaniko.go        # Kaniko job creation
    └── git.go           # Git operations
```

### Feature 3: Environment Variables
```
db/
├── migrations/
│   └── 00002_create_env_vars_table.sql
└── queries/
    └── env_vars.sql

internal/
└── service/
    └── env_service.go
```

### Feature 4: Logs (Loki)
```
internal/
└── observability/
    ├── loki.go          # Loki client
    └── logs.go          # Log querying
```

### Feature 5: Metrics (Prometheus)
```
internal/
└── observability/
    ├── prometheus.go    # Prometheus client
    └── metrics.go       # Metric querying
```

### Feature 8: Web UI
```
web/
├── src/
│   ├── routes/
│   ├── lib/
│   └── components/
├── package.json
└── Dockerfile
```

---

## Testing Strategy

### Unit Tests (Future)
```
internal/service/app_service_test.go
internal/k8s/resources_test.go
```

### Integration Tests
```
test-api.sh              # Current API tests
```

### E2E Tests (Future)
```
e2e/
└── deploy_test.go       # Full deployment flow
```

---

## Contributing

When adding new features:

1. **Database changes**:
   - Add migration in `db/migrations/`
   - Add queries in `db/queries/`
   - Run `make sqlc-generate`

2. **Business logic**:
   - Add service in `internal/service/`
   - Use dependency injection

3. **API endpoints**:
   - Add handlers in `internal/handlers/`
   - Register routes in `cmd/api/main.go`

4. **K8s resources**:
   - Add templates in `internal/k8s/resources.go`
   - Add operations in `internal/k8s/client.go`

5. **Documentation**:
   - Update API.md
   - Add examples in EXAMPLES.md
   - Update DEVELOPMENT.md

---

## Quick Reference

### Start Development
```bash
make dev              # Run API with live reload
```

### Database
```bash
make migrate          # Run migrations
make sqlc-generate    # Generate Go code
make db-shell         # Open PostgreSQL
```

### Testing
```bash
./test-api.sh         # Run API tests
```

### Kubernetes
```bash
kubectl get all -n superfly-apps        # View deployed apps
kubectl logs -n superfly-apps -l app=X  # View app logs
kubectl describe pod -n superfly-apps X # Debug pod
```

---

This structure keeps code organized, testable, and maintainable as we add more features! 🚀
