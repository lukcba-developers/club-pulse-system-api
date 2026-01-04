#!/bin/bash
echo "Logging in..."
LOGIN_RESP=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@clubpulse.com", "password": "admin123"}')

echo "Login Response: $LOGIN_RESP"

TOKEN=$(echo $LOGIN_RESP | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "Failed to get token"
  exit 1
fi

echo "Token: $TOKEN"

echo "Fetching Me..."
curl -v -X GET http://localhost:8081/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Club-ID: club-alpha"
