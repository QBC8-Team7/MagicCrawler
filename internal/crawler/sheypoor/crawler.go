package sheypoor

import (
	"fmt"
	"sync"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
)

type SheypoorCrawler struct {
	Repository repositories.CrawlerRepository
}

func GetSourceName() string {
	return "sheypoor"
}

func (sc SheypoorCrawler) GetSourceName() string {
	return "sheypoor"
}

func (sc SheypoorCrawler) GetBaseUrl() string {
	return "https://sheypoor.com"
}

func (sc SheypoorCrawler) GetRepository() repositories.CrawlerRepository {
	return sc.Repository
}

func (sc SheypoorCrawler) CrawlItemPage(job sqlc.CrawlJob, wg *sync.WaitGroup) (structs.CrawledData, error) {
	fmt.Println("start ItemPage", job.Url)

	time.Sleep(time.Second * 20)
	fmt.Println("Crawling Sheypoor item page:", job.Url)
	return structs.CrawledData{}, nil
}

func (sc SheypoorCrawler) CreateCrawlJobArchivePageLink(link string) repositories.RepoResult {
	return repositories.RepoResult{}
}

func (sc SheypoorCrawler) GetSinglePageLinksFromArchivePage(htmlContent string) ([]string, error) {
	return []string{}, nil
}
