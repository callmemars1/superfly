#!/bin/bash
set -e

# Superfly Development Environment Setup Script
# For Debian 13 (Trixie)

SUPERFLY_DIR="$HOME/superfly"
K3S_VERSION="v1.28.5+k3s1"
GO_VERSION="1.22.0"
NODE_VERSION="20"

echo "=================================="
echo "🚀 Superfly Dev Environment Setup"
echo "=================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}✓${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1"
}

# Determine if we need sudo
SUDO=""
if [ "$EUID" -ne 0 ]; then
    SUDO="sudo"
    log_info "Running as regular user (will use sudo)"
else
    log_info "Running as root"
fi

# Update system
log_info "Updating system packages..."
$SUDO apt update
$SUDO apt upgrade -y

# Install basic dependencies
log_info "Installing basic dependencies..."
$SUDO apt install -y \
    curl \
    wget \
    git \
    build-essential \
    ca-certificates \
    gnupg \
    lsb-release \
    jq \
    unzip

# Install Go
log_info "Installing Go ${GO_VERSION}..."
if ! command -v go &> /dev/null; then
    cd /tmp
    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    $SUDO rm -rf /usr/local/go
    $SUDO tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    
    # Add to PATH if not already there
    SHELL_RC="$HOME/.bashrc"
    if [ "$EUID" -eq 0 ]; then
        SHELL_RC="/root/.bashrc"
    fi
    if ! grep -q "/usr/local/go/bin" "$SHELL_RC"; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> "$SHELL_RC"
        echo 'export PATH=$PATH:$HOME/go/bin' >> "$SHELL_RC"
    fi
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
    log_info "Go installed: $(go version)"
else
    log_info "Go already installed: $(go version)"
fi

# Install Node.js (for Svelte frontend later)
log_info "Installing Node.js ${NODE_VERSION}..."
if ! command -v node &> /dev/null; then
    if [ "$EUID" -eq 0 ]; then
        curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash -
    else
        curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | sudo -E bash -
    fi
    $SUDO apt install -y nodejs
    log_info "Node.js installed: $(node --version)"
    log_info "npm installed: $(npm --version)"
else
    log_info "Node.js already installed: $(node --version)"
fi

# Install Docker (for local testing, not required but useful)
log_info "Installing Docker..."
if ! command -v docker &> /dev/null; then
    $SUDO apt install -y docker.io
    $SUDO systemctl enable docker
    $SUDO systemctl start docker
    if [ "$EUID" -ne 0 ]; then
        $SUDO usermod -aG docker $USER
        log_warn "Docker installed. You need to log out and back in for group changes to take effect"
    fi
else
    log_info "Docker already installed"
fi

# Install K3S
log_info "Installing K3S..."
if ! command -v k3s &> /dev/null; then
    curl -sfL https://get.k3s.io | sh -s - \
        --write-kubeconfig-mode 644 \
        --disable traefik
    
    # Wait for K3S to be ready
    log_info "Waiting for K3S to be ready..."
    sleep 10
    k3s kubectl wait --for=condition=Ready nodes --all --timeout=60s
    
    # Setup kubeconfig for current user
    mkdir -p ~/.kube
    if [ "$EUID" -eq 0 ]; then
        cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
    else
        $SUDO cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
        $SUDO chown $USER:$USER ~/.kube/config
    fi
    
    log_info "K3S installed successfully"
else
    log_info "K3S already installed"
fi

# Install kubectl (if not using k3s kubectl)
log_info "Setting up kubectl..."
if ! command -v kubectl &> /dev/null; then
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    $SUDO install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    rm kubectl
    log_info "kubectl installed: $(kubectl version --client --short 2>/dev/null || kubectl version --client)"
else
    log_info "kubectl already installed"
fi

# Install Helm
log_info "Installing Helm..."
if ! command -v helm &> /dev/null; then
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    log_info "Helm installed: $(helm version --short)"
else
    log_info "Helm already installed"
fi

# Install cert-manager
log_info "Installing cert-manager..."
kubectl create namespace cert-manager --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.0/cert-manager.yaml

log_info "Waiting for cert-manager to be ready..."
kubectl wait --for=condition=Available --timeout=300s \
    deployment/cert-manager \
    deployment/cert-manager-webhook \
    deployment/cert-manager-cainjector \
    -n cert-manager

# Create Let's Encrypt ClusterIssuers
log_info "Creating Let's Encrypt ClusterIssuers..."
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    email: admin@superfly.local
    privateKeySecretRef:
      name: letsencrypt-staging
    solvers:
    - http01:
        ingress:
          class: traefik
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@superfly.local
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: traefik
EOF

# Install Traefik (since we disabled it during K3S install)
log_info "Installing Traefik..."
helm repo add traefik https://traefik.github.io/charts
helm repo update

kubectl create namespace traefik --dry-run=client -o yaml | kubectl apply -f -

helm upgrade --install traefik traefik/traefik \
    --namespace traefik \
    --set ports.web.exposedPort=80 \
    --set ports.websecure.exposedPort=443 \
    --set ports.web.nodePort=30080 \
    --set ports.websecure.nodePort=30443 \
    --set service.type=LoadBalancer \
    --wait

log_info "Waiting for Traefik to be ready..."
kubectl wait --for=condition=Available --timeout=120s deployment/traefik -n traefik

# Create superfly-system and superfly-apps namespaces
log_info "Creating Superfly namespaces..."
kubectl create namespace superfly-system --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace superfly-apps --dry-run=client -o yaml | kubectl apply -f -

