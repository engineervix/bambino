-- scripts/init-db.sql
-- PostgreSQL initialization script for baby-tracker

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom types if needed (GORM will handle most of this)
-- This file is mainly for any custom database setup that GORM can't handle

-- Example: Create a function for updating updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Note: The actual tables will be created by GORM's AutoMigrate
-- This script is for PostgreSQL-specific features that GORM doesn't handle

-- Grant permissions (if needed)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO babytracker;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO babytracker;

-- Create indexes for better performance (if not created by GORM)
-- These will be created after tables exist, so they might fail on first run
-- CREATE INDEX IF NOT EXISTS idx_activities_baby_id_start_time ON activities(baby_id, start_time DESC);
-- CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
-- CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);

-- Sample data for development (optional)
-- Uncomment if you want some test data
/*
-- Create a test user (password: test123)
INSERT INTO users (id, username, password_hash, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'demo',
    '$2a$10$ZxPr3dMKLdMHpK5FY0dq6OAhkEjCLXrXxqZ4lU9Dd1n5K9wKZ3Jxu',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (username) DO NOTHING;

-- Create a test baby
INSERT INTO babies (id, user_id, name, birth_date, birth_weight, birth_height, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    id,
    'Demo Baby',
    CURRENT_DATE - INTERVAL '6 months',
    3.5,
    50.0,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM users WHERE username = 'demo'
ON CONFLICT DO NOTHING;
*/