# TL;DR - Quick Deploy

Get an app running in 10 minutes. No fluff.

## On Your Server

```bash
# 1. Clone and setup
git clone https://github.com/callmemars1/superfly.git
cd superfly
chmod +x *.sh

# 2. Install everything (takes 5 mins)
./dev-setup.sh

# 3. Verify
./verify-setup.sh

# 4. Build
make init
go mod tidy
make migrate
make sqlc-generate
make build

# 5. Start API (background)
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
Environment="DATABASE_URL=postgresql://superfly:superfly_dev_password@localhost:5432/superfly?sslmode=disable"
Environment="KUBERNETES_IN_CLUSTER=false"
Environment="KUBECONFIG=$HOME/.kube/config"

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now superfly-api

# 6. Deploy nginx
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "image": "nginx:alpine",
    "port": 80,
    "domain": "test.superfly.smartynov.com"
  }'

# 7. Watch it deploy
kubectl get pods -n superfly-apps -w

# 8. Access it
curl https://test.superfly.smartynov.com
```

## Done! 🎉

Your app is live with HTTPS.

## Common Commands

```bash
# List apps
curl localhost:8080/api/apps | jq

# Deploy app
curl -X POST localhost:8080/api/apps -d '{"name":"App","image":"nginx:alpine","port":80}'

# Scale app
curl -X PATCH localhost:8080/api/apps/ID -d '{"replicas":3}'

# Delete app
curl -X DELETE localhost:8080/api/apps/ID

# View logs
kubectl logs -n superfly-apps deployment/SLUG -f

# Check status
kubectl get all -n superfly-apps
```

## Troubleshooting

```bash
# API not working?
sudo systemctl restart superfly-api
sudo journalctl -u superfly-api -f

# Pod not starting?
kubectl describe pod -n superfly-apps POD-NAME
kubectl logs -n superfly-apps POD-NAME

# Database issue?
sudo systemctl restart superfly-postgres-forward
psql postgresql://superfly:superfly_dev_password@localhost:5432/superfly
```

## For Details

See [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) for the full guide.
