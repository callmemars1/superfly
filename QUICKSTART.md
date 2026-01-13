# Superfly Quick Start

Get Superfly running in 10 minutes!

## Step 1: Setup Environment (5 minutes)

On your Debian 13 server:

```bash
chmod +x dev-setup.sh
./dev-setup.sh
```

Wait for the script to complete. It will install:
- K3S cluster
- PostgreSQL
- Traefik
- cert-manager
- Container Registry
- Go, Node.js, and dev tools

## Step 2: Verify Installation (30 seconds)

```bash
chmod +x verify-setup.sh
./verify-setup.sh
```

All checks should pass ✅

## Step 3: Initialize Project (1 minute)

```bash
# Install Go dependencies
make init
go mod tidy

# Run database migrations
make migrate

# Generate Go code from SQL
make sqlc-generate
```

## Step 4: Start API Server (10 seconds)

```bash
make dev
```

You should see:
```
🚀 Server listening on 0.0.0.0:8080
📝 API documentation: http://0.0.0.0:8080/api
```

## Step 5: Deploy Your First App (1 minute)

In a new terminal, test the API:

```bash
# Deploy nginx
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My First App",
    "image": "nginx:alpine",
    "port": 80,
    "domain": "app.example.com"
  }'
```

Response:
```json
{
  "id": "550e8400-...",
  "slug": "my-first-app",
  "name": "My First App",
  "status": "deploying",
  ...
}
```

## Step 6: Check Deployment

```bash
# Watch pods
kubectl get pods -n superfly-apps -w

# View deployment
kubectl get all -n superfly-apps

# Check app logs
kubectl logs -n superfly-apps deployment/my-first-app -f
```

## Step 7: Run Full Test Suite

```bash
chmod +x test-api.sh
./test-api.sh
```

## What You Just Built

✅ Full Kubernetes cluster (K3S)  
✅ REST API for app deployment  
✅ Automatic TLS certificates (Let's Encrypt)  
✅ Traefik ingress routing  
✅ PostgreSQL database  
✅ Container registry  

## API Endpoints

```bash
# Health check
GET  /health

# List apps
GET  /api/apps

# Create app
POST /api/apps
{
  "name": "App Name",
  "image": "docker/image:tag",
  "port": 8080,
  "domain": "example.com"
}

# Get app
GET  /api/apps/{id}

# Update app
PATCH /api/apps/{id}
{
  "replicas": 2
}

# Restart app
POST /api/apps/{id}/restart

# Delete app
DELETE /api/apps/{id}
```

## Example: Deploy a Real App

### Node.js App

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Express API",
    "image": "node:20-alpine",
    "port": 3000,
    "domain": "api.myapp.com"
  }'
```

### Python App

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Flask App",
    "image": "python:3.11-slim",
    "port": 5000,
    "domain": "flask.myapp.com"
  }'
```

### Static Site

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Website",
    "image": "nginx:alpine",
    "port": 80,
    "domain": "www.mysite.com",
    "cpu_limit": "250m",
    "memory_limit": "128Mi"
  }'
```

## Troubleshooting

### PostgreSQL not connecting?

```bash
# Restart port forward
sudo systemctl restart superfly-postgres-forward.service

# Or manually
kubectl port-forward -n superfly-system svc/postgres 5432:5432
```

### K3S not responding?

```bash
sudo systemctl status k3s
sudo systemctl restart k3s
```

### App stuck in deploying?

```bash
# Check pods
kubectl get pods -n superfly-apps

# Describe pod
kubectl describe pod -n superfly-apps <pod-name>

# View logs
kubectl logs -n superfly-apps <pod-name>
```

## Next Steps

📚 **Read Development Guide**: See [DEVELOPMENT.md](DEVELOPMENT.md)

🔨 **Build Features**: 
- Feature 2: Build from GitHub
- Feature 3: Environment Variables
- Feature 4: Logs with Loki
- Feature 5: Metrics with Prometheus

🎨 **Create Web UI**: Svelte dashboard

## Useful Commands

```bash
# Database
make db-shell              # Open PostgreSQL
make migrate              # Run migrations
make db-reset             # Reset database

# Development
make dev                  # Run with live reload
make build               # Build binary
make test                # Run tests

# Kubernetes
kubectl get all -n superfly-apps        # View apps
kubectl logs -n superfly-apps -l app=X  # View logs
kubectl describe pod -n superfly-apps X # Debug pod
```

## Architecture

```
┌─────────────────────────────────────────┐
│  Internet → Traefik (Port 80/443)      │
└─────────────┬───────────────────────────┘
              │
         ┌────▼─────┐
         │ Ingress  │ (TLS via Let's Encrypt)
         └────┬─────┘
              │
    ┌─────────┴──────────┐
    │                    │
┌───▼────┐        ┌──────▼────┐
│ App 1  │        │ App 2     │
│ nginx  │        │ Your App  │
└────────┘        └───────────┘
```

## Success! 🎉

You now have a working homelab PaaS platform!

Deploy apps with a single API call and get:
- Automatic HTTPS
- Health checks
- Rolling updates
- Resource limits
- Multi-app networking

Happy deploying! 🚀
