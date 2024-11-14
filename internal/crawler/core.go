package crawler

import (
	"sync"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/divar"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/sheypoor"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
)

type Crawler interface {
	CrawlArchivePage(link string, wg *sync.WaitGroup, timeoutCh <-chan time.Time, statusIsPicked bool)
	CrawlItemPage(link string) (structs.CrawledData, error)
	GetBaseUrl() string
}

func NewCrawler(sourceName string, repo repositories.CrawlJobRepository) Crawler {
	switch sourceName {
	case divar.GetSourceName():
		return divar.Crawler{Repository: repo}
	case sheypoor.GetSourceName():
		return sheypoor.Crawler{Repository: repo}
	default:
		panic("Unknown source, using default crawler")
	}
}
