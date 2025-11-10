package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GET /users - Get all users
func getUsers(c *fiber.Ctx) error {
	rows, err := db.Query(`
		SELECT id, member_id, first_name, last_name, mobile_number, email, 
		       register_date, membership_level, point_balance, created_at, updated_at 
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.MemberID, &user.FirstName, &user.LastName,
			&user.MobileNumber, &user.Email, &user.RegisterDate, &user.MembershipLevel,
			&user.PointBalance, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to scan user data",
			})
		}
		users = append(users, user)
	}

	return c.JSON(fiber.Map{
		"data":  users,
		"count": len(users),
	})
}

// GET /users/:id - Get user by ID
func getUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var user User
	err = db.QueryRow(`
		SELECT id, member_id, first_name, last_name, mobile_number, email, 
		       register_date, membership_level, point_balance, created_at, updated_at 
		FROM users WHERE id = ?
	`, userID).Scan(&user.ID, &user.MemberID, &user.FirstName, &user.LastName,
		&user.MobileNumber, &user.Email, &user.RegisterDate, &user.MembershipLevel,
		&user.PointBalance, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": user,
	})
}

// POST /users - Create new user
func createUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if user.FirstName == "" || user.LastName == "" || user.MemberID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "First name, last name, and member ID are required",
		})
	}

	// Set default values
	if user.MembershipLevel == "" {
		user.MembershipLevel = "Bronze"
	}
	if user.RegisterDate == "" {
		user.RegisterDate = time.Now().Format("2006-01-02")
	}

	result, err := db.Exec(`
		INSERT INTO users (member_id, first_name, last_name, mobile_number, email, 
		                   register_date, membership_level, point_balance) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, user.MemberID, user.FirstName, user.LastName, user.MobileNumber, user.Email,
		user.RegisterDate, user.MembershipLevel, user.PointBalance)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user: " + err.Error(),
		})
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)

	return c.Status(201).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})
}

// PUT /users/:id - Update user
func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if user exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	_, err = db.Exec(`
		UPDATE users SET 
		member_id = COALESCE(NULLIF(?, ''), member_id),
		first_name = COALESCE(NULLIF(?, ''), first_name),
		last_name = COALESCE(NULLIF(?, ''), last_name),
		mobile_number = COALESCE(NULLIF(?, ''), mobile_number),
		email = COALESCE(NULLIF(?, ''), email),
		register_date = COALESCE(NULLIF(?, ''), register_date),
		membership_level = COALESCE(NULLIF(?, ''), membership_level),
		point_balance = CASE WHEN ? >= 0 THEN ? ELSE point_balance END,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, user.MemberID, user.FirstName, user.LastName, user.MobileNumber, user.Email,
		user.RegisterDate, user.MembershipLevel, user.PointBalance, user.PointBalance, userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update user: " + err.Error(),
		})
	}

	// Fetch updated user
	var updatedUser User
	err = db.QueryRow(`
		SELECT id, member_id, first_name, last_name, mobile_number, email, 
		       register_date, membership_level, point_balance, created_at, updated_at 
		FROM users WHERE id = ?
	`, userID).Scan(&updatedUser.ID, &updatedUser.MemberID, &updatedUser.FirstName,
		&updatedUser.LastName, &updatedUser.MobileNumber, &updatedUser.Email,
		&updatedUser.RegisterDate, &updatedUser.MembershipLevel,
		&updatedUser.PointBalance, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch updated user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
		"data":    updatedUser,
	})
}

// DELETE /users/:id - Delete user
func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Check if user exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// POST /transfers - Create a new point transfer
func createTransfer(c *fiber.Ctx) error {
	var req TransferCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.FromUserID <= 0 || req.ToUserID <= 0 || req.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "fromUserId, toUserId, and amount must be positive integers",
		})
	}

	// Check if trying to transfer to themselves
	if req.FromUserID == req.ToUserID {
		return c.Status(422).JSON(fiber.Map{
			"error":   "BUSINESS_ERROR",
			"message": "Cannot transfer points to yourself",
		})
	}

	// Generate idempotency key
	idemKey := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to start transaction",
		})
	}
	defer tx.Rollback()

	// Check if both users exist and get their current balances
	var fromUserBalance, toUserBalance int
	err = tx.QueryRow("SELECT point_balance FROM users WHERE id = ?", req.FromUserID).Scan(&fromUserBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"error":   "NOT_FOUND",
				"message": "From user not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to check from user",
		})
	}

	err = tx.QueryRow("SELECT point_balance FROM users WHERE id = ?", req.ToUserID).Scan(&toUserBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"error":   "NOT_FOUND",
				"message": "To user not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to check to user",
		})
	}

	// Check if from user has sufficient balance
	if fromUserBalance < req.Amount {
		return c.Status(409).JSON(fiber.Map{
			"error":   "INSUFFICIENT_BALANCE",
			"message": "Insufficient point balance",
		})
	}

	// Create transfer record
	result, err := tx.Exec(`
		INSERT INTO transfers (from_user_id, to_user_id, amount, status, note, idempotency_key, created_at, updated_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.FromUserID, req.ToUserID, req.Amount, "completed", req.Note, idemKey, now, now, now)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to create transfer: " + err.Error(),
		})
	}

	transferID, _ := result.LastInsertId()

	// Update user balances
	newFromBalance := fromUserBalance - req.Amount
	newToBalance := toUserBalance + req.Amount

	_, err = tx.Exec("UPDATE users SET point_balance = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", newFromBalance, req.FromUserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to update from user balance",
		})
	}

	_, err = tx.Exec("UPDATE users SET point_balance = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", newToBalance, req.ToUserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to update to user balance",
		})
	}

	// Create ledger entries
	_, err = tx.Exec(`
		INSERT INTO point_ledger (user_id, change, balance_after, event_type, transfer_id, reference, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.FromUserID, -req.Amount, newFromBalance, "transfer_out", transferID, fmt.Sprintf("Transfer to user %d", req.ToUserID), now)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to create from user ledger entry",
		})
	}

	_, err = tx.Exec(`
		INSERT INTO point_ledger (user_id, change, balance_after, event_type, transfer_id, reference, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.ToUserID, req.Amount, newToBalance, "transfer_in", transferID, fmt.Sprintf("Transfer from user %d", req.FromUserID), now)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to create to user ledger entry",
		})
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to commit transaction",
		})
	}

	// Prepare response
	transfer := Transfer{
		IdemKey:     idemKey,
		TransferID:  int(transferID),
		FromUserID:  req.FromUserID,
		ToUserID:    req.ToUserID,
		Amount:      req.Amount,
		Status:      "completed",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	if req.Note != "" {
		transfer.Note = &req.Note
	}
	transfer.CompletedAt = &now

	// Set response header
	c.Set("Idempotency-Key", idemKey)

	return c.Status(201).JSON(TransferCreateResponse{
		Transfer: transfer,
	})
}

