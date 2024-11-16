package main

import (
	"context"
	"flag"
	"fmt"
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
	fmt.Println("Start")

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
		// {"link": "https://divar.ir/s/tehran-province/buy-apartment", "source": divar.GetSourceName()},
		{"link": "https://divar.ir/s/tehran-province/buy-villa", "source": divar.GetSourceName()},
		{"link": "https://divar.ir/s/tehran-province/rent-apartment", "source": divar.GetSourceName()},
		{"link": "https://divar.ir/s/tehran-province/rent-villa", "source": divar.GetSourceName()},
	}

	// Set crawl duration to 10 minutes
	// TODO - use context if you can do
	timeout := time.Duration(100) * time.Second
	timeoutCh := time.After(timeout)

	// TODO - maybe you need to use buffered channel
	done := make(chan struct{})
	var crawlerVar crawler.Crawler
	var wg sync.WaitGroup

	jobRepository := repositories.JobRepository{
		Queries: queries,
	}

	repo := repositories.NewCrawlerRepository(queries)
	err = repo.MakeOldCrawlJobsStatusFailed()
	if err != nil {
		fmt.Println("×××××× ERROR in change old jobs status:", err)
		return
	}

	fmt.Println("old jobs status changed to failed")

	for index, seed := range seeds {
		fmt.Println("Seed:", index, seed["link"])
		crawlerVar = crawler.NewCrawler(seed["source"], repo)
		repoResult := crawlerVar.CreateCrawlJobArchivePageLink(seed["link"])
		if repoResult.Err != nil || repoResult.Exist {
			continue
		}
		fmt.Println("Add seed to jobs:", seed["link"])

		wg.Add(1)
	}

	go func() {
		// TODO - implement worker pool here
		for {
			fmt.Println("begin of iteration")
			crawlJob, err := jobRepository.GetFirstWaitingCrawlJob()
			if err != nil {
				fmt.Println("×××××× Error in getting job", err)
				return
			}
			fmt.Println("Job:", crawlJob.PageType, crawlJob.Url)

			workerCrawler := crawler.NewCrawler(crawlJob.SourceName, repo)

			if crawlJob.PageType == crawler.ARCHIVE_PAGE {
				crawler.CrawlArchivePage(workerCrawler, crawlJob, &wg)
			} else {
				crawledData, err := workerCrawler.CrawlItemPage(crawlJob, &wg)
				if err != nil {
					fmt.Println("×××××× Error in crawling single page. jobID:", crawlJob.ID, err)
					workerCrawler.GetRepository().UpdateCrawlJobStatus(crawlJob.ID, repositories.CRAWLJOB_STATUS_FAILED)
				} else {
					err = workerCrawler.GetRepository().CreateOrUpdateAd(crawledData)
					if err != nil {
						fmt.Println("×××××× Error in insert or update:", err)
						_, err = workerCrawler.GetRepository().UpdateCrawlJobStatus(crawlJob.ID, repositories.CRAWLJOB_STATUS_FAILED)
						if err != nil {
							fmt.Println("×××××× Error in changing job status", err)
						}
					} else {
						_, err = workerCrawler.GetRepository().UpdateCrawlJobStatus(crawlJob.ID, repositories.CRAWLJOB_STATUS_DONE)
						if err != nil {
							fmt.Println("×××××× Error in changing job status", err)
						}
					}

				}

			}
		}
	}()

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return
	case <-timeoutCh:
		fmt.Println("TIME FINISHED")
		err := repo.MakeOldCrawlJobsStatusFailed()
		if err != nil {
			fmt.Println(err)
		}
		return
	}

}
