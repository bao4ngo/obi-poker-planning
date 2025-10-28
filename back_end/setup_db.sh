#!/bin/bash

# PostgreSQL Database Setup Script for Linux/Mac

echo "Setting up PostgreSQL database for Poker Planning..."

# Database connection details
export PGPASSWORD="1234"
HOST="localhost"
PORT="5432"
USER="postgres"
DATABASE="poker_planning"

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "Error: psql command not found. Please install PostgreSQL client tools."
    exit 1
fi

echo "Creating database..."
psql -h "$HOST" -p "$PORT" -U "$USER" -c "CREATE DATABASE $DATABASE;"

if [ $? -eq 0 ]; then
    echo "Database created successfully!"
else
    echo "Database might already exist or there was an error. Continuing..."
fi

echo "Running schema migration..."
psql -h "$HOST" -p "$PORT" -U "$USER" -d "$DATABASE" -f "database/schema.sql"

if [ $? -eq 0 ]; then
    echo "Schema created successfully!"
    echo "Database setup complete!"
else
    echo "Error running schema migration."
    exit 1
fi

# Clear password from environment
unset PGPASSWORD

echo ""
echo "✅ Setup complete! You can now run: go run main.go"
