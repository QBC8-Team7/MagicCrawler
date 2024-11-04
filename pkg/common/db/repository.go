package db

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPublisherByName(name string) (*Publisher, error) {
	publisher := &Publisher{}
	query := `SELECT id, name, url FROM publisher WHERE name = $1`
	err := r.db.QueryRow(query, name).Scan(&publisher.ID, &publisher.Name, &publisher.URL)
	if err != nil {
		return nil, fmt.Errorf("could not find publisher: %w", err)
	}
	return publisher, nil
}
