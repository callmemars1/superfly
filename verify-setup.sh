#!/bin/bash

# Superfly Development Environment Verification

echo "🔍 Verifying Superfly Development Environment..."
echo ""

FAILED=0

check_command() {
    if command -v $1 &> /dev/null; then
        echo "✓ $1 is installed"
        return 0
    else
        echo "✗ $1 is NOT installed"
        FAILED=1
        return 1
    fi
}

check_k8s_resource() {
    if kubectl get $1 $2 -n $3 &> /dev/null; then
        echo "✓ $1/$2 exists in namespace $3"
        return 0
    else
        echo "✗ $1/$2 NOT found in namespace $3"
        FAILED=1
        return 1
    fi
}

# Check commands
echo "Checking installed tools..."
check_command go
check_command node
check_command npm
check_command kubectl
check_command helm
check_command docker
check_command sqlc
check_command goose
check_command air
check_command psql

echo ""
echo "Checking K3S cluster..."
if kubectl cluster-info &> /dev/null; then
    echo "✓ K3S cluster is running"
else
    echo "✗ K3S cluster is NOT accessible"
    FAILED=1
fi

echo ""
echo "Checking namespaces..."
check_k8s_resource namespace superfly-system ""
check_k8s_resource namespace superfly-apps ""
check_k8s_resource namespace cert-manager ""
check_k8s_resource namespace traefik ""

echo ""
echo "Checking deployments..."
check_k8s_resource deployment postgres superfly-system
check_k8s_resource deployment registry superfly-system
check_k8s_resource deployment cert-manager cert-manager
check_k8s_resource deployment traefik traefik

echo ""
echo "Checking PostgreSQL connection..."
if PGPASSWORD=superfly_dev_password psql -h localhost -U superfly -d superfly -c "SELECT 1" &> /dev/null; then
    echo "✓ PostgreSQL is accessible"
else
    echo "✗ PostgreSQL connection failed"
    echo "  (Port forward might not be running)"
    FAILED=1
fi

echo ""
echo "Checking ClusterIssuers..."
check_k8s_resource clusterissuer letsencrypt-staging ""
check_k8s_resource clusterissuer letsencrypt-prod ""

echo ""
if [ $FAILED -eq 0 ]; then
    echo "✅ All checks passed! Environment is ready."
    exit 0
else
    echo "❌ Some checks failed. Please review the output above."
    exit 1
fi
