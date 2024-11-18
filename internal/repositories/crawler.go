package repositories

import (
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
)

type CrawlerRepository struct {
	AdRepository
	JobRepository
}

func NewCrawlerRepository(queries *sqlc.Queries, logger *logger.AppLogger) CrawlerRepository {
	return CrawlerRepository{
		AdRepository: AdRepository{
			Queries: queries,
			Logger:  logger,
		},

		JobRepository: JobRepository{
			Queries: queries,
		},
	}
}
