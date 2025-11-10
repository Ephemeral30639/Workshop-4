#!/bin/bash

echo "=== Adding 10 More Users to LBK Membership Database ==="
echo ""

BASE_URL="http://localhost:3000"

echo "User 1: Somchai Jaidee (Gold Member)"
curl -X POST "$BASE_URL/users" \
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
  }'
echo ""

echo "User 2: Siriporn Thanakit (Platinum Member)"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001235",
    "first_name": "สิริพร",
    "last_name": "ธนกิจ",
    "mobile_number": "082-567-8901",
    "email": "siriporn.t@example.com",
    "register_date": "2022-03-10",
    "membership_level": "Platinum",
    "point_balance": 45780
  }'
echo ""

echo "User 3: John Smith (Silver Member)"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001236",
    "first_name": "John",
    "last_name": "Smith",
    "mobile_number": "089-123-4567",
    "email": "john.smith@example.com",
    "register_date": "2024-01-20",
    "membership_level": "Silver",
    "point_balance": 8950
  }'
echo ""

echo "User 4: Apinya Wongsuwan (Bronze Member)"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001237",
    "first_name": "อภิญญา",
    "last_name": "วงศ์สุวรรณ",
    "mobile_number": "095-876-5432",
    "email": "apinya.w@example.com",
    "register_date": "2024-08-05",
    "membership_level": "Bronze",
    "point_balance": 2340
  }'
echo ""

echo "User 5: Maria Garcia (Gold Member)"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001238",
    "first_name": "Maria",
    "last_name": "Garcia",
    "mobile_number": "061-345-6789",
    "email": "maria.garcia@example.com",
    "register_date": "2023-11-12",
    "membership_level": "Gold",
    "point_balance": 18750
  }'
echo ""

echo "User 6: Pongsakorn Rattanakit (Silver Member)"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001239",
    "first_name": "พงศกร",
    "last_name": "รัตนกิจ",
    "mobile_number": "098-234-5671",
    "email": "pongsakorn.r@example.com",
    "register_date": "2024-02-28",
    "membership_level": "Silver",
    "point_balance": 12680
  }'
echo ""

echo "User 7: Emily Johnson (Platinum Member)"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001240",
    "first_name": "Emily",
    "last_name": "Johnson",
    "mobile_number": "065-789-0123",
    "email": "emily.johnson@example.com",
    "register_date": "2021-09-15",
    "membership_level": "Platinum",
    "point_balance": 52340
  }'
echo ""

echo "User 8: Nattaporn Srisawat (Bronze Member)"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001241",
    "first_name": "ณัฐพร",
    "last_name": "ศรีสวัสดิ์",
    "mobile_number": "092-456-7890",
    "email": "nattaporn.s@example.com",
    "register_date": "2024-07-18",
    "membership_level": "Bronze",
    "point_balance": 4520
  }'
echo ""

echo "User 9: David Chen (Gold Member)"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001242",
    "first_name": "David",
    "last_name": "Chen",
    "mobile_number": "084-678-9012",
    "email": "david.chen@example.com",
    "register_date": "2023-04-07",
    "membership_level": "Gold",
    "point_balance": 23150
  }'
echo ""

echo "User 10: Kanlaya Promsuwan (Silver Member)"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001243",
    "first_name": "กัลยา",
    "last_name": "พรมสุวรรณ",
    "mobile_number": "087-890-1234",
    "email": "kanlaya.p@example.com",
    "register_date": "2023-12-03",
    "membership_level": "Silver",
    "point_balance": 9870
  }'
echo ""

echo "=== All 10 Users Added Successfully! ==="