# Superfly Examples

Real-world examples of deploying apps with Superfly.

## Example 1: Static Website (Nginx)

Deploy a simple nginx web server:

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

Check status:
```bash
kubectl get pods -n superfly-apps -l app=my-website
kubectl logs -n superfly-apps deployment/my-website -f
```

---

## Example 2: Node.js API

Deploy a Node.js Express API:

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Express API",
    "image": "node:20-alpine",
    "port": 3000,
    "domain": "api.myapp.com",
    "health_check_path": "/health",
    "cpu_limit": "500m",
    "memory_limit": "512Mi"
  }'
```

---

## Example 3: Python Flask App

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Flask Backend",
    "image": "python:3.11-slim",
    "port": 5000,
    "domain": "flask.myapp.com",
    "cpu_limit": "1000m",
    "memory_limit": "1Gi"
  }'
```

---

## Example 4: Redis Cache

Deploy Redis for your apps to use:

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

Access from other apps:
```bash
redis-cli -h redis-cache.superfly-apps.svc.cluster.local -p 6379
```

---

## Example 5: PostgreSQL Database

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PostgreSQL DB",
    "slug": "postgres-db",
    "image": "postgres:16-alpine",
    "port": 5432,
    "cpu_limit": "1000m",
    "memory_limit": "1Gi"
  }'
```

Connect from other apps:
```bash
postgresql://postgres@postgres-db.superfly-apps.svc.cluster.local:5432/postgres
```

---

## Example 6: Full Stack App (Frontend + API + Database)

### Step 1: Deploy Database

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "App Database",
    "slug": "app-db",
    "image": "postgres:16-alpine",
    "port": 5432,
    "memory_limit": "1Gi"
  }'
```

### Step 2: Deploy API

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "App API",
    "slug": "app-api",
    "image": "myregistry/app-api:latest",
    "port": 8080,
    "domain": "api.myapp.com",
    "health_check_path": "/health",
    "replicas": 2,
    "memory_limit": "512Mi"
  }'
```

API connects to database:
```
DATABASE_URL=postgresql://postgres@app-db.superfly-apps.svc.cluster.local:5432/appdb
```

### Step 3: Deploy Frontend

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "App Frontend",
    "slug": "app-frontend",
    "image": "myregistry/app-frontend:latest",
    "port": 80,
    "domain": "myapp.com",
    "memory_limit": "256Mi"
  }'
```

Frontend calls API:
```
API_URL=https://api.myapp.com
```

---

## Example 7: Scaling an App

Deploy initially with 1 replica:

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Scalable API",
    "image": "myapi:latest",
    "port": 8080,
    "domain": "api.example.com",
    "replicas": 1
  }'
```

Later, scale to 5 replicas:

```bash
curl -X PATCH http://localhost:8080/api/apps/{app-id} \
  -H "Content-Type: application/json" \
  -d '{
    "replicas": 5
  }'
```

Watch scaling:
```bash
kubectl get pods -n superfly-apps -l app=scalable-api -w
```

---

## Example 8: Updating an App

Deploy initial version:

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My API",
    "image": "myapi:v1.0.0",
    "port": 8080,
    "domain": "api.example.com"
  }'
```

Update to new version (rolling update):

```bash
curl -X PATCH http://localhost:8080/api/apps/{app-id} \
  -H "Content-Type: application/json" \
  -d '{
    "image": "myapi:v1.1.0"
  }'
```

Watch rolling update:
```bash
kubectl rollout status deployment/my-api -n superfly-apps
```

---

## Example 9: High Availability Setup

Deploy with multiple replicas and appropriate resources:

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "HA Production API",
    "image": "myapi:production",
    "port": 8080,
    "domain": "api.production.com",
    "replicas": 3,
    "cpu_limit": "2000m",
    "memory_limit": "2Gi",
    "health_check_path": "/health"
  }'
```

---

## Example 10: Microservices Architecture

### User Service

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "User Service",
    "slug": "user-service",
    "image": "myapp/user-service:latest",
    "port": 8001,
    "domain": "users.api.myapp.com"
  }'
```

