package handlers

import (
	"log"
	"net/http"
	"poker-planning-api/db"
	"poker-planning-api/models"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// JoinSessionMessage represents the initial message to join a session
type JoinSessionMessage struct {
	UserName string `json:"userName"`
	UserID   string `json:"userId,omitempty"`
}

// VoteMessage represents a vote submission
type VoteMessage struct {
	ItemID string `json:"itemId"`
	Vote   string `json:"vote"`
}

// ErrorMessage represents an error message to send to client
type ErrorMessage struct {
	Error string `json:"error"`
}

// isUserNameTaken checks if a username is already taken in the session (case-insensitive)
func isUserNameTaken(sessionID, userName string, excludeUserID string) bool {
	taken, err := db.IsUserNameTaken(sessionID, userName, excludeUserID)
	if err != nil {
		log.Printf("Error checking username: %v", err)
		return false
	}
	return taken
}

// HandleWebSocket handles WebSocket connections for real-time updates
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	session, exists := GetSessionByID(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Wait for the join message
	var joinMsg JoinSessionMessage
	if err := conn.ReadJSON(&joinMsg); err != nil {
		log.Printf("Failed to read join message: %v", err)
		conn.Close()
		return
	}

	// Validate username is not empty
	if strings.TrimSpace(joinMsg.UserName) == "" {
		errorMsg := models.WSMessage{
			Type:    "error",
			Payload: map[string]string{"error": "Username cannot be empty"},
		}
		conn.WriteJSON(errorMsg)
		conn.Close()
		return
	}

	// Create or retrieve user
	var user *models.User
	if joinMsg.UserID != "" {
		// Existing user reconnecting
		existingUser, err := db.GetUserByID(joinMsg.UserID)
		if err == nil {
			user = existingUser
			user.Conn = conn
			user.Connected = true
			db.UpdateUserConnection(user.ID, true)
			session.Users[user.ID] = user
		} else {
			// User ID provided but not found, check for duplicate username
			if isUserNameTaken(sessionID, joinMsg.UserName, joinMsg.UserID) {
				errorMsg := models.WSMessage{
					Type:    "error",
					Payload: map[string]string{"error": "Username is already taken in this session"},
				}
				conn.WriteJSON(errorMsg)
				conn.Close()
				return
			}
			// User ID provided but not found, create new user
			user = &models.User{
				ID:        joinMsg.UserID,
				Name:      joinMsg.UserName,
				IsHost:    joinMsg.UserID == session.HostID,
				Connected: true,
				Conn:      conn,
			}
			if err := db.CreateUser(user, sessionID); err != nil {
				log.Printf("Failed to create user: %v", err)
				errorMsg := models.WSMessage{
					Type:    "error",
					Payload: map[string]string{"error": "Failed to create user"},
				}
				conn.WriteJSON(errorMsg)
				conn.Close()
				return
			}
			session.Users[user.ID] = user
		}
	} else {
		// New user joining - check for duplicate username
		if isUserNameTaken(sessionID, joinMsg.UserName, "") {
			errorMsg := models.WSMessage{
				Type:    "error",
				Payload: map[string]string{"error": "Username is already taken in this session"},
			}
			conn.WriteJSON(errorMsg)
			conn.Close()
			return
		}

		// New user joining
		user = &models.User{
			ID:        uuid.New().String(),
			Name:      joinMsg.UserName,
			IsHost:    false,
			Connected: true,
			Conn:      conn,
		}
		if err := db.CreateUser(user, sessionID); err != nil {
			log.Printf("Failed to create user: %v", err)
			errorMsg := models.WSMessage{
				Type:    "error",
				Payload: map[string]string{"error": "Failed to create user"},
			}
			conn.WriteJSON(errorMsg)
			conn.Close()
			return
		}
		session.Users[user.ID] = user
	}

	// Send welcome message with user info and session state
	welcomeMsg := models.WSMessage{
		Type: "welcome",
		Payload: map[string]interface{}{
			"userId":  user.ID,
			"session": session,
		},
	}
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Failed to send welcome message: %v", err)
	}

	// Broadcast user joined to all other users
	BroadcastToSession(sessionID, models.WSMessage{
		Type:    "user_joined",
		Payload: user,
	})

	// Handle incoming messages
	go handleMessages(conn, session, user)
}

