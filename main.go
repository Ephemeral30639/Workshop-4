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

// Transfer represents a point transfer between users
type Transfer struct {
	IdemKey     string  `json:"idemKey"`
	TransferID  int     `json:"transferId,omitempty"`
	FromUserID  int     `json:"fromUserId"`
	ToUserID    int     `json:"toUserId"`
	Amount      int     `json:"amount"`
	Status      string  `json:"status"`
	Note        *string `json:"note,omitempty"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	CompletedAt *string `json:"completedAt,omitempty"`
	FailReason  *string `json:"failReason,omitempty"`
}

// TransferCreateRequest represents the request body for creating a transfer
type TransferCreateRequest struct {
	FromUserID int    `json:"fromUserId"`
	ToUserID   int    `json:"toUserId"`
	Amount     int    `json:"amount"`
	Note       string `json:"note,omitempty"`
}

// TransferCreateResponse represents the response for creating a transfer
type TransferCreateResponse struct {
	Transfer Transfer `json:"transfer"`
}

// TransferGetResponse represents the response for getting a transfer by ID
type TransferGetResponse struct {
	Transfer Transfer `json:"transfer"`
}

// TransferListResponse represents the response for listing transfers
type TransferListResponse struct {
	Data     []Transfer `json:"data"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Total    int        `json:"total"`
}

// PointLedgerEntry represents an entry in the point ledger
type PointLedgerEntry struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Change       int    `json:"change"`
	BalanceAfter int    `json:"balance_after"`
	EventType    string `json:"event_type"`
	TransferID   *int   `json:"transfer_id,omitempty"`
	Reference    string `json:"reference,omitempty"`
	Metadata     string `json:"metadata,omitempty"`
	CreatedAt    string `json:"created_at"`
}

var db *sql.DB

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create users table
	createUsersTableSQL := `
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

	_, err = db.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	// Create transfers table
	createTransfersTableSQL := `
	CREATE TABLE IF NOT EXISTS transfers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_user_id INTEGER NOT NULL,
		to_user_id INTEGER NOT NULL,
		amount INTEGER NOT NULL CHECK (amount > 0),
		status TEXT NOT NULL CHECK (status IN ('pending','processing','completed','failed','cancelled','reversed')),
		note TEXT,
		idempotency_key TEXT NOT NULL UNIQUE,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		completed_at TEXT,
		fail_reason TEXT,
		FOREIGN KEY (from_user_id) REFERENCES users(id),
		FOREIGN KEY (to_user_id) REFERENCES users(id)
	);`

	_, err = db.Exec(createTransfersTableSQL)
	if err != nil {
		log.Fatal("Failed to create transfers table:", err)
	}

	// Create indexes for transfers table
	indexesSQL := []string{
		`CREATE INDEX IF NOT EXISTS idx_transfers_from ON transfers(from_user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_transfers_to ON transfers(to_user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_transfers_created ON transfers(created_at);`,
	}

	for _, indexSQL := range indexesSQL {
		_, err = db.Exec(indexSQL)
		if err != nil {
			log.Fatal("Failed to create transfer indexes:", err)
		}
	}

	// Create point_ledger table
	createPointLedgerTableSQL := `
	CREATE TABLE IF NOT EXISTS point_ledger (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		change INTEGER NOT NULL,
		balance_after INTEGER NOT NULL,
		event_type TEXT NOT NULL CHECK (event_type IN ('transfer_out','transfer_in','adjust','earn','redeem')),
		transfer_id INTEGER,
		reference TEXT,
		metadata TEXT,
		created_at TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (transfer_id) REFERENCES transfers(id)
	);`

	_, err = db.Exec(createPointLedgerTableSQL)
	if err != nil {
		log.Fatal("Failed to create point_ledger table:", err)
	}

	// Create indexes for point_ledger table
	ledgerIndexesSQL := []string{
		`CREATE INDEX IF NOT EXISTS idx_ledger_user ON point_ledger(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_transfer ON point_ledger(transfer_id);`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_created ON point_ledger(created_at);`,
	}

	for _, indexSQL := range ledgerIndexesSQL {
		_, err = db.Exec(indexSQL)
		if err != nil {
			log.Fatal("Failed to create point ledger indexes:", err)
		}
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

	// Transfer routes
	app.Post("/transfers", createTransfer)
	app.Get("/transfers/:id", getTransferByID)
	app.Get("/transfers", getTransfers)

	log.Println("Server starting on :3000")
	app.Listen(":3000")
}