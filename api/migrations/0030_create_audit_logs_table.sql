CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,

    actor BIGINT NOT NULL,
    action VARCHAR(50) NOT NULL,   -- CREATE, UPDATE, DELETE
    entity VARCHAR(50) NOT NULL,   -- table name

    entity_id BIGINT,
    description TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);