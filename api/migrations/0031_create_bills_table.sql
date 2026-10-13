CREATE TABLE IF NOT EXISTS bills (
    id BIGSERIAL PRIMARY KEY,
    store_id BIGINT REFERENCES stores(id),

    code VARCHAR(50) UNIQUE,
    people_id BIGINT NOT NULL,
    
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',

    amount_billed NUMERIC(12,2) NOT NULL,
    amount_paid NUMERIC(12,2) DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL,

    deleted_at TIMESTAMP NULL,
    deleted_by BIGINT NULL
);