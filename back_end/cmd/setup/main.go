package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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
	// First, connect to postgres database to create our database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, user, password)

	log.Println("Connecting to PostgreSQL...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to PostgreSQL:", err)
	}
	defer db.Close()

	// Test connection
	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging PostgreSQL:", err)
	}
	log.Println("✓ Connected to PostgreSQL")

	// Create database if it doesn't exist
	log.Println("Creating database...")
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		// Database might already exist
		log.Printf("Database creation: %v (might already exist)\n", err)
	} else {
		log.Println("✓ Database created")
	}

	// Close connection to default database
	db.Close()

	// Connect to our new database
	connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Println("Connecting to poker_planning database...")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to poker_planning database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging poker_planning database:", err)
	}
	log.Println("✓ Connected to poker_planning database")

	// Read and execute schema
	log.Println("Creating database schema...")
	schema, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatal("Error reading schema.sql:", err)
	}

	// Execute schema
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal("Error executing schema:", err)
	}

	log.Println("✓ Database schema created successfully!")
	log.Println("\n✅ Database setup complete!")
	log.Println("\nYou can now run the application with: go run main.go")
}
