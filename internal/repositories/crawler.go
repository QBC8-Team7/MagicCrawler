package repositories

import "github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"

type CrawlerRepository struct {
	AdRepository
	JobRepository
}

func NewCrawlerRepository(queries *sqlc.Queries) CrawlerRepository {
	return CrawlerRepository{
		AdRepository: AdRepository{
			Queries: queries,
		},

		JobRepository: JobRepository{
			Queries: queries,
		},
	}
}
