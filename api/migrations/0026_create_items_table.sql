CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,

    unit_cost NUMERIC(12,2) DEFAULT 0,
    uom VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    tags TEXT[],

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL,

    deleted_at TIMESTAMP NULL,
    deleted_by BIGINT NULL
);