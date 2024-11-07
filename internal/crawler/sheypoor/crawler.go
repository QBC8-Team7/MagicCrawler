package sheypoor

import (
	"fmt"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
)

type Crawler struct{}

func (s Crawler) CrawlArchivePage(link string) []string {
	fmt.Println("Crawling Sheypoor archive page:", link)
	return []string{"item_link1", "item_link2"}
}

func (s Crawler) CrawlItemPage(link string) (structs.CrawledData, error) {
	fmt.Println("start ItemPage", link)

	time.Sleep(time.Second * 20)
	fmt.Println("Crawling Sheypoor item page:", link)
	return structs.CrawledData{}, nil
}
