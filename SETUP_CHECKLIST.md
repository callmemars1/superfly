# Superfly Setup Checklist

Complete setup checklist for **superfly.smartynov.com**

Repository: https://github.com/callmemars1/superfly

---

## Phase 1: DNS Configuration ⏱️ 5 minutes

### Configure DNS Records

Add these records to your domain registrar:

```
Type: A
Name: superfly
Value: [YOUR_SERVER_IP]
TTL: 3600

Type: A
Name: *.superfly
Value: [YOUR_SERVER_IP]
TTL: 3600
```

### Verify DNS

```bash
dig superfly.smartynov.com
dig test.superfly.smartynov.com
```

✅ DNS returns your server IP

---

## Phase 2: Server Preparation ⏱️ 2 minutes

### SSH into Server

```bash
ssh root@YOUR_SERVER_IP
```

### Update System

```bash
apt update && apt upgrade -y
```

✅ System updated

---

## Phase 3: Install Superfly ⏱️ 10 minutes

### Clone Repository

```bash
cd ~
git clone https://github.com/callmemars1/superfly.git
cd superfly
chmod +x *.sh
```

### Run Setup

```bash
./dev-setup.sh
```

**Wait for completion** (5-10 minutes)

### Verify Installation

```bash
./verify-setup.sh
```

✅ All checks pass

---

## Phase 4: Build Superfly ⏱️ 3 minutes

### Initialize Project

```bash
make init
go mod tidy
```

### Setup Database

```bash
make migrate
make sqlc-generate
```

### Build Binary

```bash
make build
```

✅ Binary created at `bin/superfly-api`

---

## Phase 5: Start API Server ⏱️ 2 minutes

### Create Systemd Service

```bash
sudo tee /etc/systemd/system/superfly-api.service > /dev/null <<'EOF'
[Unit]
Description=Superfly API Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/superfly
ExecStart=/root/superfly/bin/superfly-api
Restart=always
RestartSec=10
Environment="DATABASE_URL=postgresql://superfly:superfly_dev_password@localhost:5432/superfly?sslmode=disable"
Environment="KUBERNETES_IN_CLUSTER=false"
Environment="KUBECONFIG=/root/.kube/config"

[Install]
WantedBy=multi-user.target
EOF
```

### Start Service

```bash
sudo systemctl daemon-reload
sudo systemctl enable superfly-api
sudo systemctl start superfly-api
```

### Verify Running

```bash
sudo systemctl status superfly-api
curl http://localhost:8080/health
```

✅ Response: `{"status":"ok","version":"0.1.0"}`

---

## Phase 6: Deploy Test App ⏱️ 3 minutes

### Deploy Nginx

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test App",
    "image": "nginx:alpine",
    "port": 80,
    "domain": "test.superfly.smartynov.com"
  }'
```

### Watch Deployment

```bash
kubectl get pods -n superfly-apps -w
```

**Wait for**: `test-app-xxx  1/1  Running`

### Check Status

```bash
kubectl get all -n superfly-apps
kubectl get ingress -n superfly-apps
kubectl get certificate -n superfly-apps
```

✅ Pod running, Ingress created, Certificate ready

---

## Phase 7: Verify Access ⏱️ 2 minutes

### Test HTTP

```bash
curl http://test.superfly.smartynov.com
```

### Test HTTPS (after cert issued)

```bash
curl https://test.superfly.smartynov.com
```

✅ Returns nginx welcome page

---

## Phase 8: Deploy Real App ⏱️ 5 minutes

### Your App Example

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Production App",
    "image": "your-docker-image:tag",
    "port": 8080,
    "domain": "app.superfly.smartynov.com",
    "replicas": 2,
    "cpu_limit": "1000m",
    "memory_limit": "512Mi",
    "health_check_path": "/health"
  }'
```

### Monitor

```bash
kubectl get pods -n superfly-apps -w
kubectl logs -n superfly-apps deployment/my-production-app -f
```

✅ App running at https://app.superfly.smartynov.com

---

## Completion Checklist

### Infrastructure
- [ ] DNS configured (A records)
- [ ] DNS propagating/propagated
- [ ] Server accessible via SSH
- [ ] Ports 80, 443, 6443 open
- [ ] System updated

### Superfly Installation
- [ ] K3S installed and running
- [ ] PostgreSQL deployed
- [ ] Traefik deployed
- [ ] cert-manager deployed
- [ ] Container registry deployed
- [ ] Port forwarding configured

### Superfly Build
- [ ] Go dependencies installed
- [ ] Database migrations run
- [ ] Code generated with sqlc
- [ ] Binary compiled

### API Server
- [ ] Systemd service created
- [ ] Service enabled
- [ ] Service running
- [ ] Health check passes

### First App
- [ ] Test app deployed
- [ ] Pod running
- [ ] Ingress configured
- [ ] Certificate issued
- [ ] Accessible via HTTPS

### Production Ready
- [ ] Real app deployed
- [ ] Monitoring configured
- [ ] Logs accessible
- [ ] Backups configured (optional)
- [ ] Firewall configured

---

## Quick Commands Reference

```bash
# Check API status
curl http://localhost:8080/health

# List apps
curl http://localhost:8080/api/apps | jq

# View all deployments
kubectl get all -n superfly-apps

# View logs
kubectl logs -n superfly-apps deployment/APP-SLUG -f

# View certificates
kubectl get certificate -n superfly-apps

# Check API logs
sudo journalctl -u superfly-api -f

# Restart API
sudo systemctl restart superfly-api

# Database shell
psql postgresql://superfly:superfly_dev_password@localhost:5432/superfly
```

---

## Troubleshooting

### API not starting
```bash
sudo journalctl -u superfly-api -f
# Check for errors in logs
```

### DNS not resolving
```bash
dig superfly.smartynov.com
# Wait 5-30 minutes for propagation
```

### Certificate not issuing
```bash
kubectl describe certificate -n superfly-apps
kubectl logs -n cert-manager deployment/cert-manager -f
# Check DNS pointing to correct IP
# Check ports 80/443 accessible
```

### Pod not starting
```bash
kubectl describe pod -n superfly-apps POD-NAME
kubectl logs -n superfly-apps POD-NAME
# Common: image pull error, health check failing
```

---

## Success! 🎉

Your Superfly installation is complete when:

✅ API responds to health checks  
✅ Test app accessible via HTTPS  
✅ Certificate issued successfully  
✅ Logs visible in kubectl  
✅ Can deploy new apps via API  

**You can now deploy any Docker image with automatic HTTPS!**

---

## Next Steps

1. **Deploy your apps** - See [EXAMPLES.md](EXAMPLES.md)
2. **Setup monitoring** - Feature 4 & 5 (coming soon)
3. **Add environment variables** - Feature 3 (coming soon)
4. **GitHub integration** - Feature 2 (coming soon)
5. **Build Web UI** - Feature 8 (coming soon)

---

## Your Superfly Installation

- **Repository**: https://github.com/callmemars1/superfly
- **Main Domain**: superfly.smartynov.com
- **App Domains**: *.superfly.smartynov.com
- **API Endpoint**: http://YOUR_SERVER_IP:8080

**Have fun deploying!** 🚀
