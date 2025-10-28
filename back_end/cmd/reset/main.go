package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "poker_planning"
)

func main() {
	// Connect to the database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Println("Connecting to PostgreSQL...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to PostgreSQL:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging PostgreSQL:", err)
	}
	log.Println("✓ Connected to database")

	// Drop all tables
	log.Println("Dropping all tables...")

	dropStatements := []string{
		"DROP TABLE IF EXISTS votes CASCADE",
		"DROP TABLE IF EXISTS planning_items CASCADE",
		"DROP TABLE IF EXISTS users CASCADE",
		"DROP TABLE IF EXISTS sessions CASCADE",
		"DROP FUNCTION IF EXISTS update_updated_at_column CASCADE",
	}

	for _, stmt := range dropStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			log.Printf("Warning: %v\n", err)
		}
	}

	log.Println("✓ All tables dropped successfully!")
	log.Println("\n✅ Database reset complete!")
	log.Println("\nRun setup again with: go run cmd/setup/main.go")
}
