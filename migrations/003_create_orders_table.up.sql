CREATE TABLE IF NOT EXISTS orders
(
    id           UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    amount       NUMERIC(10, 2) CHECK (amount IS NULL OR amount >= 0),
    number       TEXT        NOT NULL UNIQUE,
    type         OrderType   NOT NULL     DEFAULT 'CREDIT',
    status       OrderStatus NOT NULL     DEFAULT 'NEW',
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);
