package crawler

import (
	"errors"
	"fmt"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/divar"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/helpers"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/sheypoor"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
)

const ARCHIVE_PAGE = "archive"
const SINGLE_PAGE = "single"

type Crawler interface {
	CreateCrawlJobArchivePageLink(link string) repositories.RepoResult
	CrawlItemPage(job sqlc.CrawlJob) (structs.CrawledData, error)
	GetSinglePageLinksFromArchivePage(htmlContent string) ([]string, error)
	GetBaseUrl() string
	GetSourceName() string
	GetRepository() repositories.CrawlerRepository
	GetLogger() *logger.AppLogger
}

func NewCrawler(sourceName string, repo repositories.CrawlerRepository, logger *logger.AppLogger) Crawler {
	switch sourceName {
	case divar.GetSourceName():
		return divar.DivarCrawler{Repository: repo, Logger: logger}
	case sheypoor.GetSourceName():
		return sheypoor.SheypoorCrawler{Repository: repo, Logger: logger}
	default:
		panic("Unknown source, using default crawler")
	}
}

func Crawl(crawlerInstance Crawler, crawlJob sqlc.CrawlJob) error {
	_, err := crawlerInstance.GetRepository().UpdateCrawlJobStatus(crawlJob.ID, repositories.CRAWLJOB_STATUS_PICKED)
	if err != nil {
		return err
	}
	crawlerInstance.GetLogger().Infof(" | job status change to picked | JobID: %d", crawlJob.ID)

	if crawlJob.PageType == ARCHIVE_PAGE {
		return CrawlArchivePage(crawlerInstance, crawlJob)
	} else {
		crawledData, crawlItemPageErr := crawlerInstance.CrawlItemPage(crawlJob)
		if crawlItemPageErr != nil {
			_, updateJobStatusErr := crawlerInstance.GetRepository().UpdateCrawlJobStatus(crawlJob.ID, repositories.CRAWLJOB_STATUS_FAILED)
			if updateJobStatusErr != nil {
				return errors.Join(crawlItemPageErr, updateJobStatusErr)
			}
			return fmt.Errorf("error in crawling single page. jobID: %d - %s", crawlJob.ID, crawlItemPageErr)
		} else {
			createOrUpdateAdErr := crawlerInstance.GetRepository().CreateOrUpdateAd(crawledData)
			if createOrUpdateAdErr != nil {
				_, updateJobStatusErr := crawlerInstance.GetRepository().UpdateCrawlJobStatus(crawlJob.ID, repositories.CRAWLJOB_STATUS_FAILED)
				if updateJobStatusErr != nil {
					return errors.Join(createOrUpdateAdErr, updateJobStatusErr)
				}

				return createOrUpdateAdErr
			} else {
				_, updateJobStatusErr := crawlerInstance.GetRepository().UpdateCrawlJobStatus(crawlJob.ID, repositories.CRAWLJOB_STATUS_DONE)
				if updateJobStatusErr != nil {
					return fmt.Errorf("error in changing job status: %s", updateJobStatusErr)
				}
				crawlerInstance.GetLogger().Infof(" | job status change to DONE | JobID: %d", crawlJob.ID)
			}
		}
	}

	return nil
}

func CrawlArchivePage(crawlerInstance Crawler, job sqlc.CrawlJob) error {
	htmlContent, getHtmlError := helpers.GetHtml(job.Url)
	if getHtmlError != nil {
		_, updateJobStatusErr := crawlerInstance.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
		if updateJobStatusErr != nil {
			return errors.Join(getHtmlError, updateJobStatusErr)
		}
		return fmt.Errorf("error in getting html: %s", getHtmlError)
	}

	crawlerInstance.GetLogger().Infof(" | got html page | JobID: %d", job.ID)

	links, getSingleLinksErr := crawlerInstance.GetSinglePageLinksFromArchivePage(htmlContent)
	if getSingleLinksErr != nil {
		_, updateJobStatusErr := crawlerInstance.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
		if updateJobStatusErr != nil {
			return errors.Join(getSingleLinksErr, updateJobStatusErr)
		}
		return fmt.Errorf("error in getting single page links in archive page: %s", getSingleLinksErr)
	}

	if len(links) > 0 {
		crawlerInstance.GetLogger().Infof(" | %d single pages links extracted from archive page | JobID: %d", len(links), job.ID)

		nextLink, MakeNextPageLinkErr := helpers.GetNextPageLink(job.Url)
		if MakeNextPageLinkErr != nil {
			_, updateJobStatusErr := crawlerInstance.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
			if updateJobStatusErr != nil {
				return errors.Join(MakeNextPageLinkErr, updateJobStatusErr)

			}
			return fmt.Errorf("error in getting next page of archive page: %s", MakeNextPageLinkErr)
		}

		allErrors := crawlerInstance.GetRepository().CreateCrawlJobForSinglePageLinks(links, crawlerInstance.GetSourceName())
		if allErrors != nil {
			_, updateJobStatusErr := crawlerInstance.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
			if updateJobStatusErr != nil {
				return errors.Join(allErrors, updateJobStatusErr)

			}
			return errors.Join(errors.New("error in making crawl jobs records for extracted single page links"), allErrors)
		}

		crawlerInstance.GetLogger().Infof(" | single pages links inserted to jobs table | JobID: %d", job.ID)

		nextLinkResult := crawlerInstance.GetRepository().CreateCrawlJobArchivePageLink(nextLink, crawlerInstance.GetSourceName())
		if nextLinkResult.Err != nil {
			_, updateJobStatusErr := crawlerInstance.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
			if updateJobStatusErr != nil {
				return errors.Join(nextLinkResult.Err, updateJobStatusErr)
			}
			return fmt.Errorf("error in creating crawl job for next page of archive page: %s", nextLinkResult.Err)
		}

		crawlerInstance.GetLogger().Infof(" | next archive page link inserted to jobs table | JobID: %d", job.ID)
	} else {
		crawlerInstance.GetLogger().Infof(" | no links found in archive page | JobID: %d", job.ID)
	}

	_, updateJobStatusErr := crawlerInstance.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_DONE)
	if updateJobStatusErr != nil {
		return errors.Join(errors.New("error in updating archive page status"), updateJobStatusErr)
	}

	crawlerInstance.GetLogger().Infof(" | archive page job status changed to DONE | JobID: %d", job.ID)

	return nil
}
