-- name: CreateTopUpOrder :one
INSERT INTO orders (user_id, number, amount, type)
VALUES ($1, $2, $3, 'CREDIT')
RETURNING *;

-- name: GetTopUpOrdersByUserID :many
SELECT *
FROM orders
WHERE type = 'CREDIT'
  AND user_id = $1;

-- name: GetOrderByNumber :one
SELECT *
FROM orders
WHERE number = $1
LIMIT 1;

-- name: UpdateOrderByNumber :exec
UPDATE orders
SET status       = $2,
    amount       = $3,
    processed_at = $4
WHERE number = $1;

-- name: GetNewTopUpOrder :one
SELECT *
FROM orders
WHERE type = 'CREDIT'
  AND status = 'NEW'
LIMIT 1;

-- name: Withdrawal :one
INSERT INTO orders (user_id, number, amount, type, status)
VALUES ($1, $2, $3, 'DEBIT', 'PROCESSED')
RETURNING *;

-- name: GetWithdrawalsByUserID :many
SELECT *
FROM orders
WHERE user_id = $1
  AND type = 'DEBIT'
  AND status = 'PROCESSED'
ORDER BY processed_at;

-- name: GetUserBalanceByUserID :one
SELECT COALESCE(SUM(
                        CASE type
                            WHEN 'CREDIT' THEN amount
                            WHEN 'DEBIT' THEN -amount
                            ELSE 0
                            END
                ), 0::NUMERIC(10, 2))::NUMERIC(10, 2)
FROM orders
WHERE user_id = $1
  AND status = 'PROCESSED';

-- name: GetUserWithdrawByUserID :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(10, 2)
FROM orders
WHERE user_id = $1
  AND type = 'DEBIT'
  AND status = 'PROCESSED';