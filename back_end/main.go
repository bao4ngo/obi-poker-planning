package main

import (
	"log"
	"net/http"
	"os"
	"poker-planning-api/db"
	"poker-planning-api/handlers"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Initialize database connection from environment variables
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	dbConfig := db.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "1234"),
		DBName:   getEnv("DB_NAME", "poker_planning"),
	}

	if err := db.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/api/sessions", handlers.CreateSession).Methods("POST")
	router.HandleFunc("/api/sessions", handlers.GetSessions).Methods("GET")
	router.HandleFunc("/api/sessions/{sessionId}", handlers.GetSession).Methods("GET")
	router.HandleFunc("/api/sessions/{sessionId}/items", handlers.AddItem).Methods("POST")
	router.HandleFunc("/api/sessions/{sessionId}/current-item", handlers.SetCurrentItem).Methods("POST")

	// WebSocket route
	router.HandleFunc("/ws/{sessionId}", handlers.HandleWebSocket)

	// CORS configuration - allow multiple origins
	allowedOrigins := strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ",")
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
