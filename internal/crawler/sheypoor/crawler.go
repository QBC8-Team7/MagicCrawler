package sheypoor

import (
	"fmt"
	"sync"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
)

type Crawler struct {
	Repository repositories.CrawlJobRepository
}

func GetSourceName() string {
	return "sheypoor"
}

func (c Crawler) GetBaseUrl() string {
	return "https://sheypoor.com"
}

func (s Crawler) CrawlArchivePage(job sqlc.CrawlJob, wg *sync.WaitGroup, timeoutCh <-chan time.Time) {
	fmt.Println("Crawling Sheypoor archive page:", job.Url)
}

func (s Crawler) CrawlItemPage(job sqlc.CrawlJob, wg *sync.WaitGroup, timeoutCh <-chan time.Time) (structs.CrawledData, error) {
	fmt.Println("start ItemPage", job.Url)

	time.Sleep(time.Second * 20)
	fmt.Println("Crawling Sheypoor item page:", job.Url)
	return structs.CrawledData{}, nil
}

func (s Crawler) CreateCrawlJobArchivePageLink(link string) repositories.RepoResult {
	return repositories.RepoResult{}
}