// GET /transfers/:id - Get transfer by idempotency key
func getTransferByID(c *fiber.Ctx) error {
	idemKey := c.Params("id")
	if idemKey == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "Transfer ID is required",
		})
	}

	var transfer Transfer
	var note, completedAt, failReason sql.NullString
	
	err := db.QueryRow(`
		SELECT idempotency_key, id, from_user_id, to_user_id, amount, status, note, 
		       created_at, updated_at, completed_at, fail_reason
		FROM transfers 
		WHERE idempotency_key = ?
	`, idemKey).Scan(&transfer.IdemKey, &transfer.TransferID, &transfer.FromUserID, 
		&transfer.ToUserID, &transfer.Amount, &transfer.Status, &note,
		&transfer.CreatedAt, &transfer.UpdatedAt, &completedAt, &failReason)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"error":   "NOT_FOUND",
				"message": "Transfer not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to fetch transfer: " + err.Error(),
		})
	}

	// Handle nullable fields
	if note.Valid {
		transfer.Note = &note.String
	}
	if completedAt.Valid {
		transfer.CompletedAt = &completedAt.String
	}
	if failReason.Valid {
		transfer.FailReason = &failReason.String
	}

	return c.JSON(TransferGetResponse{
		Transfer: transfer,
	})
}

// GET /transfers - List transfers with user filtering and pagination
func getTransfers(c *fiber.Ctx) error {
	// Get query parameters
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "userId query parameter is required",
		})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "userId must be a positive integer",
		})
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 200 {
			pageSize = ps
		}
	}

	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM transfers 
		WHERE from_user_id = ? OR to_user_id = ?
	`, userID, userID).Scan(&total)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to count transfers",
		})
	}

	// Get transfers
	rows, err := db.Query(`
		SELECT idempotency_key, id, from_user_id, to_user_id, amount, status, note,
		       created_at, updated_at, completed_at, fail_reason
		FROM transfers 
		WHERE from_user_id = ? OR to_user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, userID, userID, pageSize, offset)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to fetch transfers",
		})
	}
	defer rows.Close()

	var transfers []Transfer
	for rows.Next() {
		var transfer Transfer
		var note, completedAt, failReason sql.NullString
		
		err := rows.Scan(&transfer.IdemKey, &transfer.TransferID, &transfer.FromUserID,
			&transfer.ToUserID, &transfer.Amount, &transfer.Status, &note,
			&transfer.CreatedAt, &transfer.UpdatedAt, &completedAt, &failReason)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "INTERNAL_ERROR",
				"message": "Failed to scan transfer data",
			})
		}
		
		// Handle nullable fields
		if note.Valid {
			transfer.Note = &note.String
		}
		if completedAt.Valid {
			transfer.CompletedAt = &completedAt.String
		}
		if failReason.Valid {
			transfer.FailReason = &failReason.String
		}
		
		transfers = append(transfers, transfer)
	}

	if transfers == nil {
		transfers = []Transfer{}
	}

	return c.JSON(TransferListResponse{
		Data:     transfers,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}