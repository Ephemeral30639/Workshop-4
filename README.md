# LBK Membership User Management API

This project creates a RESTful API using Go and the Fiber framework for managing LBK membership users with SQLite database.

## Prerequisites

Make sure you have Go installed on your system. You can download it from: https://go.dev/dl/

## Installation

1. Install Go from https://go.dev/dl/ (minimum version 1.17)
2. Navigate to this project directory
3. Run the following commands:

```bash
go mod tidy
CGO_ENABLED=1 go run .
```

## Database Fields

The users table contains the following fields:
1. **First Name** (`first_name`) - User's first name
2. **Last Name** (`last_name`) - User's last name  
3. **Mobile Number** (`mobile_number`) - User's mobile phone number
4. **Email** (`email`) - User's email address (unique)
5. **Register Date** (`register_date`) - Date when user registered
6. **Membership Level** (`membership_level`) - Bronze/Silver/Gold/Platinum
7. **Point Balance** (`point_balance`) - Current points balance
8. **Member ID** (`member_id`) - Unique membership identifier

## API Endpoints

### Base URL: `http://localhost:3000`

- `GET /` - API information
- `GET /users` - List all users
- `GET /users/{id}` - Get user by ID
- `POST /users` - Create new user
- `PUT /users/{id}` - Update user by ID
- `DELETE /users/{id}` - Delete user by ID

## Example Usage

```bash
# Create a user
curl -X POST http://localhost:3000/users \
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

# List all users
curl http://localhost:3000/users

# Get user by ID
curl http://localhost:3000/users/1

# Update user
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"point_balance": 20000}'

# Delete user
curl -X DELETE http://localhost:3000/users/1
```

## Testing

Run the test script to verify all endpoints:
```bash
./test_api.sh
```