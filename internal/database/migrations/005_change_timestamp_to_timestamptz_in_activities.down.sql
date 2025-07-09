-- Revert timestamp columns from timestamptz to timestamp in activities table
ALTER TABLE activities
    ALTER COLUMN start_time TYPE TIMESTAMP,
    ALTER COLUMN end_time TYPE TIMESTAMP,
    ALTER COLUMN created_at TYPE TIMESTAMP,
    ALTER COLUMN updated_at TYPE TIMESTAMP;
