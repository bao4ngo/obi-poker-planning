# Database Setup

## PostgreSQL Configuration

- **Host**: localhost
- **Port**: 5432
- **Database**: poker_planning
- **Username**: postgres
- **Password**: 1234

## Setup Instructions

### 1. Create the Database

Run the schema creation script:

```bash
psql -U postgres -h localhost -p 5432 -f database/schema.sql
```

Or manually:

```bash
psql -U postgres -h localhost -p 5432
```

Then run the SQL commands from `schema.sql`.

### 2. Environment Variables

The application will use these default connection settings:
- DB_HOST=localhost
- DB_PORT=5432
- DB_USER=postgres
- DB_PASSWORD=1234
- DB_NAME=poker_planning

### 3. Drop Database (if needed)

To reset the database:

```bash
psql -U postgres -h localhost -p 5432 -f database/drop.sql
```

## Database Schema

### Tables

1. **sessions** - Stores planning sessions
   - id (UUID, PK)
   - name (VARCHAR)
   - host_id (UUID)
   - current_item_id (UUID, nullable)
   - created_at (TIMESTAMP)
   - updated_at (TIMESTAMP)

2. **users** - Stores session participants
   - id (UUID, PK)
   - session_id (UUID, FK -> sessions)
   - name (VARCHAR)
   - is_host (BOOLEAN)
   - connected (BOOLEAN)
   - created_at (TIMESTAMP)
   - UNIQUE(session_id, name) - Prevents duplicate names per session

3. **planning_items** - Stores items to be estimated
   - id (UUID, PK)
   - session_id (UUID, FK -> sessions)
   - title (VARCHAR)
   - description (TEXT)
   - revealed (BOOLEAN)
   - final_estimate (VARCHAR)
   - created_at (TIMESTAMP)
   - item_order (INTEGER)

4. **votes** - Stores user votes for items
   - id (SERIAL, PK)
   - planning_item_id (UUID, FK -> planning_items)
   - user_id (UUID, FK -> users)
   - vote (VARCHAR)
   - created_at (TIMESTAMP)
   - UNIQUE(planning_item_id, user_id) - One vote per user per item

## Maintenance

### View Active Sessions

```sql
SELECT s.*, COUNT(DISTINCT u.id) as user_count, COUNT(DISTINCT pi.id) as item_count
FROM sessions s
LEFT JOIN users u ON s.id = u.session_id
LEFT JOIN planning_items pi ON s.id = pi.session_id
GROUP BY s.id;
```

### View Session Details

```sql
SELECT * FROM sessions WHERE id = 'session-uuid';
SELECT * FROM users WHERE session_id = 'session-uuid';
SELECT * FROM planning_items WHERE session_id = 'session-uuid';
```

### Clean Old Sessions

```sql
DELETE FROM sessions WHERE created_at < NOW() - INTERVAL '7 days';
```
