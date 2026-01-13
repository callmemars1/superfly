# Superfly API Documentation

Base URL: `http://localhost:8080`

## Endpoints

### Health Check

#### GET /health

Check if the API server is running.

**Response**
```json
{
  "status": "ok",
  "version": "0.1.0"
}
```

---

### Apps

#### POST /api/apps

Create a new app and deploy it to Kubernetes.

**Request Body**
```json
{
  "name": "My App",               // Required: Display name
  "slug": "my-app",               // Optional: URL-safe name (auto-generated if not provided)
  "image": "nginx:alpine",        // Required: Docker image
  "port": 80,                     // Optional: Container port (default: 8080)
  "replicas": 1,                  // Optional: Number of replicas (default: 1)
  "cpu_limit": "500m",            // Optional: CPU limit (default: 500m)
  "memory_limit": "256Mi",        // Optional: Memory limit (default: 256Mi)
  "domain": "example.com",        // Optional: Domain for ingress
  "health_check_path": "/"        // Optional: Health check path (default: /)
}
```

**Response** (201 Created)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "slug": "my-app",
  "name": "My App",
  "image": "nginx:alpine",
  "port": 80,
  "replicas": 1,
  "cpu_limit": "500m",
  "memory_limit": "256Mi",
  "domain": "example.com",
  "health_check_path": "/",
  "status": "pending",
  "created_at": "2026-01-14T10:30:00Z",
  "updated_at": "2026-01-14T10:30:00Z",
  "last_deployed_at": null
}
```

**Status Values**
- `pending` - App created, not yet deploying
- `deploying` - Currently deploying to Kubernetes
- `running` - Successfully deployed and running
- `failed` - Deployment failed

**Example**
```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Nginx Web Server",
    "image": "nginx:alpine",
    "port": 80,
    "domain": "web.example.com"
  }'
```

---

#### GET /api/apps

List all apps.

**Response** (200 OK)
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "slug": "my-app",
    "name": "My App",
    "image": "nginx:alpine",
    "port": 80,
    "status": "running",
    ...
  }
]
```

**Example**
```bash
curl http://localhost:8080/api/apps
```

---

#### GET /api/apps/:id

Get details of a specific app.

**Parameters**
- `id` (UUID) - App ID

**Response** (200 OK)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "slug": "my-app",
  "name": "My App",
  "image": "nginx:alpine",
  "port": 80,
  "replicas": 1,
  "cpu_limit": "500m",
  "memory_limit": "256Mi",
  "domain": "example.com",
  "health_check_path": "/",
  "status": "running",
  "created_at": "2026-01-14T10:30:00Z",
  "updated_at": "2026-01-14T10:35:00Z",
  "last_deployed_at": "2026-01-14T10:35:00Z"
}
```

**Example**
```bash
curl http://localhost:8080/api/apps/550e8400-e29b-41d4-a716-446655440000
```

---

#### PATCH /api/apps/:id

Update an app. Provide only the fields you want to update.

**Parameters**
- `id` (UUID) - App ID

**Request Body**
```json
{
  "name": "New Name",             // Optional
  "image": "nginx:latest",        // Optional (triggers redeploy)
  "port": 8080,                   // Optional (triggers redeploy)
  "replicas": 2,                  // Optional (triggers redeploy)
  "cpu_limit": "1000m",           // Optional (triggers redeploy)
  "memory_limit": "512Mi",        // Optional (triggers redeploy)
  "domain": "newdomain.com",      // Optional (triggers redeploy)
  "health_check_path": "/health"  // Optional (triggers redeploy)
}
```

**Response** (200 OK)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "slug": "my-app",
  "name": "New Name",
  "image": "nginx:latest",
  "port": 8080,
  "replicas": 2,
  ...
}
```

**Example: Scale to 3 replicas**
```bash
curl -X PATCH http://localhost:8080/api/apps/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"replicas": 3}'
```

**Example: Update image**
```bash
curl -X PATCH http://localhost:8080/api/apps/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"image": "nginx:1.25"}'
```

---

#### POST /api/apps/:id/restart

Trigger a rolling restart of the app.

**Parameters**
- `id` (UUID) - App ID

**Response** (200 OK)
```json
{
  "message": "App restart initiated"
}
```

**Example**
```bash
curl -X POST http://localhost:8080/api/apps/550e8400-e29b-41d4-a716-446655440000/restart
```

---

#### DELETE /api/apps/:id

Delete an app and all its Kubernetes resources.

**Parameters**
- `id` (UUID) - App ID

**Response** (204 No Content)

**Example**
```bash
curl -X DELETE http://localhost:8080/api/apps/550e8400-e29b-41d4-a716-446655440000
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message"
}
```

### HTTP Status Codes

