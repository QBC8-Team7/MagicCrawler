-- Get user by Telegram ID
-- name: GetUserByTGID :one
SELECT *
FROM "user"
WHERE tg_id = sqlc.narg('tg_id')
LIMIT 1;

-- Get all users with pagination
-- name: GetAllUsers :many
SELECT *
FROM "user"
ORDER BY tg_id
LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- Create a new user
-- name: CreateUser :one
INSERT INTO "user" (tg_id, role, watchlist_period)
VALUES (sqlc.arg('tg_id'), sqlc.arg('role'), sqlc.arg('watchlist_period'))
RETURNING *;

-- Update user role and watchlist period
-- name: UpdateUser :one
UPDATE "user"
SET role             = sqlc.narg('role'),
    watchlist_period = sqlc.narg('watchlist_period')
WHERE tg_id = sqlc.narg('tg_id')
RETURNING *;

-- Delete a user by Telegram ID
-- name: DeleteUser :exec
DELETE
FROM "user"
WHERE tg_id = sqlc.narg('tg_id');