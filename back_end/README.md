# Poker Planning API - Backend

A Golang-based backend API for Agile Poker Planning sessions with real-time WebSocket support and PostgreSQL persistence.

## Features

- Create and manage poker planning sessions
- Real-time updates via WebSocket
- User management (host and participants)
- Planning item management
- Voting system with reveal/reset functionality
- Final estimate tracking
- **PostgreSQL database for persistent storage**
- Username validation (case-insensitive, no duplicates per session)

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- PostgreSQL client tools (psql)

## Database Setup

### 1. Ensure PostgreSQL is running

Make sure your PostgreSQL server is running on:
- Host: localhost
- Port: 5432
- Username: postgres
- Password: 1234

### 2. Run the setup command (Cross-platform)

**Recommended - Using Go (works on all OS):**

```bash
go run cmd/setup/main.go
```

This will automatically:
1. Create the database
2. Create all tables
3. Set up indexes and constraints

**Alternative - Using psql (if you prefer):**

```bash
# Windows PowerShell
.\setup_db.ps1

# Linux/Mac
./setup_db.sh
```

## Installation

1. Navigate to the backend directory:
```bash
cd back_end
```

2. Install dependencies:
```bash
go mod download
```

If you encounter network issues, try:
```bash
# Windows PowerShell
$env:GOPROXY="https://goproxy.io,direct"
go mod download

# Linux/Mac
export GOPROXY=https://goproxy.io,direct
go mod download
```

3. Set up the database:
```bash
go run cmd/setup/main.go
```

## Running the Server

```bash
go run main.go
```

The server will:
1. Connect to PostgreSQL database
2. Start on `http://localhost:8080`
3. Accept WebSocket connections

## Database Configuration

The database configuration is in `main.go`:

```go
dbConfig := db.Config{
    Host:     "localhost",
    Port:     5432,
    User:     "postgres",
    Password: "1234",
    DBName:   "poker_planning",
}
```

## API Endpoints

### REST API

- `POST /api/sessions` - Create a new planning session
- `GET /api/sessions` - Get all active sessions
- `GET /api/sessions/{sessionId}` - Get session details
- `POST /api/sessions/{sessionId}/items` - Add a planning item
- `POST /api/sessions/{sessionId}/current-item` - Set the current item

### WebSocket

- `WS /ws/{sessionId}` - Connect to a session for real-time updates

## WebSocket Message Types

### Client to Server:
- `vote` - Submit a vote for an item
- `reveal_votes` - Reveal all votes (host only)
- `reset_votes` - Reset all votes (host only)
- `set_final_estimate` - Set final estimate (host only)

### Server to Client:
- `welcome` - Initial connection confirmation
- `error` - Error message (e.g., username taken)
- `user_joined` - New user joined the session
- `user_left` - User left the session
- `item_added` - New item added
- `vote_submitted` - Vote was submitted
- `votes_revealed` - Votes were revealed
- `votes_reset` - Votes were reset
- `current_item_changed` - Current item changed
- `final_estimate_set` - Final estimate was set

## Database Schema

### Tables

1. **sessions** - Planning sessions
2. **users** - Session participants (with username uniqueness per session)
3. **planning_items** - Items to estimate
4. **votes** - User votes for items

See `database/README.md` for detailed schema information.

## Project Structure

```
back_end/
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── db/
│   ├── db.go           # Database connection
│   └── queries.go      # Database queries and operations
├── database/
│   ├── schema.sql      # Database schema
│   ├── drop.sql        # Drop tables script
│   └── README.md       # Database documentation
├── handlers/
│   ├── session.go      # REST API handlers
│   └── websocket.go    # WebSocket handlers
└── models/
    └── models.go       # Data models

```

## Card Values

The default poker planning card values are: 0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, ?

## Troubleshooting

### Database Connection Issues

1. Check PostgreSQL is running:
```bash
# Windows
Get-Service postgresql*

# Linux
sudo systemctl status postgresql

# Mac
brew services list | grep postgresql
```

2. Test connection:
```bash
# Try connecting directly
psql -h localhost -p 5432 -U postgres -c "SELECT 1;"
```

3. Run database setup:
```bash
go run cmd/setup/main.go
```

### Go Module Download Issues

If you get network errors downloading modules:

```bash
# Windows PowerShell
$env:GOPROXY="https://goproxy.io,direct"
$env:GOSUMDB="off"
go mod download

# Linux/Mac
export GOPROXY=https://goproxy.io,direct
export GOSUMDB=off
go mod download
```

### Reset Database

To reset all data:

```bash
go run cmd/reset/main.go
```

Then set up again:
```bash
go run cmd/setup/main.go
```

## Development

### View Database Contents

Connect to the database:

```bash
# All platforms
psql -h localhost -p 5432 -U postgres -d poker_planning
```

Or use the Go migration tool to check tables:
```bash
go run cmd/setup/main.go
```

Then run SQL queries:
```sql
SELECT * FROM sessions;
SELECT * FROM users;
SELECT * FROM planning_items;
SELECT * FROM votes;
```

## Production Considerations

- [ ] Change database password
- [ ] Use environment variables for configuration
- [ ] Enable SSL for database connections
- [ ] Implement connection pooling limits
- [ ] Add database migrations tool
- [ ] Implement session cleanup/archiving
- [ ] Add logging and monitoring
- [ ] Restrict CORS origins
- [ ] Add rate limiting
- [ ] Implement authentication