- `200 OK` - Request succeeded
- `201 Created` - Resource created
- `204 No Content` - Resource deleted
- `400 Bad Request` - Invalid input
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

### Common Errors

**Invalid app ID**
```json
{
  "error": "Invalid app ID"
}
```

**App not found**
```json
{
  "error": "App not found"
}
```

**Slug already exists**
```json
{
  "error": "app with slug 'my-app' already exists"
}
```

**Domain already in use**
```json
{
  "error": "domain 'example.com' already in use"
}
```

**Invalid slug**
```json
{
  "error": "slug must consist of lowercase alphanumeric characters or '-', and must start and end with an alphanumeric character"
}
```

---

## Resource Limits

### CPU Limits
Format: `<number>m` (millicores)
- `100m` = 0.1 CPU core
- `500m` = 0.5 CPU core
- `1000m` = 1 CPU core
- `2000m` = 2 CPU cores

### Memory Limits
Format: `<number>Mi` or `<number>Gi`
- `128Mi` = 128 megabytes
- `256Mi` = 256 megabytes
- `512Mi` = 512 megabytes
- `1Gi` = 1 gigabyte

---

## Kubernetes Resources Created

When you create an app, Superfly creates these Kubernetes resources:

### Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: superfly-apps
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  template:
    spec:
      containers:
      - name: app
        image: nginx:alpine
        ports:
        - containerPort: 80
        resources:
          limits:
            cpu: "500m"
            memory: "256Mi"
        livenessProbe:
          httpGet:
            path: /
            port: 80
        readinessProbe:
          httpGet:
            path: /
            port: 80
```

### Service
```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app
  namespace: superfly-apps
spec:
  selector:
    app: my-app
  ports:
  - port: 80
    targetPort: 80
```

### Ingress
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-app
  namespace: superfly-apps
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    traefik.ingress.kubernetes.io/router.tls: "true"
spec:
  tls:
  - hosts:
    - example.com
    secretName: my-app-tls
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        backend:
          service:
            name: my-app
            port:
              number: 80
```

---

## Advanced Examples

### Deploy Postgres Database

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PostgreSQL",
    "image": "postgres:16-alpine",
    "port": 5432,
    "cpu_limit": "1000m",
    "memory_limit": "1Gi"
  }'
```

### Deploy Redis

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Redis Cache",
    "image": "redis:7-alpine",
    "port": 6379,
    "cpu_limit": "500m",
    "memory_limit": "512Mi"
  }'
```

### Deploy React App (via nginx)

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My React App",
    "image": "myregistry/react-app:latest",
    "port": 80,
    "domain": "app.example.com",
    "cpu_limit": "250m",
    "memory_limit": "128Mi"
  }'
```

### High-Performance App

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "High Performance API",
    "image": "myapi:latest",
    "port": 8080,
    "replicas": 3,
    "cpu_limit": "2000m",
    "memory_limit": "2Gi",
    "domain": "api.example.com",
    "health_check_path": "/health"
  }'
```

---

## Accessing Your Apps

### Within Kubernetes Cluster

Apps can call each other using Kubernetes DNS:

```
http://<app-slug>.superfly-apps.svc.cluster.local
```

Example:
```bash
# From another pod in the cluster
curl http://my-app.superfly-apps.svc.cluster.local
```

### From Internet (with domain)

If you configured a domain, the app is accessible at:

```
https://<your-domain>
```

Make sure your domain's DNS points to your server's IP address.

### Port Forwarding (testing)

For testing without a domain:

```bash
kubectl port-forward -n superfly-apps svc/my-app 8080:80
```

Then access at `http://localhost:8080`

---

## Monitoring Deployment Status

### Via API

Poll the app endpoint to check status:

```bash
watch -n 2 'curl -s http://localhost:8080/api/apps/550e8400-e29b-41d4-a716-446655440000 | jq .status'
```

### Via Kubernetes

```bash
# Watch pods
kubectl get pods -n superfly-apps -w

# Check deployment status
kubectl rollout status deployment/my-app -n superfly-apps

# View events
kubectl get events -n superfly-apps --sort-by='.lastTimestamp'
```

---

## Tips & Best Practices

1. **Use specific image tags** instead of `latest` for production
2. **Set appropriate resource limits** to prevent resource exhaustion
3. **Use health check paths** that don't have external dependencies
4. **Start with 1 replica** and scale up as needed
5. **Use unique slugs** to avoid conflicts
6. **Configure domains** only after DNS is properly set up

---

## Coming Soon

- ⏳ Environment variables
- ⏳ Build from GitHub repositories
- ⏳ View logs via API
- ⏳ View metrics via API
- ⏳ GitHub webhooks (auto-deploy on push)
- ⏳ Multiple domains per app
- ⏳ Persistent volumes
- ⏳ Cron jobs
