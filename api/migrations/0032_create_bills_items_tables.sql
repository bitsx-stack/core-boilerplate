CREATE TABLE IF NOT EXISTS bill_items (
    id BIGSERIAL PRIMARY KEY,

    bill_id BIGINT NOT NULL REFERENCES bills(id),
    item_id BIGINT NOT NULL REFERENCES items(id),

    item_price NUMERIC(12,2) NOT NULL,
    item_quantity NUMERIC(12,2) NOT NULL,

    discount_amount NUMERIC(12,2) DEFAULT 0,
    discount_rate NUMERIC(5,2) DEFAULT 0
);