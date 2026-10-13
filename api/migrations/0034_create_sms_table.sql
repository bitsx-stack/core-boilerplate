CREATE TABLE IF NOT EXISTS sms (
    id BIGSERIAL PRIMARY KEY,

    content TEXT NOT NULL,
    sender VARCHAR(50),
    receiver VARCHAR(50),

    status VARCHAR(20) DEFAULT 'pending',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL
);