### Order Service

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Order Service",
    "slug": "order-service",
    "image": "myapp/order-service:latest",
    "port": 8002,
    "domain": "orders.api.myapp.com"
  }'
```

### Payment Service

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Payment Service",
    "slug": "payment-service",
    "image": "myapp/payment-service:latest",
    "port": 8003,
    "domain": "payments.api.myapp.com"
  }'
```

Services communicate internally:
```bash
# Order service calls user service
http://user-service.superfly-apps.svc.cluster.local:8001/users/123

# Order service calls payment service
http://payment-service.superfly-apps.svc.cluster.local:8003/charge
```

---

## Example 11: Testing Without Domain

Deploy without a domain for testing:

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test App",
    "image": "nginx:alpine",
    "port": 80
  }'
```

Port forward to access:
```bash
kubectl port-forward -n superfly-apps svc/test-app 8080:80
# Access at http://localhost:8080
```

---

## Example 12: Low Resource App

Deploy a lightweight app with minimal resources:

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Lightweight Service",
    "image": "alpine:latest",
    "port": 8080,
    "cpu_limit": "100m",
    "memory_limit": "64Mi"
  }'
```

---

## Example 13: Restart an App

Sometimes you need to restart an app (e.g., to pick up config changes):

```bash
curl -X POST http://localhost:8080/api/apps/{app-id}/restart
```

This triggers a rolling restart (zero downtime).

---

## Example 14: Delete an App

```bash
curl -X DELETE http://localhost:8080/api/apps/{app-id}
```

This removes:
- Deployment
- Service
- Ingress
- Database record

---

## Common Patterns

### Pattern 1: Database + API + Worker

1. **Database**: `postgres-db`
2. **API**: Connects to `postgres-db.superfly-apps.svc.cluster.local`
3. **Worker**: Also connects to same database

All can share the same database using internal DNS.

### Pattern 2: Multiple Frontends, One API

1. **API**: `api.example.com`
2. **Web App**: `app.example.com` → calls API
3. **Admin Panel**: `admin.example.com` → calls API
4. **Mobile API**: Same backend API

### Pattern 3: Multi-tenant Apps

Deploy separate instances per tenant:

```bash
# Tenant A
curl -X POST http://localhost:8080/api/apps \
  -d '{"name":"Tenant A App","slug":"tenant-a","image":"myapp:latest","domain":"a.myapp.com"}'

# Tenant B
curl -X POST http://localhost:8080/api/apps \
  -d '{"name":"Tenant B App","slug":"tenant-b","image":"myapp:latest","domain":"b.myapp.com"}'
```

---

## Troubleshooting Examples

### Check Pod Status

```bash
kubectl get pods -n superfly-apps
kubectl describe pod -n superfly-apps my-app-xxx
```

### View Logs

```bash
kubectl logs -n superfly-apps deployment/my-app -f
```

### Check Resource Usage

```bash
kubectl top pods -n superfly-apps
```

### Debug Networking

```bash
# Test from another pod
kubectl run -it --rm debug --image=alpine --restart=Never -n superfly-apps -- sh
# Inside pod:
apk add curl
curl http://my-app.superfly-apps.svc.cluster.local
```

---

## Tips

1. **Use slugs wisely**: They become DNS names, keep them short and descriptive
2. **Start small**: Begin with minimal resources, scale up as needed
3. **Health checks**: Use endpoints that don't depend on external services
4. **Domains**: Set up DNS before deploying with domain configuration
5. **Rolling updates**: Update image to deploy new versions with zero downtime
6. **Internal networking**: Use K8s DNS for service-to-service communication
7. **Resource limits**: Set appropriate limits to prevent resource exhaustion

---

## Next: Build from Source

Currently, you deploy pre-built images. Coming soon:

```bash
# Feature 2: Build from GitHub
curl -X POST http://localhost:8080/api/apps \
  -d '{
    "name": "My App",
    "git_repo": "https://github.com/user/repo",
    "git_branch": "main",
    "dockerfile_path": "./Dockerfile"
  }'
```

Superfly will:
1. Clone repo
2. Build with Kaniko
3. Push to registry
4. Deploy

Stay tuned! 🚀
