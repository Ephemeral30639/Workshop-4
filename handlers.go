package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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