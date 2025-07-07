-- Drop activities table and related indexes
DROP INDEX IF EXISTS idx_activities_baby_id;
DROP INDEX IF EXISTS idx_activities_type;
DROP INDEX IF EXISTS idx_activities_start_time;
DROP INDEX IF EXISTS idx_activities_baby_id_start_time;
DROP TABLE IF EXISTS activities;