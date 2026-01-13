# Superfly

> Homelab platform for software developers - Deploy apps with SaaS-like experience

**Current Status**: ✅ Feature 1 Complete - Deploy Pre-built Docker Images

## Quick Start (Development)

**⚡ Want to get started fast?** See [QUICKSTART.md](QUICKSTART.md) for a 10-minute setup guide.

## What Works Now

✅ **Deploy any Docker image** with automatic TLS  
✅ **REST API** for app management  
✅ **Kubernetes orchestration** (K3S)  
✅ **Automatic certificates** (Let's Encrypt)  
✅ **Ingress routing** (Traefik)  
✅ **Health checks** and rolling updates  
✅ **Resource limits** (CPU/Memory)  

### Example

```bash
# Deploy an app in one command
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My App",
    "image": "nginx:alpine",
    "port": 80,
    "domain": "myapp.example.com"
  }'

# App is now live at https://myapp.example.com with TLS! 🎉
```

## 📖 Documentation

Choose your path:

### Getting Started
- **🚀 [TL;DR](TLDR.md)** - Commands only, 10 minutes
- **📋 [Deployment Guide](DEPLOYMENT_GUIDE.md)** - Complete step-by-step for production
- **👀 [Visual Guide](VISUAL_GUIDE.md)** - Diagrams and flowcharts
- **⚡ [Quick Start](QUICKSTART.md)** - Fast setup for local development

### Reference
- **📚 [API Reference](API.md)** - Complete API documentation
- **💡 [Examples](EXAMPLES.md)** - Real-world usage examples
- **🔧 [Development Guide](DEVELOPMENT.md)** - Contributing guide
- **🏗️ [Architecture](PROJECT_STRUCTURE.md)** - How everything works

---

## Full Setup Guide

### 1. Setup Development Environment

On your Debian 13 server:

```bash
# Download and run the setup script
curl -fsSL https://raw.githubusercontent.com/callmemars1/superfly/main/dev-setup.sh -o dev-setup.sh
chmod +x dev-setup.sh
./dev-setup.sh
```

Or if you already have the repo:

```bash
chmod +x dev-setup.sh
./dev-setup.sh
```

This will install:
- ✅ K3S (Kubernetes)
- ✅ Traefik (Ingress Controller)
- ✅ cert-manager (Let's Encrypt)
- ✅ PostgreSQL (in K3S)
- ✅ Container Registry (in K3S)
- ✅ Go 1.22
- ✅ Node.js 20
- ✅ Development tools (sqlc, goose, air)

**Note:** If Docker was installed, you'll need to log out and back in for group permissions.

### 2. Verify Installation

```bash
chmod +x verify-setup.sh
./verify-setup.sh
```

All checks should pass ✅

### 3. Access Services

After setup:

**PostgreSQL:**
```bash
# Connection string
postgresql://superfly:superfly_dev_password@localhost:5432/superfly

# Connect with psql
psql postgresql://superfly:superfly_dev_password@localhost:5432/superfly
```

**Container Registry:**
```bash
# From within K8S cluster
registry.superfly-system.svc.cluster.local:5000
```

**Kubernetes:**
```bash
kubectl get pods -n superfly-system
kubectl get pods -n superfly-apps
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      K3S Cluster                            │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Traefik Ingress (with Let's Encrypt)               │  │
│  └────────┬──────────────────────────────┬──────────────┘  │
│           │                               │                 │
│  ┌────────▼────────┐          ┌──────────▼──────────┐     │
│  │  User Apps      │          │  Control Plane      │     │
│  │  (superfly-apps)│          │  (superfly-system)  │     │
│  │                 │          │  - API Server       │     │
│  │  - app1         │          │  - PostgreSQL       │     │
│  │  - app2         │          │  - Registry         │     │
│  │  - app3         │          │                     │     │
│  └─────────────────┘          └─────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

## Tech Stack

### Backend (Control Plane)
- **Language:** Go 1.22
- **Database:** PostgreSQL 16
- **ORM:** sqlc + pgx v5
- **K8S Client:** client-go

### Frontend
- **Framework:** SvelteKit
- **UI Library:** shadcn-svelte or Skeleton UI

### Infrastructure
- **Orchestration:** K3S
- **Ingress:** Traefik
- **Certificates:** cert-manager + Let's Encrypt
- **Registry:** Docker Registry v2
- **Observability:** Prometheus + Loki + Grafana (coming soon)

## Project Structure

```
superfly/
├── cmd/
│   ├── api/              # Control plane API server
│   ├── installer/        # Migration runner
│   └── cli/              # CLI tool (future)
│
├── internal/
│   ├── service/          # Business logic
│   ├── k8s/              # Kubernetes helpers
│   ├── handlers/         # HTTP handlers
│   └── middleware/       # HTTP middleware
│
├── db/
│   ├── migrations/       # SQL migrations (goose)
│   ├── queries/          # SQL queries (sqlc)
│   └── sqlc.yaml         # sqlc configuration
│
├── web/                  # Svelte frontend
│   ├── src/
│   ├── package.json
│   └── Dockerfile
│
├── manifests/            # K8s manifests for system
│   ├── postgres.yaml
│   ├── registry.yaml
│   └── ...
│
├── dev-setup.sh          # Development environment setup
├── verify-setup.sh       # Verify environment
└── README.md
```

## Development Workflow

### Start PostgreSQL Port Forward (if not running)
```bash
kubectl port-forward -n superfly-system svc/postgres 5432:5432
```

### Run Database Migrations
```bash
# Create a new migration
goose -dir db/migrations create add_apps_table sql

# Run migrations
goose -dir db/migrations postgres "postgresql://superfly:superfly_dev_password@localhost:5432/superfly" up

# Rollback
goose -dir db/migrations postgres "postgresql://superfly:superfly_dev_password@localhost:5432/superfly" down
```

### Generate Code with sqlc
```bash
sqlc generate
```

### Run API Server (with live reload)
```bash
cd cmd/api
air  # or: go run main.go
```

### Run Frontend
```bash
cd web
npm install
npm run dev
```

## MVP Features

### Feature 1: Deploy Pre-built Docker Image ✅ (Current)
Deploy an app from an existing Docker image with automatic TLS.

**API:**
```bash
POST /api/apps
{
  "name": "My App",
  "image": "nginx:latest",
  "port": 80,
  "domain": "myapp.example.com"
}
```

### Feature 2: Build from GitHub + Dockerfile
Build Docker images from GitHub repositories using Kaniko.

### Feature 3: Environment Variables
Manage environment variables and secrets per app.

### Feature 4: View Logs (Loki)
Real-time log viewing from Loki.

### Feature 5: View Metrics (Prometheus + Grafana)
Monitor CPU, memory, and request metrics.

### Feature 6: GitHub Webhooks
Auto-deploy on git push.

### Feature 7: Multiple Domains
Support multiple domains per app.

### Feature 8: Web UI
Full-featured web dashboard.

## Useful Commands

### Kubernetes
```bash
# Get all resources in superfly-system
kubectl get all -n superfly-system

# Get all apps
kubectl get all -n superfly-apps

# View logs
kubectl logs -n superfly-system deployment/postgres
kubectl logs -n superfly-apps deployment/my-app

# Describe resource
kubectl describe ingress -n superfly-apps my-app

# Get ingress with IPs
kubectl get ingress -A

# Port forward
kubectl port-forward -n superfly-system svc/postgres 5432:5432
```

### Database
```bash
# Connect to PostgreSQL
psql postgresql://superfly:superfly_dev_password@localhost:5432/superfly

# List tables
\dt

# Describe table
\d apps

# Query
SELECT * FROM apps;
```

### Testing
```bash
# Test deploying nginx
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Nginx",
    "image": "nginx:latest",
    "port": 80,
    "domain": "test.example.com"
  }'

# List apps
curl http://localhost:8080/api/apps

# Get app details
curl http://localhost:8080/api/apps/{id}

# Delete app
curl -X DELETE http://localhost:8080/api/apps/{id}
```

## Environment Variables

See `.env.example` for all available environment variables.

Key variables:
- `DATABASE_URL` - PostgreSQL connection string
- `KUBECONFIG` - Path to kubeconfig (for local dev)
- `KUBERNETES_IN_CLUSTER` - Set to `true` when running in K8S
- `REGISTRY_URL` - Container registry URL
- `API_PORT` - API server port (default: 8080)

## Troubleshooting

### PostgreSQL connection refused
```bash
# Check if port forward is running
sudo systemctl status superfly-postgres-forward.service

# Restart port forward
sudo systemctl restart superfly-postgres-forward.service

# Or manually forward
kubectl port-forward -n superfly-system svc/postgres 5432:5432
```

### K3S not responding
```bash
# Check K3S status
sudo systemctl status k3s

# Restart K3S
sudo systemctl restart k3s

# View logs
sudo journalctl -u k3s -f
```

### Pods not starting
```bash
# Check pod status
kubectl get pods -n superfly-system

# Describe pod
kubectl describe pod <pod-name> -n superfly-system

# View logs
kubectl logs <pod-name> -n superfly-system
```

## License

MIT

## Links

- **Repository**: https://github.com/callmemars1/superfly
- **Domain**: superfly.smartynov.com
- **Documentation**: See guides above

## Contributing

This is a work in progress. Contributions welcome!

## License

MIT
