#!/bin/bash

echo "=== Testing Beautified JSON Output ==="
echo ""

BASE_URL="http://localhost:3000"

echo "1. Testing Root Endpoint (should show beautified JSON)"
echo "curl $BASE_URL/"
echo ""
curl -s "$BASE_URL/"
echo ""
echo ""

echo "2. Testing GET /users (empty list with beautified JSON)"
echo "curl $BASE_URL/users"
echo ""
curl -s "$BASE_URL/users"
echo ""
echo ""

echo "3. Creating a user with beautified response"
echo "curl -X POST $BASE_URL/users -H 'Content-Type: application/json' -d '{...}'"
echo ""
curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001234",
    "first_name": "Somchai",
    "last_name": "Jaidee",
    "mobile_number": "081-234-5678",
    "email": "somchai@example.com",
    "register_date": "2023-06-15",
    "membership_level": "Gold",
    "point_balance": 15420
  }'
echo ""
echo ""

echo "4. Getting user by ID (beautified response)"
echo "curl $BASE_URL/users/1"
echo ""
curl -s "$BASE_URL/users/1"
echo ""
echo ""

echo "=== Beautified JSON Testing Complete ==="