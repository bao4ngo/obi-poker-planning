# PostgreSQL Integration - Setup Guide

## Overview

Your Poker Planning application now uses PostgreSQL for persistent data storage. All sessions, users, planning items, and votes are stored in the database.

## What Was Added

### 1. Database Schema (`back_end/database/`)
- `schema.sql` - Complete database schema with 4 tables
- `drop.sql` - Script to reset the database
- `README.md` - Database documentation

### 2. Database Package (`back_end/db/`)
- `db.go` - Database connection management
- `queries.go` - All CRUD operations for sessions, users, items, and votes

### 3. Updated Application Code
- `main.go` - Initializes database connection on startup
- `handlers/session.go` - Uses database instead of in-memory storage
- `handlers/websocket.go` - Saves all user actions to database
- `go.mod` - Added `github.com/lib/pq` PostgreSQL driver

## Setup Steps

### Step 1: Install Dependencies

```bash
cd back_end
go mod download
```

If you encounter network issues:
```bash
# Windows PowerShell
$env:GOPROXY="https://goproxy.io,direct"

# Linux/Mac
export GOPROXY=https://goproxy.io,direct

# Then run
go mod download
```

### Step 2: Setup Database (Cross-Platform)

**Using Go (Recommended - works on all OS):**
```bash
cd back_end
go run cmd/setup/main.go
```

This will:
- Create the `poker_planning` database
- Create all tables (sessions, users, planning_items, votes)
- Set up indexes and constraints
- Create triggers

**Alternative methods:**
```bash
# Windows PowerShell
.\setup_db.ps1

# Linux/Mac with psql
psql -h localhost -p 5432 -U postgres -f database/schema.sql
```

### Step 3: Run the Application

```bash
cd back_end
go run main.go
```

You should see:
```
Successfully connected to PostgreSQL database
Server starting on :8080
```

## Database Tables

### 1. sessions
Stores planning sessions
- id, name, host_id, current_item_id, created_at, updated_at

### 2. users  
Stores participants (UNIQUE constraint on session_id + name)
- id, session_id, name, is_host, connected, created_at

### 3. planning_items
Stores items to estimate
- id, session_id, title, description, revealed, final_estimate, created_at, item_order

### 4. votes
Stores user votes (UNIQUE constraint on planning_item_id + user_id)
- id, planning_item_id, user_id, vote, created_at

## Key Features

✅ **Persistent Storage** - All data survives server restarts
✅ **Username Validation** - Case-insensitive duplicate checking
✅ **Atomic Operations** - Database transactions ensure data integrity
✅ **Connection Pooling** - Efficient database connection management
✅ **Cascading Deletes** - Deleting a session removes all related data
✅ **Indexed Queries** - Optimized for performance

## Verify Setup

### 1. Check Database Connection
```bash
# All platforms - using psql
psql -h localhost -p 5432 -U postgres -d poker_planning -c "\dt"

# Or use the Go setup tool which will verify automatically
go run cmd/setup/main.go
```

Expected output:
```
              List of relations
 Schema |      Name       | Type  |  Owner   
--------+-----------------+-------+----------
 public | planning_items  | table | postgres
 public | sessions        | table | postgres
 public | users           | table | postgres
 public | votes           | table | postgres
```

### 2. Test the Application

1. Start backend: `go run main.go`
2. Start frontend: `cd ../front_end; npm run dev`
3. Create a session
4. Add planning items
5. Invite users
6. Vote on items

### 3. View Database Contents

```bash
# Connect to database
psql -h localhost -p 5432 -U postgres -d poker_planning
```

Then run:
```sql
-- View all sessions
SELECT * FROM sessions;

-- View users in a specific session
SELECT * FROM users WHERE session_id = 'your-session-id';

-- View items and votes
SELECT 
    pi.title, 
    pi.revealed,
    pi.final_estimate,
    COUNT(v.id) as vote_count
FROM planning_items pi
LEFT JOIN votes v ON pi.id = v.planning_item_id
WHERE pi.session_id = 'your-session-id'
GROUP BY pi.id, pi.title, pi.revealed, pi.final_estimate;
```

## Troubleshooting

### Error: "Failed to initialize database"

1. Check PostgreSQL is running:
```bash
# Windows
Get-Service postgresql*

# Linux
sudo systemctl status postgresql

# Mac
brew services list | grep postgresql

# Or use pg_isready
pg_isready -h localhost -p 5432
```

2. Verify credentials:
```bash
psql -h localhost -p 5432 -U postgres -c "SELECT 1;"
```

3. Check if database exists:
```bash
psql -h localhost -p 5432 -U postgres -c "\l" | grep poker_planning
```

4. Try running setup again:
```bash
go run cmd/setup/main.go
```

### Error: "Username is already taken"

This is expected behavior! The database enforces unique usernames per session (case-insensitive).

### Reset Everything

**Using Go (Cross-platform):**
```bash
go run cmd/reset/main.go
go run cmd/setup/main.go
```

**Using SQL files:**
```bash
# Windows PowerShell
$env:PGPASSWORD="1234"
psql -h localhost -p 5432 -U postgres -d poker_planning -f database/drop.sql
psql -h localhost -p 5432 -U postgres -d poker_planning -f database/schema.sql

# Linux/Mac
PGPASSWORD=1234 psql -h localhost -p 5432 -U postgres -d poker_planning -f database/drop.sql
PGPASSWORD=1234 psql -h localhost -p 5432 -U postgres -d poker_planning -f database/schema.sql
```

## What Changed from Before

| Before | After |
|--------|-------|
| In-memory storage | PostgreSQL database |
| Data lost on restart | Data persists |
| No username validation | Case-insensitive unique usernames |
| Simple map storage | Proper relational database |
| No data history | All actions tracked with timestamps |

## Next Steps

1. ✅ Run `go run cmd/setup/main.go` (cross-platform!)
2. ✅ Run `go mod download`  
3. ✅ Start the server with `go run main.go`
4. ✅ Test creating sessions and voting
5. ✅ Verify data persists after restart

Your application is now production-ready with persistent storage! 🎉
