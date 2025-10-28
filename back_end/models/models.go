package models

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// User represents a participant in a planning session
type User struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	IsHost    bool            `json:"isHost"`
	Vote      string          `json:"vote,omitempty"`
	Connected bool            `json:"connected"`
	Conn      *websocket.Conn `json:"-"`
}

// PlanningItem represents a single item to be estimated
type PlanningItem struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Votes         map[string]string `json:"votes"` // userID -> vote
	Revealed      bool              `json:"revealed"`
	FinalEstimate string            `json:"finalEstimate,omitempty"`
}

// Session represents a poker planning session
type Session struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	HostID        string           `json:"hostId"`
	Users         map[string]*User `json:"users"`
	Items         []PlanningItem   `json:"items"`
	CurrentItemID string           `json:"currentItemId,omitempty"`
	CreatedAt     time.Time        `json:"createdAt"`
	Mutex         sync.RWMutex     `json:"-"`
}

// Message types for WebSocket communication
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Available card values for poker planning
var CardValues = []string{"0", "1", "2", "3", "5", "8", "13", "21", "34", "55", "89", "?"}

// NewSession creates a new planning session
func NewSession(id, name, hostID string) *Session {
	return &Session{
		ID:        id,
		Name:      name,
		HostID:    hostID,
		Users:     make(map[string]*User),
		Items:     []PlanningItem{},
		CreatedAt: time.Now(),
	}
}

// AddUser adds a user to the session
func (s *Session) AddUser(user *User) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Users[user.ID] = user
}

// RemoveUser removes a user from the session
func (s *Session) RemoveUser(userID string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if user, exists := s.Users[userID]; exists {
		user.Connected = false
		delete(s.Users, userID)
	}
}

// GetUser gets a user by ID
func (s *Session) GetUser(userID string) (*User, bool) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	user, exists := s.Users[userID]
	return user, exists
}

// AddItem adds a planning item to the session
func (s *Session) AddItem(item PlanningItem) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	item.Votes = make(map[string]string)
	s.Items = append(s.Items, item)
}

// GetCurrentItem returns the current item being voted on
func (s *Session) GetCurrentItem() *PlanningItem {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	if s.CurrentItemID == "" {
		return nil
	}

	for i := range s.Items {
		if s.Items[i].ID == s.CurrentItemID {
			return &s.Items[i]
		}
	}
	return nil
}

// UpdateItemVote updates a user's vote for the current item
func (s *Session) UpdateItemVote(itemID, userID, vote string) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for i := range s.Items {
		if s.Items[i].ID == itemID {
			s.Items[i].Votes[userID] = vote
			return true
		}
	}
	return false
}

// RevealVotes reveals all votes for an item
func (s *Session) RevealVotes(itemID string) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for i := range s.Items {
		if s.Items[i].ID == itemID {
			s.Items[i].Revealed = true
			return true
		}
	}
	return false
}

// ResetVotes resets all votes for an item
func (s *Session) ResetVotes(itemID string) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for i := range s.Items {
		if s.Items[i].ID == itemID {
			s.Items[i].Votes = make(map[string]string)
			s.Items[i].Revealed = false
			return true
		}
	}
	return false
}

// SetFinalEstimate sets the final estimate for an item
func (s *Session) SetFinalEstimate(itemID, estimate string) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for i := range s.Items {
		if s.Items[i].ID == itemID {
			s.Items[i].FinalEstimate = estimate
			return true
		}
	}
	return false
}
