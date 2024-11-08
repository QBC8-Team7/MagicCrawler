-- name: GetUserByTGID :one
SELECT *
FROM "user"
WHERE tg_id = $1
LIMIT 1;

-- name: GetAllUsers :many
SELECT *
FROM "user"
ORDER BY tg_id
LIMIT $1
OFFSET $2;

-- name: CreateAuthor :one
INSERT INTO "user" (tg_id, role, watchlist_period)
VALUES ($1, $2, 0)
RETURNING *;

-- name: UpdateUser :exec
UPDATE "user"
set role             = $2,
    watchlist_period = $3
WHERE tg_id = $1
RETURNING *;


-- name: DeleteUser :exec
DELETE
FROM "user"
WHERE tg_id = $1;