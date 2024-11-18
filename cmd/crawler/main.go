package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/loggers"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	"github.com/QBC8-Team7/MagicCrawler/pkg/notification"
	"github.com/QBC8-Team7/MagicCrawler/pkg/utils"
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
	timeout := time.Duration(conf.Crawler.Time) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	metricLogger := logger.NewAppLogger(conf)
	metricLogger.InitLogger(conf.Crawler.MetricLogPath, conf.Logger.SysPath)

	mainLogger := logger.NewAppLogger(conf)
	mainLogger.InitCustomLogger(conf.Crawler.GeneralLogPath, conf.Logger.SysPath)

	repo := repositories.NewCrawlerRepository(queries, mainLogger)
	err = repo.ChangeWaitingOrPickedCrawlJobsStatusToFailed()
	if err != nil {
		mainLogger.Error("error in change old jobs status:", err)
		return
	}

	mainLogger.Info("unfinished jobs marked as failed")

	seeds := conf.Crawler.Seeds
	for _, seed := range seeds {
		repoResult := crawler.NewCrawler(seed["source"], repo, mainLogger).CreateCrawlJobArchivePageLink(seed["url"])
		if repoResult.Err != nil {
			mainLogger.Error("error in inserting seed to db", repoResult.Err)
			continue
		}

		if repoResult.Exist {
			mainLogger.Info("seed is exist in db")
			continue
		}

		mainLogger.Info("seed inserted to jobs table | ", seed["url"])
	}

	adminNotifier, err := notification.NewAdminNotifier(conf, queries)
	if err != nil {
		mainLogger.Error("failed to initialize notification service", err)
		return
	}

	jobRepository := repositories.JobRepository{Queries: queries}
	for {
		select {
		case <-ctx.Done():
			mainLogger.Info("Time For Crawl Has Finished")
			err = repo.ChangeWaitingOrPickedCrawlJobsStatusToFailed()
			if err != nil {
				mainLogger.Error("error in change old jobs status:", err)
			}
			return
		default:
			crawlJob, err := jobRepository.GetFirstWaitingCrawlJob()
			if err != nil {
				mainLogger.Info("no more jobs found: ", err)
				return
			}

			mainLogger.Infof("=> Working on %s page | JobID: %d | link: %s", crawlJob.PageType, crawlJob.ID, crawlJob.Url)
			crawlerInstance := crawler.NewCrawler(crawlJob.SourceName, repo, mainLogger)

			err, usage := utils.RunAndMeasureUsage(mainLogger, func() error {
				return crawler.Crawl(crawlerInstance, crawlJob)
			})

			loggers.MetricLog(*metricLogger, err, usage, crawlJob)

			if err != nil {
				mainLogger.Errorf(" | [FAILED] | error: %s | ID: %d | type: %s", err, crawlJob.ID, crawlJob.PageType)
				adminNotifier.Send(fmt.Sprintf("Crawling Failed\nJobID: %d\n%s", crawlJob.ID, err))
			} else {
				mainLogger.Infof(" | [DONE] | type: %s | ID: %d| link: %s", crawlJob.PageType, crawlJob.ID, crawlJob.Url)
			}
			mainLogger.Info(" |--------------------------------------------------------------------------------------------------------")
		}
	}
}
