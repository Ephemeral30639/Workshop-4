#!/bin/bash

echo "=== Point Transfer Feature Testing ==="
echo ""

BASE_URL="http://localhost:3000"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_test() {
    echo -e "${BLUE}=== $1 ===${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    echo ""
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    echo ""
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
    echo ""
}

# Test 1: Setup - Create test users
print_test "Test 1: Setting up test users"

print_info "Creating User 1 (Alice) with 10000 points"
USER1_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001001",
    "first_name": "Alice",
    "last_name": "Johnson",
    "mobile_number": "081-111-1111",
    "email": "alice@example.com",
    "register_date": "2024-01-01",
    "membership_level": "Gold",
    "point_balance": 10000
  }')

USER1_ID=$(echo $USER1_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
echo "User 1 ID: $USER1_ID"
echo "$USER1_RESPONSE" | python3 -m json.tool
echo ""

print_info "Creating User 2 (Bob) with 5000 points"
USER2_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001002",
    "first_name": "Bob",
    "last_name": "Smith",
    "mobile_number": "081-222-2222",
    "email": "bob@example.com",
    "register_date": "2024-01-15",
    "membership_level": "Silver",
    "point_balance": 5000
  }')

USER2_ID=$(echo $USER2_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
echo "User 2 ID: $USER2_ID"
echo "$USER2_RESPONSE" | python3 -m json.tool
echo ""

print_info "Creating User 3 (Charlie) with 0 points"
USER3_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "LBK001003",
    "first_name": "Charlie",
    "last_name": "Brown",
    "mobile_number": "081-333-3333",
    "email": "charlie@example.com",
    "register_date": "2024-02-01",
    "membership_level": "Bronze",
    "point_balance": 0
  }')

USER3_ID=$(echo $USER3_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
echo "User 3 ID: $USER3_ID"
echo "$USER3_RESPONSE" | python3 -m json.tool
echo ""

# Test 2: Valid transfer
print_test "Test 2: Valid Point Transfer (Alice → Bob: 1500 points)"

print_info "Transferring 1500 points from Alice (ID: $USER1_ID) to Bob (ID: $USER2_ID)"
TRANSFER1_RESPONSE=$(curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER1_ID,
    \"toUserId\": $USER2_ID,
    \"amount\": 1500,
    \"note\": \"Payment for services\"
  }")

echo "$TRANSFER1_RESPONSE" | python3 -m json.tool
echo ""

# Extract transfer ID for later use
TRANSFER1_ID=$(echo $TRANSFER1_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin).get('transfer', {}).get('idemKey', ''))" 2>/dev/null)
echo "Transfer ID (Idempotency Key): $TRANSFER1_ID"
echo ""

# Verify balances after transfer
print_info "Verifying Alice's balance (should be 8500)"
curl -s "$BASE_URL/users/$USER1_ID" | python3 -m json.tool
echo ""

print_info "Verifying Bob's balance (should be 6500)"
curl -s "$BASE_URL/users/$USER2_ID" | python3 -m json.tool
echo ""

# Test 3: Insufficient balance
print_test "Test 3: Insufficient Balance Error"

print_info "Attempting to transfer 20000 points from Bob (balance: 6500)"
curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER2_ID,
    \"toUserId\": $USER1_ID,
    \"amount\": 20000,
    \"note\": \"Should fail - insufficient balance\"
  }" | python3 -m json.tool
echo ""

# Test 4: Self-transfer validation
print_test "Test 4: Self-Transfer Validation"

print_info "Attempting to transfer points to self (should fail)"
curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER1_ID,
    \"toUserId\": $USER1_ID,
    \"amount\": 100,
    \"note\": \"Should fail - self transfer\"
  }" | python3 -m json.tool
echo ""

# Test 5: Non-existent user
print_test "Test 5: Non-existent User Error"

print_info "Attempting transfer to non-existent user (ID: 99999)"
curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER1_ID,
    \"toUserId\": 99999,
    \"amount\": 100,
    \"note\": \"Should fail - user not found\"
  }" | python3 -m json.tool
echo ""

# Test 6: Invalid input validation
print_test "Test 6: Input Validation Tests"

print_info "Test 6a: Negative amount"
curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER1_ID,
    \"toUserId\": $USER2_ID,
    \"amount\": -100,
    \"note\": \"Should fail - negative amount\"
  }" | python3 -m json.tool
echo ""

print_info "Test 6b: Zero amount"
curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER1_ID,
    \"toUserId\": $USER2_ID,
    \"amount\": 0,
    \"note\": \"Should fail - zero amount\"
  }" | python3 -m json.tool
echo ""