func handleMessages(conn *websocket.Conn, session *models.Session, user *models.User) {
	defer func() {
		// Mark user as disconnected in database
		db.UpdateUserConnection(user.ID, false)
		delete(session.Users, user.ID)
		conn.Close()

		// Broadcast user left
		BroadcastToSession(session.ID, models.WSMessage{
			Type:    "user_left",
			Payload: map[string]string{"userId": user.ID},
		})
	}()

	for {
		var msg models.WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		handleMessage(session, user, msg)
	}
}

func handleMessage(session *models.Session, user *models.User, msg models.WSMessage) {
	switch msg.Type {
	case "vote":
		handleVote(session, user, msg)
	case "reveal_votes":
		handleRevealVotes(session, user, msg)
	case "reset_votes":
		handleResetVotes(session, user, msg)
	case "set_final_estimate":
		handleSetFinalEstimate(session, user, msg)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func handleVote(session *models.Session, user *models.User, msg models.WSMessage) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}

	itemID, _ := payload["itemId"].(string)
	vote, _ := payload["vote"].(string)

	if itemID == "" || vote == "" {
		return
	}

	// Save vote to database
	if err := db.SaveVote(itemID, user.ID, vote); err != nil {
		log.Printf("Failed to save vote: %v", err)
		return
	}

	// Broadcast vote update (without revealing the vote value)
	BroadcastToSession(session.ID, models.WSMessage{
		Type: "vote_submitted",
		Payload: map[string]interface{}{
			"itemID":   itemID,
			"userId":   user.ID,
			"hasVoted": true,
		},
	})
}

func handleRevealVotes(session *models.Session, user *models.User, msg models.WSMessage) {
	if !user.IsHost {
		return
	}

	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}

	itemID, _ := payload["itemId"].(string)
	if itemID == "" {
		return
	}

	// Update in database
	if err := db.UpdateItemRevealed(itemID, true); err != nil {
		log.Printf("Failed to reveal votes: %v", err)
		return
	}

	// Get the updated item with votes from database
	item, err := db.GetPlanningItemByID(itemID)
	if err != nil {
		log.Printf("Failed to get item: %v", err)
		return
	}

	BroadcastToSession(session.ID, models.WSMessage{
		Type:    "votes_revealed",
		Payload: item,
	})
}

func handleResetVotes(session *models.Session, user *models.User, msg models.WSMessage) {
	if !user.IsHost {
		return
	}

	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}

	itemID, _ := payload["itemId"].(string)
	if itemID == "" {
		return
	}

	// Delete all votes from database
	if err := db.DeleteItemVotes(itemID); err != nil {
		log.Printf("Failed to delete votes: %v", err)
		return
	}

	// Update revealed status
	if err := db.UpdateItemRevealed(itemID, false); err != nil {
		log.Printf("Failed to update revealed status: %v", err)
		return
	}

	BroadcastToSession(session.ID, models.WSMessage{
		Type:    "votes_reset",
		Payload: map[string]string{"itemId": itemID},
	})
}

func handleSetFinalEstimate(session *models.Session, user *models.User, msg models.WSMessage) {
	if !user.IsHost {
		return
	}

	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}

	itemID, _ := payload["itemId"].(string)
	estimate, _ := payload["estimate"].(string)

	if itemID == "" {
		return
	}

	// Update in database
	if err := db.UpdateItemFinalEstimate(itemID, estimate); err != nil {
		log.Printf("Failed to set final estimate: %v", err)
		return
	}

	BroadcastToSession(session.ID, models.WSMessage{
		Type: "final_estimate_set",
		Payload: map[string]string{
			"itemId":   itemID,
			"estimate": estimate,
		},
	})
}

// BroadcastToSession sends a message to all connected users in a session
func BroadcastToSession(sessionID string, msg models.WSMessage) {
	session, exists := GetSessionByID(sessionID)
	if !exists {
		return
	}

	session.Mutex.RLock()
	defer session.Mutex.RUnlock()

	for _, user := range session.Users {
		if user.Connected && user.Conn != nil {
			if err := user.Conn.WriteJSON(msg); err != nil {
				log.Printf("Failed to send message to user %s: %v", user.ID, err)
			}
		}
	}
}
