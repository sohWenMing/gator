-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * from users
WHERE name = $1
LIMIT 1;

-- name: GetAllUsers :many
SELECT id, created_at, updated_at, name from users;

-- name: DeleteUserByName :exec
DELETE FROM users
WHERE name = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;