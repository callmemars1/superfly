# Step-by-Step Deployment Guide

Complete guide to deploy apps with Superfly on a remote Debian 13 server.

---

## Prerequisites

- ✅ Debian 13 server (fresh install recommended)
- ✅ Root or sudo access
- ✅ At least 2GB RAM, 20GB disk space
- ✅ SSH access to the server
- ✅ Domain pointing to your server's IP (optional, for public access)

---

## Part 1: Server Preparation (10 minutes)

### Step 1: SSH into Your Server

```bash
# From your local machine
ssh root@your-server-ip

# Or if using a non-root user
ssh your-user@your-server-ip
```

### Step 2: Update System

```bash
sudo apt update
sudo apt upgrade -y
```

### Step 3: Clone Superfly Repository

```bash
cd ~
git clone https://github.com/callmemars1/superfly.git
cd superfly
```

**OR** if you're developing, copy files to the server:

```bash
# From your local machine
scp -r /path/to/superfly root@your-server-ip:~/
```

### Step 4: Make Scripts Executable

```bash
chmod +x dev-setup.sh
chmod +x verify-setup.sh
chmod +x test-api.sh
```

---

## Part 2: Install Superfly Environment (5 minutes)

### Step 5: Run Setup Script

```bash
./dev-setup.sh
```

**What this does** (takes 5-10 minutes):
- ✓ Installs K3S (Kubernetes)
- ✓ Installs Traefik (ingress controller)
- ✓ Installs cert-manager (for Let's Encrypt)
- ✓ Deploys PostgreSQL database
- ✓ Deploys container registry
- ✓ Installs Go 1.22
- ✓ Installs Node.js 20
- ✓ Installs development tools
- ✓ Sets up port forwarding for PostgreSQL

**Expected output** (at the end):
```
==================================
✅ Setup Complete!
==================================

✓ K3S cluster running
✓ PostgreSQL: localhost:5432
  - Database: superfly
  - User: superfly
  - Password: superfly_dev_password
✓ Container Registry: registry.superfly-system.svc.cluster.local:5000
...
```

### Step 6: Verify Installation

```bash
./verify-setup.sh
```

**All checks should pass** ✅

If any check fails, see troubleshooting section below.

---

## Part 3: Build and Start Superfly API (3 minutes)

### Step 7: Initialize Go Project

```bash
# Install Go dependencies
make init
go mod tidy
```

**Expected output**:
```
go: downloading github.com/go-chi/chi/v5 v5.0.11
go: downloading github.com/jackc/pgx/v5 v5.5.1
...
```

### Step 8: Setup Database

```bash
# Run migrations (creates tables)
make migrate
```

**Expected output**:
```
OK    00001_create_apps_table.sql (123.45ms)
goose: no migrations to run. current version: 1
```

### Step 9: Generate Database Code

```bash
# Generate Go code from SQL queries
make sqlc-generate
```

This creates files in `internal/db/`:
- `db.go`
- `models.go`
- `apps.sql.go`

### Step 10: Build the API Server

```bash
make build
```

**Expected output**:
```
Building API server...
```

Binary will be at: `bin/superfly-api`

### Step 11: Start the API Server

For development (with auto-reload):
```bash
make dev
```

For production (runs in background):
```bash
# Create systemd service (recommended)
sudo tee /etc/systemd/system/superfly-api.service > /dev/null <<EOF
[Unit]
Description=Superfly API Server
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$HOME/superfly
ExecStart=$HOME/superfly/bin/superfly-api
Restart=always
RestartSec=10
Environment="DATABASE_URL=postgresql://superfly:superfly_dev_password@localhost:5432/superfly?sslmode=disable"
Environment="KUBERNETES_IN_CLUSTER=false"
Environment="KUBECONFIG=$HOME/.kube/config"

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable superfly-api
sudo systemctl start superfly-api
```

**Check if running**:
```bash
sudo systemctl status superfly-api

# Or check logs
sudo journalctl -u superfly-api -f
```

**Test health endpoint**:
```bash
curl http://localhost:8080/health
```

**Expected response**:
```json
{"status":"ok","version":"0.1.0"}
```

✅ **API server is now running!**

---

## Part 4: Deploy Your First App (2 minutes)

### Step 12: Deploy Nginx Test App

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Nginx",
    "image": "nginx:alpine",
    "port": 80,
    "domain": "test.superfly.smartynov.com"
  }'
```

**Response** (save the `id`):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "slug": "test-nginx",
  "name": "Test Nginx",
  "image": "nginx:alpine",
  "port": 80,
  "status": "pending",
  "domain": "test.yourdomain.com",
  ...
}
```

### Step 13: Watch Deployment

**Check deployment status**:
```bash
# Watch pods starting
kubectl get pods -n superfly-apps -w

# Check all resources
kubectl get all -n superfly-apps

# View deployment details
kubectl describe deployment test-nginx -n superfly-apps
```

**Wait for pod to be Ready**:
```
NAME                          READY   STATUS    RESTARTS   AGE
test-nginx-7d4b4c9f6d-abcde   1/1     Running   0          45s
```

