CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    permissions JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL,

    deleted_at TIMESTAMP NULL,
    deleted_by BIGINT NULL
);