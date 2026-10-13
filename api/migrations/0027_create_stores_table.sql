CREATE TABLE IF NOT EXISTS stores (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    category VARCHAR(50),

    location TEXT,
    owner BIGINT,

    status VARCHAR(20) DEFAULT 'active',
    tags TEXT[],

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL,

    deleted_at TIMESTAMP NULL,
    deleted_by BIGINT NULL
);