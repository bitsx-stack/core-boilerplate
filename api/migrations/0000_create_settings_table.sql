CREATE TABLE IF NOT EXISTS settings (
    id BIGSERIAL PRIMARY KEY,
    
    code VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    value TEXT,
    category VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL,
    
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    deleted_by BIGINT NULL
);

-- Indexes
CREATE INDEX idx_settings_name ON settings(name);
CREATE INDEX idx_settings_code ON settings(code);
CREATE INDEX idx_settings_category ON settings(category);