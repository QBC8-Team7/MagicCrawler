package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/divar"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
)

func main() {
	configPath := flag.String("c", "config.yml", "Path to the configuration file")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}

	dbContext := context.Background()

	dbUri := db.GetDbUri(conf)
	dbConn, err := db.GetDBConnection(dbContext, dbUri)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	defer func() {
		if err := dbConn.Close(dbContext); err != nil {
			log.Fatalf("could not close connection with database: %v", err)
		}
	}()

	queries := sqlc.New(dbConn)

	seeds := []map[string]string{
		{"link": "https://divar.ir/s/tehran-province/buy-apartment", "source": divar.GetSourceName()},
		// {"link": "https://divar.ir/s/tehran-province/buy-villa", "source": divar.GetSourceName()},
		// {"link": "https://divar.ir/s/tehran-province/rent-apartment", "source": divar.GetSourceName()},
		// {"link": "https://divar.ir/s/tehran-province/rent-villa", "source": divar.GetSourceName()},
	}

	// Set crawl duration to 10 minutes
	// TODO - use context if you can
	timeout := time.Duration(10) * time.Minute
	timeoutCh := time.After(timeout)

	// TODO - better to use buffered channel
	done := make(chan struct{})
	var crawlerVar crawler.Crawler
	var wg sync.WaitGroup

	crawlJobRepository := repositories.CrawlJobRepository{
		Queries: queries,
	}

	for _, seed := range seeds {
		crawlerVar = crawler.NewCrawler(seed["source"], crawlJobRepository)
		wg.Add(1)
		go crawlerVar.CrawlArchivePage(seed["link"], &wg, timeoutCh, true)
		time.Sleep(time.Millisecond * 500)
	}

	// RUN WORKER POOL HERE

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return
	case <-timeoutCh:
		return
	}

}
