// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package sqlc

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "user" (tg_id, role, watchlist_period)
VALUES ($1, $2, $3)
RETURNING tg_id, role, watchlist_period
`

type CreateUserParams struct {
	TgID            string       `json:"tg_id"`
	Role            NullUserRole `json:"role"`
	WatchlistPeriod *int32       `json:"watchlist_period"`
}

// Create a new user
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.TgID, arg.Role, arg.WatchlistPeriod)
	var i User
	err := row.Scan(&i.TgID, &i.Role, &i.WatchlistPeriod)
	return i, err
}

const createUserAd = `-- name: CreateUserAd :exec
INSERT INTO user_ads (user_id, ad_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`

type CreateUserAdParams struct {
	UserID *string `json:"user_id"`
	AdID   *int64  `json:"ad_id"`
}

// Assign an ad to a user as creator of that ad
func (q *Queries) CreateUserAd(ctx context.Context, arg CreateUserAdParams) error {
	_, err := q.db.Exec(ctx, createUserAd, arg.UserID, arg.AdID)
	return err
}

const createUserFavoriteAd = `-- name: CreateUserFavoriteAd :exec


INSERT INTO favorite_ads (user_id, ad_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`

type CreateUserFavoriteAdParams struct {
	UserID *string `json:"user_id"`
	AdID   *int64  `json:"ad_id"`
}

// Avoid duplicate entries
// Assign an ad to a user as creator of that ad
func (q *Queries) CreateUserFavoriteAd(ctx context.Context, arg CreateUserFavoriteAdParams) error {
	_, err := q.db.Exec(ctx, createUserFavoriteAd, arg.UserID, arg.AdID)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM "user"
WHERE tg_id = $1
`

// Delete a user by Telegram ID
func (q *Queries) DeleteUser(ctx context.Context, tgID string) error {
	_, err := q.db.Exec(ctx, deleteUser, tgID)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT tg_id, role, watchlist_period
FROM "user"
ORDER BY tg_id
LIMIT $2 OFFSET $1
`

type GetAllUsersParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

// Get all users with pagination
func (q *Queries) GetAllUsers(ctx context.Context, arg GetAllUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllUsers, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.TgID, &i.Role, &i.WatchlistPeriod); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByTGID = `-- name: GetUserByTGID :one
SELECT tg_id, role, watchlist_period
FROM "user"
WHERE tg_id = $1
LIMIT 1
`

// Get user by Telegram ID
func (q *Queries) GetUserByTGID(ctx context.Context, tgID string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByTGID, tgID)
	var i User
	err := row.Scan(&i.TgID, &i.Role, &i.WatchlistPeriod)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE "user"
SET role             = $1,
    watchlist_period = $2
WHERE tg_id = $3
RETURNING tg_id, role, watchlist_period
`

type UpdateUserParams struct {
	Role            NullUserRole `json:"role"`
	WatchlistPeriod *int32       `json:"watchlist_period"`
	TgID            string       `json:"tg_id"`
}

// Update user role and watchlist period
func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser, arg.Role, arg.WatchlistPeriod, arg.TgID)
	var i User
	err := row.Scan(&i.TgID, &i.Role, &i.WatchlistPeriod)
	return i, err
}
