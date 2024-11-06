package repositories

import (
	"database/sql"
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/models"
)

type UserRepository interface {
	GetUserByID(tgID string) (*models.User, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) GetUserByID(tgID string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT tg_id FROM "user" WHERE tg_id = $1`
	err := r.db.QueryRow(query, tgID).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %w", err)
	}
	return user, nil
}
