CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    bill_id BIGINT NOT NULL REFERENCES bills(id),
    people_id BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL,
    category VARCHAR(50),
    payment_amount NUMERIC(12,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'completed',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL
);