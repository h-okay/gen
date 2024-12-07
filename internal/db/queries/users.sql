-- name: CreateUser :one
INSERT INTO users (email, name)
VALUES (@email, @name)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = @id;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;