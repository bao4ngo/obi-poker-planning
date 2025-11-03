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
		SSLMode:  getEnv("DB_SSLMODE", "disable"), // Use 'require' for Neon, 'disable' for local
	}

	// Try to connect to database, but don't fail if it doesn't work immediately
	if err := db.InitDB(dbConfig); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Println("Server will start anyway. Database connection will be retried on first request.")
	} else {
		log.Println("Successfully connected to database")
		defer db.CloseDB()
	}

	router := mux.NewRouter()

	// Health check endpoint for Cloud Run
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	}).Methods("GET")

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
