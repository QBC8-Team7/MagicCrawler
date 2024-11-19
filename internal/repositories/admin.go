package repositories

import (
	"context"
	"errors"

	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
)

type AdminRepository struct {
	Queries         *sqlc.Queries
	LastAdminOffset int32
}

func NewAdminRepository(queries *sqlc.Queries) *AdminRepository {
	return &AdminRepository{
		Queries:         queries,
		LastAdminOffset: 0,
	}
}

func (r *AdminRepository) GetNextAdmin() (sqlc.User, error) {
	admin, err := r.Queries.GetNextAdmin(context.Background(), r.LastAdminOffset)
	if err != nil {
		if r.LastAdminOffset == 0 {
			return sqlc.User{}, errors.Join(errors.New("error in getting next admin"), err)
		} else {
			r.LastAdminOffset = 0
			return r.GetNextAdmin()
		}
	}
	r.LastAdminOffset = r.LastAdminOffset + 1
	return admin, nil
}