# Database Management Scripts

This directory contains cross-platform database management tools written in Go.

## Setup Database

Creates the database and runs migrations:

```bash
go run cmd/setup/main.go
```

This will:
1. Connect to PostgreSQL
2. Create the `poker_planning` database
3. Create all tables (sessions, users, planning_items, votes)
4. Create indexes and triggers

## Reset Database

Drops all tables and resets the database:

```bash
go run cmd/reset/main.go
```

This will:
1. Drop all tables
2. Drop all functions
3. Clean the database completely

After reset, run setup again:
```bash
go run cmd/setup/main.go
```

## Configuration

Both scripts use the following default configuration:
- Host: `localhost`
- Port: `5432`
- User: `postgres`
- Password: `1234`
- Database: `poker_planning`

To change these values, edit the constants in the respective `main.go` files.

## Prerequisites

- PostgreSQL server running
- Go 1.21 or higher
- `github.com/lib/pq` driver (automatically downloaded by `go run`)

## Troubleshooting

### Connection refused
Make sure PostgreSQL is running:
```bash
# Windows
Get-Service postgresql*

# Linux/Mac
sudo systemctl status postgresql
```

### Authentication failed
Check your PostgreSQL password and update the `password` constant in the scripts.

### Database already exists
This is normal! The setup script will continue and create the schema.
