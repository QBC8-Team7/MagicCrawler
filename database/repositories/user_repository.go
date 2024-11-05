package repositories

import (
	"database/sql"
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/database/models"
)

type UserRepository interface {
	GetUserByID(id int) (*models.User, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %w", err)
	}
	return user, nil
}
