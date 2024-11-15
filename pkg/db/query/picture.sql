-- Assign a picture to a Ad
-- name: CreateAdPicture :one
INSERT INTO ad_picture(ad_id, url)
VALUES (sqlc.arg('ad_id'), sqlc.arg('url'))
RETURNING *;

-- name: GetAdPictures :many
SELECT * FROM ad_picture
WHERE ad_id = sqlc.arg('ad_id');

-- name: DeletePictureByID :exec
DELETE FROM ad_picture
WHERE id = sqlc.arg('id');

-- name: GetPictureByID :one
SELECT *
FROM ad_picture
WHERE id = sqlc.arg('id');