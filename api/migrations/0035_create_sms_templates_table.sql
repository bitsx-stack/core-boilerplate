CREATE TABLE IF NOT EXISTS sms_templates (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL,

    deleted_at TIMESTAMP NULL,
    deleted_by BIGINT NULL
);