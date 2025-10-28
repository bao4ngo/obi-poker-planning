package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"poker-planning-api/db"
	"poker-planning-api/models"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	// In-memory cache for active WebSocket connections
	activeSessions = make(map[string]*models.Session)
	sessionsMutex  sync.RWMutex
)

// CreateSessionRequest represents the request to create a new session
type CreateSessionRequest struct {
	Name     string `json:"name"`
	HostName string `json:"hostName"`
}

// CreateSessionResponse represents the response after creating a session
type CreateSessionResponse struct {
	SessionID string `json:"sessionId"`
	HostID    string `json:"hostId"`
}

// CreateSession handles creating a new poker planning session
func CreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.HostName == "" {
		http.Error(w, "Session name and host name are required", http.StatusBadRequest)
		return
	}

	sessionID := uuid.New().String()
	hostID := uuid.New().String()

	session := models.NewSession(sessionID, req.Name, hostID)

	// Save session to database
	if err := db.CreateSession(session); err != nil {
		log.Printf("Failed to create session: %v", err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Create host user
	host := &models.User{
		ID:        hostID,
		Name:      req.HostName,
		IsHost:    true,
		Connected: false,
	}

	if err := db.CreateUser(host, sessionID); err != nil {
		log.Printf("Failed to create host user: %v", err)
		http.Error(w, "Failed to create host", http.StatusInternalServerError)
		return
	}

	// Cache session for WebSocket connections
	sessionsMutex.Lock()
	activeSessions[sessionID] = session
	activeSessions[sessionID].Users[hostID] = host
	sessionsMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateSessionResponse{
		SessionID: sessionID,
		HostID:    hostID,
	})
}

// GetSession returns session information
func GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	// Try to get from database
	session, err := db.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// AddItemRequest represents the request to add a planning item
type AddItemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AddItem adds a new planning item to a session
func AddItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify session exists
	_, err := db.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	item := models.PlanningItem{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Votes:       make(map[string]string),
		Revealed:    false,
	}

	// Save item to database
	if err := db.CreatePlanningItem(&item, sessionID); err != nil {
		log.Printf("Failed to create item: %v", err)
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	// Broadcast the update to all connected clients
	BroadcastToSession(sessionID, models.WSMessage{
		Type:    "item_added",
		Payload: item,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// SetCurrentItemRequest represents the request to set the current item
type SetCurrentItemRequest struct {
	ItemID string `json:"itemId"`
}

// SetCurrentItem sets the current item being voted on
func SetCurrentItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var req SetCurrentItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update in database
	if err := db.UpdateSessionCurrentItem(sessionID, req.ItemID); err != nil {
		log.Printf("Failed to update current item: %v", err)
		http.Error(w, "Failed to update current item", http.StatusInternalServerError)
		return
	}

	// Update in cache if exists
	sessionsMutex.Lock()
	if session, exists := activeSessions[sessionID]; exists {
		session.CurrentItemID = req.ItemID
	}
	sessionsMutex.Unlock()

	// Broadcast the update to all connected clients
	BroadcastToSession(sessionID, models.WSMessage{
		Type:    "current_item_changed",
		Payload: map[string]string{"itemId": req.ItemID},
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GetSessions returns all active sessions (for debugging)
func GetSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := db.GetAllSessions()
	if err != nil {
		log.Printf("Failed to get sessions: %v", err)
		http.Error(w, "Failed to retrieve sessions", http.StatusInternalServerError)
		return
	}

	sessionList := make([]map[string]interface{}, 0)
	for _, session := range sessions {
		users, _ := db.GetSessionUsers(session.ID)
		items, _ := db.GetSessionItems(session.ID)

		sessionList = append(sessionList, map[string]interface{}{
			"id":        session.ID,
			"name":      session.Name,
			"userCount": len(users),
			"itemCount": len(items),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessionList)
}

// GetSessionByID returns a session by ID (used internally)
func GetSessionByID(sessionID string) (*models.Session, bool) {
	// Try cache first
	sessionsMutex.RLock()
	session, exists := activeSessions[sessionID]
	sessionsMutex.RUnlock()

	if exists {
		return session, true
	}

	// Try database
	session, err := db.GetSession(sessionID)
	if err != nil {
		return nil, false
	}

	// Add to cache
	sessionsMutex.Lock()
	activeSessions[sessionID] = session
	sessionsMutex.Unlock()

	return session, true
}
