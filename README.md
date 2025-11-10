# LBK Membership User Management & Point Transfer API

This project creates a comprehensive RESTful API using Go and the Fiber framework for managing LBK membership users with point transfer functionality, built on SQLite database with transaction management and audit logging.

## Features

- 🏪 **User Management**: Complete CRUD operations for membership users
- 💰 **Point Transfer System**: Secure point transfers between users with transaction integrity  
- 📊 **Transaction Ledger**: Comprehensive audit trail for all point movements
- 🔒 **Data Validation**: Input validation and business rule enforcement
- 📄 **Pagination**: Efficient data retrieval with pagination support
- 🔄 **Idempotency**: Transfer operations with idempotency keys for reliability
- 🏗️ **Database Schema**: Normalized schema with proper foreign key relationships

## Prerequisites

Make sure you have Go installed on your system. You can download it from: https://go.dev/dl/

## Installation

1. Install Go from https://go.dev/dl/ (minimum version 1.21)
2. Navigate to this project directory
3. Run the following commands:

```bash
go mod tidy
CGO_ENABLED=1 go run .
```

The server will start on port 3000.

## Database Schema

### Users Table
- **ID** (`id`) - Auto-increment primary key
- **Member ID** (`member_id`) - Unique membership identifier
- **First Name** (`first_name`) - User's first name
- **Last Name** (`last_name`) - User's last name  
- **Mobile Number** (`mobile_number`) - User's mobile phone number
- **Email** (`email`) - User's email address (unique)
- **Register Date** (`register_date`) - Date when user registered
- **Membership Level** (`membership_level`) - Bronze/Silver/Gold/Platinum
- **Point Balance** (`point_balance`) - Current points balance
- **Created At** (`created_at`) - Record creation timestamp
- **Updated At** (`updated_at`) - Record last update timestamp

### Transfers Table
- **ID** (`id`) - Auto-increment primary key
- **From User ID** (`from_user_id`) - Source user for transfer
- **To User ID** (`to_user_id`) - Destination user for transfer
- **Amount** (`amount`) - Transfer amount (must be positive)
- **Status** (`status`) - Transfer status (pending/processing/completed/failed/cancelled/reversed)
- **Note** (`note`) - Optional transfer description
- **Idempotency Key** (`idempotency_key`) - Unique key for idempotent operations
- **Created At** (`created_at`) - Transfer creation timestamp
- **Updated At** (`updated_at`) - Transfer last update timestamp
- **Completed At** (`completed_at`) - Transfer completion timestamp
- **Fail Reason** (`fail_reason`) - Failure reason if transfer failed

### Point Ledger Table
- **ID** (`id`) - Auto-increment primary key
- **User ID** (`user_id`) - User associated with the entry
- **Change** (`change`) - Point change amount (positive/negative)
- **Balance After** (`balance_after`) - User's balance after this transaction
- **Event Type** (`event_type`) - Type of event (transfer_out/transfer_in/adjust/earn/redeem)
- **Transfer ID** (`transfer_id`) - Associated transfer ID (if applicable)
- **Reference** (`reference`) - Additional reference information
- **Metadata** (`metadata`) - JSON metadata for the transaction
- **Created At** (`created_at`) - Entry creation timestamp

## API Endpoints

### Base URL: `http://localhost:3000`

#### User Management
- `GET /` - API information and version
- `GET /users` - List all users with count
- `GET /users/{id}` - Get user by ID
- `POST /users` - Create new user
- `PUT /users/{id}` - Update user by ID (partial updates supported)
- `DELETE /users/{id}` - Delete user by ID

#### Point Transfer System
- `POST /transfers` - Create a new point transfer
- `GET /transfers/{idempotencyKey}` - Get transfer by idempotency key
- `GET /transfers?userId={id}&page={page}&pageSize={size}` - List transfers for a user (paginated)

## Example Usage

### User Management

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

# Update user (partial update)
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"point_balance": 20000, "membership_level": "Platinum"}'

# Delete user
curl -X DELETE http://localhost:3000/users/1
```

### Point Transfer Operations

```bash
# Transfer points between users
curl -X POST http://localhost:3000/transfers \
  -H "Content-Type: application/json" \
  -d '{
    "fromUserId": 1,
    "toUserId": 2,
    "amount": 1500,
    "note": "Payment for services"
  }'

# Get transfer by idempotency key
curl http://localhost:3000/transfers/{idempotency-key}

# List transfers for a user (paginated)
curl "http://localhost:3000/transfers?userId=1&page=1&pageSize=10"
```

### Error Responses

The API returns structured error responses:

```json
{
  "error": "VALIDATION_ERROR",
  "message": "Invalid request body"
}
```

Common error codes:
- `VALIDATION_ERROR` - Input validation failed
- `BUSINESS_ERROR` - Business rule violation (e.g., self-transfer)
- `NOT_FOUND` - Resource not found
- `INSUFFICIENT_BALANCE` - Not enough points for transfer
- `INTERNAL_ERROR` - Server-side error

## Testing

### Basic API Testing
Run the basic test script to verify user management endpoints:
```bash
chmod +x test_api.sh
./test_api.sh
```

### Point Transfer Testing
Run comprehensive point transfer tests:
```bash
chmod +x test_transfer_feature.sh
./test_transfer_feature.sh
```

### Add Sample Data
Add 10 sample users with various membership levels:
```bash
chmod +x add_10_users.sh
./add_10_users.sh
```

### Beautified Test Output
Run tests with formatted output:
```bash
chmod +x test_beautified.sh
./test_beautified.sh
```

## Business Rules

### Point Transfers
1. **Positive Amounts**: Transfer amounts must be greater than zero
2. **Sufficient Balance**: Users must have enough points to transfer
3. **No Self-Transfers**: Users cannot transfer points to themselves
4. **User Validation**: Both sender and receiver must exist
5. **Atomicity**: All transfer operations are atomic (all-or-nothing)
6. **Audit Trail**: Every point movement is logged in the ledger

### Data Validation
1. **Required Fields**: First name, last name, and member ID are required for users
2. **Unique Constraints**: Member ID and email must be unique
3. **Membership Levels**: Bronze, Silver, Gold, Platinum
4. **Point Balance**: Cannot be negative

## Project Structure

```
├── main.go              # Application entry point and database setup
├── handlers.go          # HTTP request handlers for all endpoints
├── go.mod              # Go module dependencies
├── README.md           # This documentation
├── test_api.sh         # Basic API testing script
├── test_transfer_feature.sh # Comprehensive point transfer testing
├── test_beautified.sh  # Formatted test output script
└── add_10_users.sh     # Sample data creation script
```

## Dependencies

- **Fiber v2**: Fast HTTP web framework
- **SQLite3**: Embedded SQL database
- **Google UUID**: UUID generation for idempotency keys

## Version History

- **v1.0.0**: Initial user management API
- **v2.0.0**: Added point transfer system with transaction management
- **v2.1.0**: Added ledger tracking and comprehensive audit trail