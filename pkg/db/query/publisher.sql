-- Get publisher by its name
-- name: GetPublisherByName :one
SELECT *
FROM publisher
WHERE name = sqlc.narg('name')
LIMIT 1;

-- Insert a new publisher
-- name: CreatePublisher :one
INSERT INTO publisher (name, url)
VALUES (sqlc.arg('name'), sqlc.arg('url'))
RETURNING *;

-- Update publisher URL by ID
-- name: UpdatePublisherUrl :one
UPDATE publisher
SET url = sqlc.narg('url')
WHERE id = sqlc.narg('id')
RETURNING *;

-- Delete publisher by ID
-- name: DeletePublisher :exec
DELETE
FROM publisher
WHERE id = sqlc.narg('id');