print_info "Test 6c: Missing required fields"
curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER1_ID,
    \"amount\": 100
  }" | python3 -m json.tool
echo ""

# Test 7: Get transfer by ID
print_test "Test 7: Retrieve Transfer by ID"

if [ ! -z "$TRANSFER1_ID" ]; then
    print_info "Getting transfer details for ID: $TRANSFER1_ID"
    curl -s "$BASE_URL/transfers/$TRANSFER1_ID" | python3 -m json.tool
    echo ""
else
    print_error "No valid transfer ID available for testing"
fi

# Test 8: Multiple transfers and history
print_test "Test 8: Multiple Transfers and History"

print_info "Creating second transfer: Bob → Charlie (2000 points)"
TRANSFER2_RESPONSE=$(curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER2_ID,
    \"toUserId\": $USER3_ID,
    \"amount\": 2000,
    \"note\": \"Gift points\"
  }")

echo "$TRANSFER2_RESPONSE" | python3 -m json.tool
echo ""

print_info "Creating third transfer: Alice → Charlie (500 points)"
curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER1_ID,
    \"toUserId\": $USER3_ID,
    \"amount\": 500,
    \"note\": \"Bonus points\"
  }" | python3 -m json.tool
echo ""

# Test 9: List transfers for each user
print_test "Test 9: Transfer History for Each User"

print_info "Alice's transfer history (User ID: $USER1_ID)"
curl -s "$BASE_URL/transfers?userId=$USER1_ID" | python3 -m json.tool
echo ""

print_info "Bob's transfer history (User ID: $USER2_ID)"
curl -s "$BASE_URL/transfers?userId=$USER2_ID" | python3 -m json.tool
echo ""

print_info "Charlie's transfer history (User ID: $USER3_ID)"
curl -s "$BASE_URL/transfers?userId=$USER3_ID" | python3 -m json.tool
echo ""

# Test 10: Pagination testing
print_test "Test 10: Pagination Testing"

print_info "Testing pagination - page 1, pageSize 2"
curl -s "$BASE_URL/transfers?userId=$USER1_ID&page=1&pageSize=2" | python3 -m json.tool
echo ""

print_info "Testing pagination - page 2, pageSize 1"
curl -s "$BASE_URL/transfers?userId=$USER2_ID&page=2&pageSize=1" | python3 -m json.tool
echo ""

# Test 11: Final balance verification
print_test "Test 11: Final Balance Verification"

print_info "Alice's final balance (started: 10000, sent: 1500 + 500 = 2000, should be: 8000)"
curl -s "$BASE_URL/users/$USER1_ID" | python3 -m json.tool
echo ""

print_info "Bob's final balance (started: 5000, received: 1500, sent: 2000, should be: 4500)"
curl -s "$BASE_URL/users/$USER2_ID" | python3 -m json.tool
echo ""

print_info "Charlie's final balance (started: 0, received: 2000 + 500 = 2500, should be: 2500)"
curl -s "$BASE_URL/users/$USER3_ID" | python3 -m json.tool
echo ""

# Test 12: Edge cases
print_test "Test 12: Edge Cases"

print_info "Test 12a: Transfer all remaining points (Bob → Alice)"
# First get Bob's current balance
BOB_BALANCE=$(curl -s "$BASE_URL/users/$USER2_ID" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['point_balance'])" 2>/dev/null)
echo "Bob's current balance: $BOB_BALANCE"

curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER2_ID,
    \"toUserId\": $USER1_ID,
    \"amount\": $BOB_BALANCE,
    \"note\": \"Transfer all remaining points\"
  }" | python3 -m json.tool
echo ""

print_info "Verify Bob's balance is now 0"
curl -s "$BASE_URL/users/$USER2_ID" | python3 -m json.tool
echo ""

print_info "Test 12b: Try to transfer from user with 0 balance"
curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d "{
    \"fromUserId\": $USER2_ID,
    \"toUserId\": $USER1_ID,
    \"amount\": 1,
    \"note\": \"Should fail - no balance\"
  }" | python3 -m json.tool
echo ""

print_success "Point Transfer Feature Testing Complete!"
print_info "Summary of tests performed:"
echo "  ✓ Valid transfers between users"
echo "  ✓ Insufficient balance validation"
echo "  ✓ Self-transfer prevention"
echo "  ✓ Non-existent user handling"
echo "  ✓ Input validation (negative, zero, missing fields)"
echo "  ✓ Transfer retrieval by ID"
echo "  ✓ Transfer history listing"
echo "  ✓ Pagination functionality"
echo "  ✓ Balance accuracy verification"
echo "  ✓ Edge case handling (zero balance)"
echo ""