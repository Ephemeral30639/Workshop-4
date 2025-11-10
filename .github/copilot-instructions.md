# LBK Membership User Management & Point Transfer API - Copilot Instructions

## Repository Overview

This is a Go-based RESTful API using the Fiber v2 framework for managing LBK membership users with point transfer functionality. The project implements a comprehensive point transfer system with SQLite database, transaction integrity, and audit logging.

### High-Level Details
- **Project Type**: Go REST API server using Fiber v2 framework
- **Database**: SQLite with CGO dependency (requires CGO_ENABLED=1)
- **Language**: Go 1.21+ 
- **Framework**: Fiber v2.52.0
- **Key Dependencies**: SQLite3 driver, Google UUID for idempotency
- **Repository Size**: Small (~15 files, single module)
- **Target Runtime**: Standalone HTTP server on port 3000

## Build & Run Instructions

**CRITICAL**: This project requires Go modules and CGO enabled due to SQLite dependency.

### Environment Setup (ALWAYS Required)
```bash
# Enable Go modules (required on systems with GO111MODULE=off)
export GO111MODULE=on
export CGO_ENABLED=1
```

### Build Steps (Execute in order)
1. **Install Dependencies** (always run first):
   ```bash
   go mod tidy
   ```

2. **Run Application**:
   ```bash
   CGO_ENABLED=1 GO111MODULE=on go run .
   ```
   - Server starts on localhost:3000
   - Database initializes automatically (creates users.db)
   - Creates three tables: users, transfers, point_ledger

### Validation Commands
- **Test API is running**: `curl http://localhost:3000/`
- **Expected response**: `{"message": "LBK Membership API Server", "version": "1.0.0"}`

### Testing Scripts (All executable, use Python3 for JSON formatting)
```bash
# Make scripts executable
chmod +x test_api.sh test_transfer_feature.sh add_10_users.sh test_beautified.sh

# Basic API functionality test
./test_api.sh

# Comprehensive point transfer testing  
./test_transfer_feature.sh

# Add sample data (10 users with various membership levels)
./add_10_users.sh

# Test JSON beautification feature
./test_beautified.sh
```

### Build Validation Checklist
1. ✅ Go 1.21+ installed
2. ✅ `GO111MODULE=on` set
3. ✅ `CGO_ENABLED=1` set
4. ✅ `go mod tidy` completes without errors
5. ✅ Server starts and shows Fiber banner on port 3000
6. ✅ Database file `users.db` created automatically
7. ✅ Root endpoint returns JSON with API info

## Project Architecture & Layout

### Core Files Structure
```
├── main.go              # Entry point, database initialization, routes
├── handlers.go          # HTTP handlers for all API endpoints
├── go.mod              # Module dependencies (Fiber, SQLite3, UUID)
├── go.sum              # Dependency checksums
└── users.db            # SQLite database (auto-created)
```

### Key Source Files

**main.go**: Application entry point
- Database connection and table creation
- Fiber app configuration with JSON beautification
- Route definitions for users and transfers
- Struct definitions for User, Transfer, PointLedgerEntry

**handlers.go**: Request handlers
- User CRUD operations (GET, POST, PUT, DELETE /users)
- Point transfer operations (POST /transfers, GET /transfers)  
- Transaction management for atomic point transfers
- Input validation and business rule enforcement

### Database Schema (SQLite)
- **users**: Member information and point balances
- **transfers**: Point transfer transactions with idempotency
- **point_ledger**: Complete audit trail of point movements

### API Endpoints
- `GET /` - API information
- `GET /users` - List users with count
- `GET /users/{id}` - Get user by ID
- `POST /users` - Create user
- `PUT /users/{id}` - Update user (partial updates)
- `DELETE /users/{id}` - Delete user
- `POST /transfers` - Create point transfer
- `GET /transfers/{idempotencyKey}` - Get transfer by key
- `GET /transfers?userId={id}&page={page}&pageSize={size}` - List user transfers

### Business Rules & Constraints
1. **Point Transfers**:
   - Amount must be positive (> 0)
   - No self-transfers allowed
   - Sender must have sufficient balance
   - Both users must exist
   - Atomic operations with rollback on failure

2. **Data Validation**:
   - Required: first_name, last_name, member_id
   - Unique: member_id, email
   - Membership levels: Bronze, Silver, Gold, Platinum
   - Point balances cannot be negative

3. **Transaction Integrity**:
   - All transfers use database transactions
   - Idempotency keys prevent duplicate transfers
   - Complete audit trail in point_ledger table

### Dependencies
- **github.com/gofiber/fiber/v2**: HTTP framework
- **github.com/mattn/go-sqlite3**: SQLite driver (requires CGO)
- **github.com/google/uuid**: UUID generation for idempotency

### Configuration
- **Port**: 3000 (hardcoded)
- **Database**: SQLite file `users.db` in project root
- **JSON Output**: Beautified with 2-space indentation
- **CORS**: Enabled for all origins

## Testing & Quality Assurance

### Test Scripts Overview
1. **test_api.sh**: Basic CRUD operations for users
2. **test_transfer_feature.sh**: Comprehensive transfer testing with error cases
3. **add_10_users.sh**: Sample data creation (mixed Thai/English names)
4. **test_beautified.sh**: JSON formatting verification

### Common Issues & Troubleshooting
1. **CGO Error**: Ensure `CGO_ENABLED=1` is set
2. **Module Error**: Set `GO111MODULE=on` 
3. **Port 3000 in use**: Kill existing processes with `pkill -f "go run"`
4. **Permission denied on scripts**: Run `chmod +x *.sh`
5. **Python JSON formatting**: Requires `python3` with json module

### Validation Steps for Changes
1. Always run `go mod tidy` after dependency changes
2. Test with provided scripts after modifications
3. Verify database integrity with sample data
4. Check error handling with invalid inputs
5. Validate JSON response formatting

## Key Development Notes

### Code Patterns
- Error responses use structured JSON with error codes
- Database operations use prepared statements
- Transactions for atomic operations
- Null handling for optional fields (note, completed_at, fail_reason)
- Pagination support with default limits

### Performance Considerations
- Database indexes on frequently queried fields
- Pagination for large result sets (max 200 per page)
- Connection reuse for SQLite
- Efficient JSON marshaling with indentation

### Security Features
- Input validation on all endpoints
- Business rule enforcement
- SQL injection prevention via prepared statements
- Idempotency for reliable operations

---

**Trust these instructions**: Only search for additional information if these instructions are incomplete or found to be incorrect. This documentation covers all essential aspects for efficient development on this codebase.