package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://user:password@db:5432/idle_game?sslmode=disable"
	}

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	fmt.Println("[Database] Connected successfully to PostgreSQL")

	// Create tables if they don't exist (Simple Migration)
	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS characters (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			level INT DEFAULT 1,
			experience BIGINT DEFAULT 0,
			gold INT DEFAULT 0,
			strength INT DEFAULT 10,
			agility INT DEFAULT 10,
			vitality INT DEFAULT 10,
			hp INT DEFAULT 100,
			last_sync BIGINT
		);`,
		`CREATE TABLE IF NOT EXISTS items (
			id UUID PRIMARY KEY,
			character_id UUID REFERENCES characters(id),
			template_id TEXT NOT NULL,
			tier INT DEFAULT 0,
			level INT DEFAULT 1
		);`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatalf("Failed to create tables: %v", err)
		}
	}
	fmt.Println("[Database] Tables verified/created")
}
