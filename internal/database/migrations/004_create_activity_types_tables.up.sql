-- Create feed activities table
CREATE TABLE IF NOT EXISTS feed_activities (
    activity_id VARCHAR(36) PRIMARY KEY,
    feed_type VARCHAR(20) NOT NULL,
    amount_ml DECIMAL(5,1),
    duration_minutes INTEGER,
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE
);

-- Create pump activities table
CREATE TABLE IF NOT EXISTS pump_activities (
    activity_id VARCHAR(36) PRIMARY KEY,
    breast VARCHAR(10),
    amount_ml DECIMAL(5,1),
    duration_minutes INTEGER,
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE
);

-- Create diaper activities table
CREATE TABLE IF NOT EXISTS diaper_activities (
    activity_id VARCHAR(36) PRIMARY KEY,
    wet BOOLEAN DEFAULT FALSE,
    dirty BOOLEAN DEFAULT FALSE,
    color VARCHAR(20),
    consistency VARCHAR(20),
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE
);

-- Create sleep activities table
CREATE TABLE IF NOT EXISTS sleep_activities (
    activity_id VARCHAR(36) PRIMARY KEY,
    location VARCHAR(50),
    quality INTEGER CHECK (quality >= 1 AND quality <= 5),
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE
);

-- Create growth measurements table
CREATE TABLE IF NOT EXISTS growth_measurements (
    activity_id VARCHAR(36) PRIMARY KEY,
    weight_kg DECIMAL(5,2),
    height_cm DECIMAL(5,1),
    head_circumference_cm DECIMAL(4,1),
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE
);

-- Create health records table
CREATE TABLE IF NOT EXISTS health_records (
    activity_id VARCHAR(36) PRIMARY KEY,
    record_type VARCHAR(20) NOT NULL,
    provider VARCHAR(100),
    vaccine_name VARCHAR(100),
    symptoms TEXT,
    treatment TEXT,
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE
);

-- Create milestones table
CREATE TABLE IF NOT EXISTS milestones (
    activity_id VARCHAR(36) PRIMARY KEY,
    milestone_type VARCHAR(50) NOT NULL,
    description TEXT,
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE
);