package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID              int    `json:"id"`
	MemberID        string `json:"member_id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MobileNumber    string `json:"mobile_number"`
	Email           string `json:"email"`
	RegisterDate    string `json:"register_date"`
	MembershipLevel string `json:"membership_level"`
	PointBalance    int    `json:"point_balance"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

var db *sql.DB

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create users table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		member_id TEXT UNIQUE NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		mobile_number TEXT,
		email TEXT UNIQUE,
		register_date TEXT,
		membership_level TEXT DEFAULT 'Bronze',
		point_balance INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	log.Println("Database initialized successfully")
}

func main() {
	// Initialize database
	initDatabase()
	defer db.Close()

	// Create a new Fiber app with JSON encoder configuration
	app := fiber.New(fiber.Config{
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		},
	})

	// Enable CORS
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "LBK Membership API Server",
			"version": "1.0.0",
		})
	})

	// User CRUD routes
	app.Get("/users", getUsers)
	app.Get("/users/:id", getUserByID)
	app.Post("/users", createUser)
	app.Put("/users/:id", updateUser)
	app.Delete("/users/:id", deleteUser)

	log.Println("Server starting on :3000")
	app.Listen(":3000")
}