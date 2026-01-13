#!/bin/bash

# Test script for Superfly API

API_URL="http://localhost:8080"
APP_ID=""

echo "🧪 Testing Superfly API"
echo "======================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

test_passed() {
    echo -e "${GREEN}✓${NC} $1"
}

test_failed() {
    echo -e "${RED}✗${NC} $1"
}

# Test 1: Health check
echo "Test 1: Health check"
response=$(curl -s "$API_URL/health")
if echo "$response" | grep -q "ok"; then
    test_passed "Health check passed"
else
    test_failed "Health check failed"
    echo "Response: $response"
fi
echo ""

# Test 2: List apps (should be empty initially)
echo "Test 2: List apps"
response=$(curl -s "$API_URL/api/apps")
if [ "$response" != "null" ]; then
    test_passed "List apps endpoint works"
    echo "Apps: $response"
else
    test_failed "List apps failed"
fi
echo ""

# Test 3: Create app
echo "Test 3: Create app (nginx)"
response=$(curl -s -X POST "$API_URL/api/apps" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Test Nginx",
        "image": "nginx:alpine",
        "port": 80,
        "domain": "test.local"
    }')

if echo "$response" | grep -q "id"; then
    test_passed "App created successfully"
    APP_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "App ID: $APP_ID"
    echo "Full response: $response"
else
    test_failed "App creation failed"
    echo "Response: $response"
    exit 1
fi
echo ""

# Wait a bit for deployment
echo "Waiting 10 seconds for deployment..."
sleep 10
echo ""

# Test 4: Get app by ID
echo "Test 4: Get app by ID"
response=$(curl -s "$API_URL/api/apps/$APP_ID")
if echo "$response" | grep -q "$APP_ID"; then
    test_passed "Get app by ID works"
    status=$(echo "$response" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
    echo "App status: $status"
else
    test_failed "Get app by ID failed"
fi
echo ""

# Test 5: Check Kubernetes resources
echo "Test 5: Check Kubernetes resources"
echo "Checking if deployment exists..."
if kubectl get deployment test-nginx -n superfly-apps &> /dev/null; then
    test_passed "Deployment exists in Kubernetes"
    kubectl get deployment test-nginx -n superfly-apps
else
    test_failed "Deployment not found in Kubernetes"
fi
echo ""

echo "Checking if service exists..."
if kubectl get service test-nginx -n superfly-apps &> /dev/null; then
    test_passed "Service exists in Kubernetes"
    kubectl get service test-nginx -n superfly-apps
else
    test_failed "Service not found in Kubernetes"
fi
echo ""

echo "Checking if ingress exists..."
if kubectl get ingress test-nginx -n superfly-apps &> /dev/null; then
    test_passed "Ingress exists in Kubernetes"
    kubectl get ingress test-nginx -n superfly-apps
else
    test_failed "Ingress not found in Kubernetes"
fi
echo ""

echo "Checking pod status..."
kubectl get pods -n superfly-apps -l app=test-nginx
echo ""

# Test 6: Update app
echo "Test 6: Update app (scale to 2 replicas)"
response=$(curl -s -X PATCH "$API_URL/api/apps/$APP_ID" \
    -H "Content-Type: application/json" \
    -d '{
        "replicas": 2
    }')

if echo "$response" | grep -q "replicas"; then
    test_passed "App updated successfully"
    replicas=$(echo "$response" | grep -o '"replicas":[0-9]*' | cut -d':' -f2)
    echo "New replicas: $replicas"
else
    test_failed "App update failed"
fi
echo ""

# Wait for scaling
echo "Waiting 10 seconds for scaling..."
sleep 10
kubectl get pods -n superfly-apps -l app=test-nginx
echo ""

# Test 7: Restart app
echo "Test 7: Restart app"
response=$(curl -s -X POST "$API_URL/api/apps/$APP_ID/restart")
if echo "$response" | grep -q "restart"; then
    test_passed "App restart initiated"
else
    test_failed "App restart failed"
fi
echo ""

# Test 8: List apps again
echo "Test 8: List apps (should show our app)"
response=$(curl -s "$API_URL/api/apps")
if echo "$response" | grep -q "test-nginx"; then
    test_passed "List apps shows our app"
else
    test_failed "List apps doesn't show our app"
fi
echo ""

# Test 9: Delete app
echo "Test 9: Delete app"
response=$(curl -s -X DELETE "$API_URL/api/apps/$APP_ID" -w "%{http_code}")
if [ "$response" = "204" ]; then
    test_passed "App deleted successfully"
else
    test_failed "App deletion failed (HTTP $response)"
fi
echo ""

# Wait for cleanup
echo "Waiting 5 seconds for cleanup..."
sleep 5
echo ""

# Test 10: Verify Kubernetes resources are deleted
echo "Test 10: Verify Kubernetes cleanup"
if ! kubectl get deployment test-nginx -n superfly-apps &> /dev/null; then
    test_passed "Deployment cleaned up"
else
    test_failed "Deployment still exists"
fi

if ! kubectl get service test-nginx -n superfly-apps &> /dev/null; then
    test_passed "Service cleaned up"
else
    test_failed "Service still exists"
fi

if ! kubectl get ingress test-nginx -n superfly-apps &> /dev/null; then
    test_passed "Ingress cleaned up"
else
    test_failed "Ingress still exists"
fi
echo ""

echo "======================="
echo "✅ All tests completed!"
echo "======================="
