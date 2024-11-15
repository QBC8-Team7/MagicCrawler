package crawler

import (
	"fmt"
	"sync"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/divar"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/helpers"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/sheypoor"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
)

const ARCHIVE_PAGE = "archive"
const SINGLE_PAGE = "single"

type Crawler interface {
	CreateCrawlJobArchivePageLink(link string) repositories.RepoResult
	CrawlItemPage(job sqlc.CrawlJob, wg *sync.WaitGroup) (structs.CrawledData, error)
	GetSinglePageLinksFromArchivePage(htmlContent string) ([]string, error)
	GetBaseUrl() string
	GetSourceName() string
	GetRepository() repositories.CrawlJobRepository
}

func NewCrawler(sourceName string, repo repositories.CrawlJobRepository) Crawler {
	switch sourceName {
	case divar.GetSourceName():
		return divar.DivarCrawler{Repository: repo}
	case sheypoor.GetSourceName():
		return sheypoor.SheypoorCrawler{Repository: repo}
	default:
		panic("Unknown source, using default crawler")
	}
}

func CrawlArchivePage(crawler Crawler, job sqlc.CrawlJob, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := crawler.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_PICKED)
	if err != nil {
		fmt.Println(err)
		return
	}

	htmlContent, err := helpers.GetHtml(job.Url)
	if err != nil {
		// TODO - log here
		fmt.Println(err)
		// TODO - maybe we need to put error in db
		// TODO - maybe we need to save resource usage and time
		// TODO - maybe we can add try fields for job
		crawler.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
		return
	}

	links, err := crawler.GetSinglePageLinksFromArchivePage(htmlContent)
	if err != nil {
		fmt.Println(err)
		// TODO - maybe we need to put error in db
		// TODO - maybe we need to save resource usage and time
		// TODO - maybe we can add try fields for job
		crawler.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
		return
	}

	if len(links) > 0 {
		fmt.Println("links count:", len(links))
		nextLink, err := helpers.GetNextPageLink(job.Url)
		if err != nil {
			fmt.Println(err)
			// TODO - error handling
			// TODO - maybe we need to save resource usage and time
			// TODO - maybe we can add try fields for job
			crawler.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
			return
		}

		// TODO - maybe we need to use transactions to make sure all links with next link inserted successfuly together

		errors := crawler.GetRepository().CreateCrawlJobForSinglePageLinks(links, crawler.GetSourceName())
		if len(errors) > 0 {
			fmt.Println(errors[0])
			// TODO - error handling
			// TODO - maybe we need to save resource usage and time
			// TODO - maybe we can add try fields for job
			crawler.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
			return
		}

		nextLinkResult := crawler.GetRepository().CreateCrawlJobArchivePageLink(nextLink, crawler.GetSourceName())
		fmt.Println("next link:", nextLink)
		if nextLinkResult.Err != nil {
			// TODO - log here
			fmt.Println(nextLinkResult.Err)
			// TODO - maybe we need to put error in db
			// TODO - maybe we need to save resource usage and time
			// TODO - maybe we can add try fields for job
			crawler.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_FAILED)
			return
		}

		wg.Add(len(links) + 1)
	}

	crawler.GetRepository().UpdateCrawlJobStatus(job.ID, repositories.CRAWLJOB_STATUS_DONE)
}
