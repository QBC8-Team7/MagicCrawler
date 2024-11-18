package sheypoor

import (
	"fmt"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
)

type SheypoorCrawler struct {
	Repository repositories.CrawlerRepository
	Logger     *logger.AppLogger
}

func (c SheypoorCrawler) GetLogger() *logger.AppLogger {
	return c.Logger
}

func GetSourceName() string {
	return "sheypoor"
}

func (c SheypoorCrawler) GetSourceName() string {
	return "sheypoor"
}

func (c SheypoorCrawler) GetBaseUrl() string {
	return "https://sheypoor.com"
}

func (c SheypoorCrawler) GetRepository() repositories.CrawlerRepository {
	return c.Repository
}

func (c SheypoorCrawler) CrawlItemPage(job sqlc.CrawlJob) (structs.CrawledData, error) {
	fmt.Println("start ItemPage", job.Url)

	time.Sleep(time.Second * 20)
	fmt.Println("Crawling Sheypoor item page:", job.Url)
	return structs.CrawledData{}, nil
}

func (c SheypoorCrawler) CreateCrawlJobArchivePageLink(link string) repositories.RepoResult {
	return repositories.RepoResult{}
}

func (c SheypoorCrawler) GetSinglePageLinksFromArchivePage(htmlContent string) ([]string, error) {
	return []string{}, nil
}
