-- Get user by Telegram ID
-- name: GetUserByTGID :one
SELECT *
FROM "user"
WHERE tg_id = sqlc.arg('tg_id')
LIMIT 1;

-- Get all users with pagination
-- name: GetAllUsers :many
SELECT *
FROM "user"
ORDER BY tg_id
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

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
WHERE tg_id = sqlc.arg('tg_id')
RETURNING *;

-- Delete a user by Telegram ID
-- name: DeleteUser :exec
DELETE
FROM "user"
WHERE tg_id = sqlc.arg('tg_id');


-- Assign an ad to a user as creator of that ad
-- name: CreateUserAd :exec
INSERT INTO user_ads (user_id, ad_id)
VALUES (sqlc.arg('user_id'), sqlc.arg('ad_id'))
ON CONFLICT DO NOTHING;
-- Avoid duplicate entries


-- Assign an ad to a user as creator of that ad
-- name: CreateUserFavoriteAd :exec
INSERT INTO favorite_ads (user_id, ad_id)
VALUES (sqlc.arg('user_id'), sqlc.arg('ad_id'))
ON CONFLICT DO NOTHING; -- Avoid duplicate entries
