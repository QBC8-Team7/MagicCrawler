package repositories

import (
	"database/sql"
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/internal/db/models"
)

type PublisherRepository interface {
	GetPublisherByName(name string) (*models.Publisher, error)
}

type publisherRepositoryImpl struct {
	db *sql.DB
}

func NewPublisherRepository(db *sql.DB) PublisherRepository {
	return &publisherRepositoryImpl{db: db}
}

func (p publisherRepositoryImpl) GetPublisherByName(name string) (*models.Publisher, error) {
	publisher := &models.Publisher{}
	query := `SELECT id, name, url FROM publisher WHERE name = $1`
	err := p.db.QueryRow(query, name).Scan(&publisher.ID, &publisher.Name, &publisher.URL)
	if err != nil {
		return nil, fmt.Errorf("could not find publisher: %w", err)
	}
	return publisher, nil
}
