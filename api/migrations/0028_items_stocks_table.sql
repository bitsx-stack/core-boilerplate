CREATE TABLE IF NOT EXISTS items_stocks (
    id BIGSERIAL PRIMARY KEY,

    item_id BIGINT NOT NULL REFERENCES items(id),
    store_id BIGINT NOT NULL REFERENCES stores(id),

    batch VARCHAR(50), 

    quantity NUMERIC(12,2) DEFAULT 0,

    selling_price NUMERIC(12,2) NOT NULL,
    discount_amount NUMERIC(12,2) DEFAULT 0,
    discount_rate NUMERIC(5,2) DEFAULT 0, -- percentage

    expiry_date DATE,

    UNIQUE(item_id, store_id, batch)
);