### Step 14: Check App Status via API

```bash
# Replace with your app ID
APP_ID="550e8400-e29b-41d4-a716-446655440000"

curl http://localhost:8080/api/apps/$APP_ID
```

**Status should be "running"**:
```json
{
  "status": "running",
  "last_deployed_at": "2026-01-14T12:30:00Z",
  ...
}
```

### Step 15: View App Logs

```bash
kubectl logs -n superfly-apps deployment/test-nginx -f
```

**Expected output**:
```
/docker-entrypoint.sh: Configuration complete; ready for start up
```

---

## Part 5: Access Your App (2 minutes)

### Option A: With Domain (Production)

**Prerequisites**:
1. You have a domain (e.g., `yourdomain.com`)
2. DNS A record points to your server's IP

**Example DNS setup**:
```
Type: A
Name: test
Value: 123.45.67.89 (your server IP)
TTL: 3600
```

**Wait for DNS to propagate** (5-30 minutes):
```bash
# Check DNS
dig test.yourdomain.com

# Should return your server IP
```

**Access your app**:
```bash
# Test locally first
curl http://test.yourdomain.com

# With TLS (after cert-manager gets certificate)
curl https://test.yourdomain.com
```

**Check certificate status**:
```bash
kubectl get certificate -n superfly-apps
kubectl describe certificate test-nginx-tls -n superfly-apps
```

**Certificate should be Ready**:
```
NAME              READY   SECRET            AGE
test-nginx-tls    True    test-nginx-tls    2m
```

🎉 **Your app is now live with HTTPS!**

### Option B: Without Domain (Testing)

If you don't have a domain, use port forwarding:

```bash
# Port forward the service
kubectl port-forward -n superfly-apps svc/test-nginx 8080:80
```

**Access in browser or curl**:
```bash
curl http://localhost:8080
```

### Option C: Direct Server IP (Quick Test)

Get the Traefik LoadBalancer IP:

```bash
kubectl get svc -n traefik
```

**Access directly** (without domain):
```bash
curl http://YOUR_SERVER_IP
```

---

## Part 6: Deploy a Real App (5 minutes)

### Example: Deploy a Node.js API

**Assuming you have a Docker image**:

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My API",
    "image": "your-dockerhub-username/your-api:latest",
    "port": 3000,
    "domain": "api.yourdomain.com",
    "replicas": 2,
    "cpu_limit": "1000m",
    "memory_limit": "512Mi",
    "health_check_path": "/health"
  }'
```

**Watch deployment**:
```bash
kubectl get pods -n superfly-apps -w
```

**Check logs**:
```bash
kubectl logs -n superfly-apps deployment/my-api -f
```

**Test endpoint**:
```bash
curl https://api.yourdomain.com
```

---

## Part 7: Manage Apps

### List All Apps

```bash
curl http://localhost:8080/api/apps | jq
```

### Get App Details

```bash
curl http://localhost:8080/api/apps/$APP_ID | jq
```

### Update App (Scale to 3 replicas)

```bash
curl -X PATCH http://localhost:8080/api/apps/$APP_ID \
  -H "Content-Type: application/json" \
  -d '{"replicas": 3}'
```

### Update App (Change Image)

```bash
curl -X PATCH http://localhost:8080/api/apps/$APP_ID \
  -H "Content-Type: application/json" \
  -d '{"image": "nginx:latest"}'
```

### Restart App

```bash
curl -X POST http://localhost:8080/api/apps/$APP_ID/restart
```

### Delete App

```bash
curl -X DELETE http://localhost:8080/api/apps/$APP_ID
```

---

## Part 8: Monitoring & Troubleshooting

### Check System Status

```bash
# All pods in superfly-system
kubectl get pods -n superfly-system

# All deployed apps
kubectl get pods -n superfly-apps

# All ingresses
kubectl get ingress -A
```

### View Logs

```bash
# API server logs
sudo journalctl -u superfly-api -f

# App logs
kubectl logs -n superfly-apps deployment/APP-SLUG -f

# PostgreSQL logs
kubectl logs -n superfly-system deployment/postgres -f

# Traefik logs
kubectl logs -n traefik deployment/traefik -f
```

### Check Resource Usage

```bash
# Node resources
kubectl top nodes

# Pod resources
kubectl top pods -n superfly-apps
```

### Database Operations

```bash
# Connect to database
psql postgresql://superfly:superfly_dev_password@localhost:5432/superfly

# List apps
SELECT id, slug, name, status, created_at FROM apps;

# Check app count
SELECT COUNT(*) FROM apps;
```

---

## Common Issues & Solutions

### Issue 1: Port 5432 Already in Use

**Solution**:
```bash
# Check what's using port 5432
sudo lsof -i :5432

# If it's PostgreSQL, stop it
sudo systemctl stop postgresql

# Restart port forward
sudo systemctl restart superfly-postgres-forward
```

### Issue 2: Pod Stuck in "Pending"

**Check**:
```bash
kubectl describe pod -n superfly-apps POD-NAME
```

**Common causes**:
- Insufficient resources
- Image pull errors
- Storage issues

**Solution**:
```bash
# Check node resources
kubectl describe nodes

