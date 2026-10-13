CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    categorytype_id INT NOT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by BIGINT
)