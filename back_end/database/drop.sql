-- Drop all tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS planning_items CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;
