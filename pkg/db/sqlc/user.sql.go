// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO "user" (tg_id, role, watchlist_period)
VALUES ($1, $2, 0)
RETURNING tg_id, role, watchlist_period
`

type CreateAuthorParams struct {
	TgID string
	Role NullUserRole
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (User, error) {
	row := q.db.QueryRow(ctx, createAuthor, arg.TgID, arg.Role)
	var i User
	err := row.Scan(&i.TgID, &i.Role, &i.WatchlistPeriod)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM "user"
WHERE tg_id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, tgID string) error {
	_, err := q.db.Exec(ctx, deleteUser, tgID)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT tg_id, role, watchlist_period
FROM "user"
ORDER BY tg_id
LIMIT $1
OFFSET $2
`

type GetAllUsersParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetAllUsers(ctx context.Context, arg GetAllUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllUsers, arg.Limit, arg.Offset)
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

func (q *Queries) GetUserByTGID(ctx context.Context, tgID string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByTGID, tgID)
	var i User
	err := row.Scan(&i.TgID, &i.Role, &i.WatchlistPeriod)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE "user"
set role             = $2,
    watchlist_period = $3
WHERE tg_id = $1
RETURNING tg_id, role, watchlist_period
`

type UpdateUserParams struct {
	TgID            string
	Role            NullUserRole
	WatchlistPeriod pgtype.Int4
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser, arg.TgID, arg.Role, arg.WatchlistPeriod)
	return err
}