# Deploy PostgreSQL in K3S for development
log_info "Deploying PostgreSQL..."
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: postgres-secret
  namespace: superfly-system
type: Opaque
stringData:
  POSTGRES_PASSWORD: superfly_dev_password
  POSTGRES_USER: superfly
  POSTGRES_DB: superfly
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: superfly-system
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: superfly-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:16-alpine
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: POSTGRES_PASSWORD
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: POSTGRES_USER
        - name: POSTGRES_DB
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: POSTGRES_DB
        - name: PGDATA
          value: /var/lib/postgresql/data/pgdata
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: superfly-system
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
  type: ClusterIP
EOF

log_info "Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=Available --timeout=120s deployment/postgres -n superfly-system

# Deploy Container Registry
log_info "Deploying Container Registry..."
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: registry-pvc
  namespace: superfly-system
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: registry
  namespace: superfly-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: registry
  template:
    metadata:
      labels:
        app: registry
    spec:
      containers:
      - name: registry
        image: registry:2
        ports:
        - containerPort: 5000
        volumeMounts:
        - name: registry-storage
          mountPath: /var/lib/registry
        env:
        - name: REGISTRY_STORAGE_DELETE_ENABLED
          value: "true"
      volumes:
      - name: registry-storage
        persistentVolumeClaim:
          claimName: registry-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: registry
  namespace: superfly-system
spec:
  selector:
    app: registry
  ports:
  - port: 5000
    targetPort: 5000
  type: ClusterIP
EOF

log_info "Waiting for Registry to be ready..."
kubectl wait --for=condition=Available --timeout=120s deployment/registry -n superfly-system

# Install Go tools
log_info "Installing Go development tools..."
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/air-verse/air@latest  # For live reload during dev

# Install PostgreSQL client tools (for local access)
log_info "Installing PostgreSQL client..."
$SUDO apt install -y postgresql-client

# Port forward PostgreSQL for local development (as a background service)
log_info "Setting up PostgreSQL port forward..."
SERVICE_USER="$USER"
if [ "$EUID" -eq 0 ]; then
    SERVICE_USER="root"
fi

cat <<EOF | $SUDO tee /etc/systemd/system/superfly-postgres-forward.service > /dev/null
[Unit]
Description=Superfly PostgreSQL Port Forward
After=network.target

[Service]
Type=simple
User=$SERVICE_USER
ExecStart=/usr/local/bin/kubectl port-forward -n superfly-system svc/postgres 5432:5432
Restart=always
RestartSec=10
Environment="KUBECONFIG=$HOME/.kube/config"

[Install]
WantedBy=multi-user.target
EOF

$SUDO systemctl daemon-reload
$SUDO systemctl enable superfly-postgres-forward.service
$SUDO systemctl start superfly-postgres-forward.service

# Create development .env file template
log_info "Creating development environment template..."
mkdir -p $SUPERFLY_DIR
cat <<EOF > $SUPERFLY_DIR/.env.example
# Database
DATABASE_URL=postgresql://superfly:superfly_dev_password@localhost:5432/superfly?sslmode=disable

# Kubernetes
KUBERNETES_IN_CLUSTER=false
KUBECONFIG=$HOME/.kube/config

# Registry
REGISTRY_URL=registry.superfly-system.svc.cluster.local:5000

# API Server
API_PORT=8080
API_HOST=0.0.0.0

# Environment
ENV=development
LOG_LEVEL=debug
EOF

cat <<EOF > $SUPERFLY_DIR/.env
# Database
DATABASE_URL=postgresql://superfly:superfly_dev_password@localhost:5432/superfly?sslmode=disable

# Kubernetes
KUBERNETES_IN_CLUSTER=false
KUBECONFIG=$HOME/.kube/config

# Registry
REGISTRY_URL=registry.superfly-system.svc.cluster.local:5000

# API Server
API_PORT=8080
API_HOST=0.0.0.0

# Environment
ENV=development
LOG_LEVEL=debug
EOF

# Print summary
echo ""
echo "=================================="
echo "✅ Setup Complete!"
echo "=================================="
echo ""
log_info "K3S cluster running"
log_info "PostgreSQL: localhost:5432"
log_info "  - Database: superfly"
log_info "  - User: superfly"
log_info "  - Password: superfly_dev_password"
log_info "Container Registry: registry.superfly-system.svc.cluster.local:5000"
echo ""
echo "📝 Database connection string:"
echo "   postgresql://superfly:superfly_dev_password@localhost:5432/superfly"
echo ""
echo "🔧 Development tools installed:"
echo "   - Go $(go version | awk '{print $3}')"
echo "   - Node.js $(node --version)"
echo "   - kubectl $(kubectl version --client --short 2>/dev/null || echo '')"
echo "   - helm $(helm version --short)"
echo "   - sqlc"
echo "   - goose (migrations)"
echo "   - air (live reload)"
echo ""
echo "📁 Development directory: $SUPERFLY_DIR"
echo "   Environment file: $SUPERFLY_DIR/.env"
echo ""
echo "🚀 Next steps:"
echo "   1. cd $SUPERFLY_DIR"
echo "   2. Initialize your Go project"
echo "   3. Start coding!"
echo ""
log_warn "If you installed Docker, please log out and back in for group permissions"
echo ""
echo "🔍 Useful commands:"
echo "   kubectl get pods -n superfly-system    # Check system pods"
echo "   kubectl get pods -n superfly-apps      # Check app pods"
echo "   kubectl logs -n superfly-system postgres-xxx  # View postgres logs"
echo "   psql postgresql://superfly:superfly_dev_password@localhost:5432/superfly  # Connect to DB"
echo ""