# Check events
kubectl get events -n superfly-apps --sort-by='.lastTimestamp'
```

### Issue 3: Pod Stuck in "ImagePullBackOff"

**Cause**: Can't pull Docker image

**Solution**:
- Verify image exists: `docker pull IMAGE:TAG`
- Check image name spelling
- For private registries, add credentials

### Issue 4: App Shows "deploying" for Too Long

**Check**:
```bash
# Pod status
kubectl get pods -n superfly-apps -l app=YOUR-SLUG

# Pod events
kubectl describe pod -n superfly-apps POD-NAME

# Pod logs
kubectl logs -n superfly-apps POD-NAME
```

**Common causes**:
- Health check failing
- App crashing on startup
- Wrong port configuration

### Issue 5: Certificate Not Issuing

**Check**:
```bash
kubectl get certificate -n superfly-apps
kubectl describe certificate CERT-NAME -n superfly-apps
kubectl get challenges -A
```

**Common causes**:
- DNS not pointing to server
- Ports 80/443 not accessible
- Rate limit hit (use staging issuer first)

**Solution**:
```bash
# Check DNS
dig YOUR-DOMAIN

# Test port 80 accessibility
curl http://YOUR-DOMAIN

# View cert-manager logs
kubectl logs -n cert-manager deployment/cert-manager -f
```

### Issue 6: Can't Access API Server

**Check**:
```bash
# Is it running?
sudo systemctl status superfly-api

# Check logs
sudo journalctl -u superfly-api -f

# Test locally
curl http://localhost:8080/health
```

**Solution**:
```bash
# Restart API
sudo systemctl restart superfly-api

# Check firewall
sudo ufw status
sudo ufw allow 8080/tcp
```

---

## Production Checklist

Before going to production:

- [ ] Change PostgreSQL password
- [ ] Use production Let's Encrypt issuer (not staging)
- [ ] Configure firewall (ufw)
- [ ] Setup backups for PostgreSQL
- [ ] Configure monitoring (Feature 5)
- [ ] Setup log aggregation (Feature 4)
- [ ] Configure resource limits appropriately
- [ ] Setup alerts for pod failures
- [ ] Document your domains
- [ ] Setup DNS properly
- [ ] Test disaster recovery
- [ ] Configure auto-scaling (if needed)

---

## Security Hardening

### 1. Change Default Passwords

```bash
# Update PostgreSQL password
kubectl exec -it -n superfly-system postgres-xxx -- psql -U superfly
ALTER USER superfly WITH PASSWORD 'new_secure_password';

# Update .env file
nano .env
# Change DATABASE_URL password
```

### 2. Configure Firewall

```bash
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 6443/tcp  # K3s API
sudo ufw enable
```

### 3. Restrict API Access

```bash
# Only allow localhost
sudo ufw allow from 127.0.0.1 to any port 8080

# Or specific IP
sudo ufw allow from YOUR_IP to any port 8080
```

### 4. Enable HTTPS Only

Update Ingress to redirect HTTP to HTTPS (already configured).

---

## Performance Tips

### 1. Resource Limits

Set appropriate limits for your apps:
- Small apps: 250m CPU, 128Mi memory
- Medium apps: 500m CPU, 512Mi memory
- Large apps: 1000m+ CPU, 1Gi+ memory

### 2. Replicas

Scale apps based on traffic:
- Low traffic: 1-2 replicas
- Medium traffic: 3-5 replicas
- High traffic: 5+ replicas

### 3. Health Checks

Use lightweight health check endpoints:
- Don't check database in health checks
- Use `/health` or `/ping` endpoints
- Return 200 quickly

---

## Quick Command Reference

```bash
# Check system
kubectl get pods -A
kubectl get all -n superfly-apps

# Deploy app
curl -X POST http://localhost:8080/api/apps -d '{...}'

# List apps
curl http://localhost:8080/api/apps | jq

# Update app
curl -X PATCH http://localhost:8080/api/apps/ID -d '{...}'

# Delete app
curl -X DELETE http://localhost:8080/api/apps/ID

# View logs
kubectl logs -n superfly-apps deployment/SLUG -f

# Port forward
kubectl port-forward -n superfly-apps svc/SLUG 8080:80

# Database
psql postgresql://superfly:PASSWORD@localhost:5432/superfly

# Restart API
sudo systemctl restart superfly-api

# View API logs
sudo journalctl -u superfly-api -f
```

---

## Next Steps

1. ✅ **Deploy your first app** - You just did this!
2. 📚 **Read the examples** - See [EXAMPLES.md](EXAMPLES.md)
3. 🔨 **Build Feature 2** - GitHub + Dockerfile builds
4. 🎨 **Create Web UI** - Svelte dashboard
5. 📊 **Add observability** - Logs + Metrics

---

## Success! 🎉

You now have:
- ✅ Kubernetes cluster running
- ✅ API server accepting requests
- ✅ Apps deploying with one command
- ✅ Automatic HTTPS certificates
- ✅ Production-ready infrastructure

**Deploy anything!** 🚀
