package sheypoor

import (
	"fmt"
	"sync"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
)

type Crawler struct{}

func (s Crawler) CrawlArchivePage(link string, wg *sync.WaitGroup) {
	fmt.Println("Crawling Sheypoor archive page:", link)
}

func (s Crawler) CrawlItemPage(link string) (structs.CrawledData, error) {
	fmt.Println("start ItemPage", link)

	time.Sleep(time.Second * 20)
	fmt.Println("Crawling Sheypoor item page:", link)
	return structs.CrawledData{}, nil
}
