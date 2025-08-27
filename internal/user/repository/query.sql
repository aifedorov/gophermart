-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByCredentials :one
SELECT *
FROM users
WHERE username = $1
  AND password_hash = $2;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (username, password_hash)
VALUES ($1, $2)
RETURNING *;