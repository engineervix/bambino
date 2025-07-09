-- PostgreSQL initialization script for bambino

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

