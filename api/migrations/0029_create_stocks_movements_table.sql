CREATE TABLE IF NOT EXISTS stocks_movements (
    id BIGSERIAL PRIMARY KEY,

    item_id BIGINT NOT NULL REFERENCES items(id),
    store_id BIGINT NOT NULL REFERENCES stores(id),

    type VARCHAR(20) NOT NULL, -- IN, OUT, ADJUSTMENT
    category VARCHAR(50),      -- purchase, sale, damage, transfer

    quantity NUMERIC(12,2) NOT NULL,

    status VARCHAR(20) DEFAULT 'completed',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL
);