CREATE TABLE IF NOT EXISTS category_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by BIGINT
);