#!/bin/bash

# Test script for Gateway Service

set -e

GATEWAY_URL="http://localhost:8000"
USER_SERVICE_URL="http://localhost:8080"

echo "🚀 Testing Gateway Service"

# Test health check
echo "Testing health check..."
curl -s "$GATEWAY_URL/health" | grep -q "ok" && echo "✅ Health check passed" || echo "❌ Health check failed"

# Test public endpoints (should work without auth)
echo "Testing public auth endpoints..."

# Test signup (should be proxied to user-service via internal ingress)
echo "Testing signup endpoint..."
SIGNUP_RESPONSE=$(curl -s -w "%{http_code}" -X POST "$GATEWAY_URL/api/v1/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "username": "testuser",
    "first_name": "Test",
    "last_name": "User"
  }')

echo "Signup response: $SIGNUP_RESPONSE"

# Test signin
echo "Testing signin endpoint..."
SIGNIN_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/v1/auth/signin" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com", 
    "password": "password123"
  }')

echo "Signin response: $SIGNIN_RESPONSE"

# Extract token from signin response (if successful)
if echo "$SIGNIN_RESPONSE" | grep -q "access_token"; then
  TOKEN=$(echo "$SIGNIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
  echo "✅ Got access token: ${TOKEN:0:20}..."
  
  # Test protected endpoint
  echo "Testing protected endpoint with token..."
  PROFILE_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$GATEWAY_URL/api/v1/profile" \
    -H "Authorization: Bearer $TOKEN")
  
  echo "Profile response: $PROFILE_RESPONSE"
  
  # Test endpoint without token (should fail)
  echo "Testing protected endpoint without token..."
  NO_AUTH_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$GATEWAY_URL/api/v1/profile")
  
  if echo "$NO_AUTH_RESPONSE" | grep -q "401"; then
    echo "✅ Correctly rejected request without token"
  else
    echo "❌ Should have rejected request without token"
  fi
  
else
  echo "❌ Could not get access token from signin"
fi

# Test non-existent route
echo "Testing non-existent route..."
NOT_FOUND_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$GATEWAY_URL/api/v1/nonexistent")

if echo "$NOT_FOUND_RESPONSE" | grep -q "404"; then
  echo "✅ Correctly returned 404 for non-existent route"
else
  echo "❌ Should have returned 404 for non-existent route"
fi

echo "🏁 Gateway tests completed"