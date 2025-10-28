package db

import (
	"database/sql"
	"poker-planning-api/models"
	"time"
)

// CreateSession creates a new session in the database
func CreateSession(session *models.Session) error {
	query := `
		INSERT INTO sessions (id, name, host_id, current_item_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := DB.Exec(query, session.ID, session.Name, session.HostID,
		sql.NullString{String: session.CurrentItemID, Valid: session.CurrentItemID != ""},
		session.CreatedAt, time.Now())
	return err
}

// GetSession retrieves a session by ID
func GetSession(sessionID string) (*models.Session, error) {
	query := `SELECT id, name, host_id, current_item_id, created_at FROM sessions WHERE id = $1`

	session := &models.Session{
		Users: make(map[string]*models.User),
		Items: []models.PlanningItem{},
	}

	var currentItemID sql.NullString
	err := DB.QueryRow(query, sessionID).Scan(
		&session.ID, &session.Name, &session.HostID, &currentItemID, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if currentItemID.Valid {
		session.CurrentItemID = currentItemID.String
	}

	// Load users
	users, err := GetSessionUsers(sessionID)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		session.Users[user.ID] = user
	}

	// Load items
	items, err := GetSessionItems(sessionID)
	if err != nil {
		return nil, err
	}
	session.Items = items

	return session, nil
}

// GetAllSessions retrieves all sessions
func GetAllSessions() ([]*models.Session, error) {
	query := `SELECT id, name, host_id, current_item_id, created_at FROM sessions ORDER BY created_at DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := []*models.Session{}
	for rows.Next() {
		session := &models.Session{
			Users: make(map[string]*models.User),
			Items: []models.PlanningItem{},
		}
		var currentItemID sql.NullString
		err := rows.Scan(&session.ID, &session.Name, &session.HostID, &currentItemID, &session.CreatedAt)
		if err != nil {
			return nil, err
		}
		if currentItemID.Valid {
			session.CurrentItemID = currentItemID.String
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// UpdateSessionCurrentItem updates the current item being voted on
func UpdateSessionCurrentItem(sessionID, itemID string) error {
	query := `UPDATE sessions SET current_item_id = $1, updated_at = $2 WHERE id = $3`
	_, err := DB.Exec(query, sql.NullString{String: itemID, Valid: itemID != ""}, time.Now(), sessionID)
	return err
}

// DeleteSession deletes a session and all related data (cascades)
func DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := DB.Exec(query, sessionID)
	return err
}

// CreateUser creates a new user in the database
func CreateUser(user *models.User, sessionID string) error {
	query := `
		INSERT INTO users (id, session_id, name, is_host, connected, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := DB.Exec(query, user.ID, sessionID, user.Name, user.IsHost, user.Connected, time.Now())
	return err
}

// GetSessionUsers retrieves all users for a session
func GetSessionUsers(sessionID string) ([]*models.User, error) {
	query := `SELECT id, name, is_host, connected FROM users WHERE session_id = $1`

	rows, err := DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.IsHost, &user.Connected)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID string) (*models.User, error) {
	query := `SELECT id, name, is_host, connected FROM users WHERE id = $1`

	user := &models.User{}
	err := DB.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.IsHost, &user.Connected)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserConnection updates a user's connection status
func UpdateUserConnection(userID string, connected bool) error {
	query := `UPDATE users SET connected = $1 WHERE id = $2`
	_, err := DB.Exec(query, connected, userID)
	return err
}

// DeleteUser deletes a user from the database
func DeleteUser(userID string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := DB.Exec(query, userID)
	return err
}

// IsUserNameTaken checks if a username is already taken in a session (case-insensitive)
func IsUserNameTaken(sessionID, userName, excludeUserID string) (bool, error) {
	var query string
	var args []interface{}

	if excludeUserID == "" {
		// No user to exclude, check all users in the session
		query = `
			SELECT COUNT(*) FROM users 
			WHERE session_id = $1 AND LOWER(TRIM(name)) = LOWER(TRIM($2))
		`
		args = []interface{}{sessionID, userName}
	} else {
		// Exclude a specific user (for reconnection scenarios)
		query = `
			SELECT COUNT(*) FROM users 
			WHERE session_id = $1 AND LOWER(TRIM(name)) = LOWER(TRIM($2)) AND id != $3
		`
		args = []interface{}{sessionID, userName, excludeUserID}
	}

	var count int
	err := DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreatePlanningItem creates a new planning item
func CreatePlanningItem(item *models.PlanningItem, sessionID string) error {
	// Get the next order number
	var maxOrder int
	orderQuery := `SELECT COALESCE(MAX(item_order), 0) FROM planning_items WHERE session_id = $1`
	DB.QueryRow(orderQuery, sessionID).Scan(&maxOrder)

	query := `
		INSERT INTO planning_items (id, session_id, title, description, revealed, final_estimate, created_at, item_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := DB.Exec(query, item.ID, sessionID, item.Title, item.Description, item.Revealed,
		sql.NullString{String: item.FinalEstimate, Valid: item.FinalEstimate != ""},
		time.Now(), maxOrder+1)
	return err
}

// GetSessionItems retrieves all planning items for a session
func GetSessionItems(sessionID string) ([]models.PlanningItem, error) {
	query := `
		SELECT id, title, description, revealed, final_estimate 
		FROM planning_items 
		WHERE session_id = $1 
		ORDER BY item_order, created_at
	`

	rows, err := DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.PlanningItem{}
	for rows.Next() {
		item := models.PlanningItem{
			Votes: make(map[string]string),
		}
		var finalEstimate sql.NullString
		err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Revealed, &finalEstimate)
		if err != nil {
			return nil, err
		}
		if finalEstimate.Valid {
			item.FinalEstimate = finalEstimate.String
		}

		// Load votes for this item
		votes, err := GetItemVotes(item.ID)
		if err != nil {
			return nil, err
		}
		item.Votes = votes

		items = append(items, item)
	}

	return items, nil
}

// GetPlanningItemByID retrieves a planning item by ID
func GetPlanningItemByID(itemID string) (*models.PlanningItem, error) {
	query := `SELECT id, title, description, revealed, final_estimate FROM planning_items WHERE id = $1`

	item := &models.PlanningItem{
		Votes: make(map[string]string),
	}
	var finalEstimate sql.NullString
	err := DB.QueryRow(query, itemID).Scan(&item.ID, &item.Title, &item.Description, &item.Revealed, &finalEstimate)
	if err != nil {
		return nil, err
	}
	if finalEstimate.Valid {
		item.FinalEstimate = finalEstimate.String
	}

	// Load votes
	votes, err := GetItemVotes(item.ID)
	if err != nil {
		return nil, err
	}
	item.Votes = votes

	return item, nil
}

// UpdateItemRevealed updates the revealed status of an item
func UpdateItemRevealed(itemID string, revealed bool) error {
	query := `UPDATE planning_items SET revealed = $1 WHERE id = $2`
	_, err := DB.Exec(query, revealed, itemID)
	return err
}

// UpdateItemFinalEstimate updates the final estimate of an item
func UpdateItemFinalEstimate(itemID, estimate string) error {
	query := `UPDATE planning_items SET final_estimate = $1 WHERE id = $2`
	_, err := DB.Exec(query, sql.NullString{String: estimate, Valid: estimate != ""}, itemID)
	return err
}

// SaveVote saves or updates a user's vote for an item
func SaveVote(itemID, userID, vote string) error {
	query := `
		INSERT INTO votes (planning_item_id, user_id, vote, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (planning_item_id, user_id) 
		DO UPDATE SET vote = $3, created_at = $4
	`
	_, err := DB.Exec(query, itemID, userID, vote, time.Now())
	return err
}

// GetItemVotes retrieves all votes for a planning item
func GetItemVotes(itemID string) (map[string]string, error) {
	query := `SELECT user_id, vote FROM votes WHERE planning_item_id = $1`

	rows, err := DB.Query(query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	votes := make(map[string]string)
	for rows.Next() {
		var userID, vote string
		err := rows.Scan(&userID, &vote)
		if err != nil {
			return nil, err
		}
		votes[userID] = vote
	}

	return votes, nil
}

// DeleteItemVotes deletes all votes for a planning item
func DeleteItemVotes(itemID string) error {
	query := `DELETE FROM votes WHERE planning_item_id = $1`
	_, err := DB.Exec(query, itemID)
	return err
}
