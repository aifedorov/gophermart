CREATE TABLE users
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

CREATE TYPE OrderStatus AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TYPE OrderType AS ENUM ('DEBIT', 'CREDIT');

CREATE TABLE IF NOT EXISTS orders
(
    id           UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    amount NUMERIC(10, 2) NOT NULL CHECK (amount >= 0),
    number       TEXT        NOT NULL UNIQUE,
    type         OrderType   NOT NULL     DEFAULT 'CREDIT',
    status       OrderStatus NOT NULL     DEFAULT 'NEW',
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_number ON orders (number);