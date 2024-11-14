package sheypoor

import (
	"fmt"
	"sync"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
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

func (s Crawler) CrawlArchivePage(link string, wg *sync.WaitGroup, timeoutCh <-chan time.Time, statusIsPicked bool) {
	fmt.Println("Crawling Sheypoor archive page:", link)
}

func (s Crawler) CrawlItemPage(link string) (structs.CrawledData, error) {
	fmt.Println("start ItemPage", link)

	time.Sleep(time.Second * 20)
	fmt.Println("Crawling Sheypoor item page:", link)
	return structs.CrawledData{}, nil
}
