// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: picture.sql

package sqlc

import (
	"context"
)

const createAdPicture = `-- name: CreateAdPicture :one
INSERT INTO ad_picture(ad_id, url)
VALUES ($1, $2)
RETURNING id, ad_id, url
`

type CreateAdPictureParams struct {
	AdID *int64  `json:"ad_id"`
	Url  *string `json:"url"`
}

// Assign a picture to a Ad
func (q *Queries) CreateAdPicture(ctx context.Context, arg CreateAdPictureParams) (AdPicture, error) {
	row := q.db.QueryRow(ctx, createAdPicture, arg.AdID, arg.Url)
	var i AdPicture
	err := row.Scan(&i.ID, &i.AdID, &i.Url)
	return i, err
}

const deleteAllPicturesOfAd = `-- name: DeleteAllPicturesOfAd :exec
DELETE FROM ad_picture WHERE ad_id = $1
`

func (q *Queries) DeleteAllPicturesOfAd(ctx context.Context, adID *int64) error {
	_, err := q.db.Exec(ctx, deleteAllPicturesOfAd, adID)
	return err
}

const deletePictureByID = `-- name: DeletePictureByID :exec
DELETE FROM ad_picture
WHERE id = $1
`

func (q *Queries) DeletePictureByID(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deletePictureByID, id)
	return err
}

const getAdPictures = `-- name: GetAdPictures :many
SELECT id, ad_id, url FROM ad_picture
WHERE ad_id = $1
`

func (q *Queries) GetAdPictures(ctx context.Context, adID *int64) ([]AdPicture, error) {
	rows, err := q.db.Query(ctx, getAdPictures, adID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AdPicture
	for rows.Next() {
		var i AdPicture
		if err := rows.Scan(&i.ID, &i.AdID, &i.Url); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPictureByID = `-- name: GetPictureByID :one
SELECT id, ad_id, url
FROM ad_picture
WHERE id = $1
`

func (q *Queries) GetPictureByID(ctx context.Context, id int64) (AdPicture, error) {
	row := q.db.QueryRow(ctx, getPictureByID, id)
	var i AdPicture
	err := row.Scan(&i.ID, &i.AdID, &i.Url)
	return i, err
}
