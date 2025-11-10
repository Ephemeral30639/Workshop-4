#!/bin/bash

echo "=== LBK Membership User Management API Testing ==="
echo ""

BASE_URL="http://localhost:3000"

echo "1. Testing Root Endpoint"
curl -s "$BASE_URL/" | python3 -m json.tool
echo ""

echo "2. GET /users - List all users (should be empty initially)"
curl -s "$BASE_URL/users" | python3 -m json.tool
echo ""

echo "3. POST /users - Create the user from the Thai screenshot (Somchai Jaidee)"
curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001234",
    "first_name": "สมชาย",
    "last_name": "ใจดี",
    "mobile_number": "081-234-5678",
    "email": "somchai@example.com",
    "register_date": "2023-06-15",
    "membership_level": "Gold",
    "point_balance": 15420
  }' | python3 -m json.tool
echo ""

echo "4. POST /users - Create another test user"
curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001235",
    "first_name": "วิลลี่",
    "last_name": "สมิท",
    "mobile_number": "081-999-8888",
    "email": "willy@example.com",
    "register_date": "2024-01-10",
    "membership_level": "Silver",
    "point_balance": 5000
  }' | python3 -m json.tool
echo ""

echo "5. GET /users - List all users (should now have 2 users)"
curl -s "$BASE_URL/users" | python3 -m json.tool
echo ""

echo "6. GET /users/1 - Get user by ID"
curl -s "$BASE_URL/users/1" | python3 -m json.tool
echo ""

echo "7. PUT /users/1 - Update user points"
curl -s -X PUT "$BASE_URL/users/1" \
  -H "Content-Type: application/json" \
  -d '{
    "point_balance": 20000,
    "membership_level": "Platinum"
  }' | python3 -m json.tool
echo ""

echo "8. GET /users/1 - Verify the update"
curl -s "$BASE_URL/users/1" | python3 -m json.tool
echo ""

echo "9. DELETE /users/2 - Delete second user"
curl -s -X DELETE "$BASE_URL/users/2" | python3 -m json.tool
echo ""

echo "10. GET /users - Final list (should have only 1 user)"
curl -s "$BASE_URL/users" | python3 -m json.tool
echo ""

echo "=== API Testing Complete ==="