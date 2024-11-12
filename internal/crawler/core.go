package crawler

import (
	"fmt"
	"sync"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/divar"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/sheypoor"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
)

type Crawler interface {
	CrawlArchivePage(link string, wg *sync.WaitGroup)
	CrawlItemPage(link string) (structs.CrawledData, error)
}

const (
	SOURCE_DIVAR    = "divar"
	SOURCE_SHEYPOOR = "sheypoor"
)

func Start(seeds []map[string]string, timeout time.Duration) {
	done := make(chan struct{})
	timeoutCh := time.After(timeout)

	var wg sync.WaitGroup

	for _, seed := range seeds {
		wg.Add(1)
		go func(link, source string) {
			defer wg.Done()
			crawler := getCrawler(source)
			crawler.CrawlArchivePage(link, &wg)
		}(seed["link"], seed["source"])
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		fmt.Println("All goroutines finished")
	case <-timeoutCh:
		fmt.Println("Time is over")
	}
}

func getCrawler(source string) Crawler {
	switch source {
	case SOURCE_DIVAR:
		return divar.Crawler{}
	case SOURCE_SHEYPOOR:
		return sheypoor.Crawler{}
	default:
		panic("Unknown source, using default crawler")
	}
}
