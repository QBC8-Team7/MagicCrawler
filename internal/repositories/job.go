package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
)

const (
	CRAWLJOB_TYPE_ARCHIVE = "archive"
	CRAWLJOB_TYPE_SINGLE  = "single"

	CRAWLJOB_STATUS_WAITING = "waiting"
	CRAWLJOB_STATUS_DONE    = "done"
	CRAWLJOB_STATUS_FAILED  = "failed"
	CRAWLJOB_STATUS_PICKED  = "picked"
)

type JobRepository struct {
	Queries *sqlc.Queries
}

type RepoResult struct {
	Err   error
	Exist bool
	Job   sqlc.CrawlJob
}

func (r JobRepository) CreateCrawlJobForSinglePageLinks(links []string, sourceName string) []error {
	// TODO - maybe we need to use transaction here to make sure if all links inserted successfully
	var errors []error
	for _, link := range links {
		result := r.createCrawlJob(link, CRAWLJOB_TYPE_SINGLE, CRAWLJOB_STATUS_WAITING, sourceName)
		if result.Err != nil {
			errors = append(errors, result.Err)
		}
	}

	return errors
}

func (r JobRepository) CreateCrawlJobArchivePageLink(link string, sourceName string) RepoResult {
	return r.createCrawlJob(link, CRAWLJOB_TYPE_ARCHIVE, CRAWLJOB_STATUS_WAITING, sourceName)
}

func (r JobRepository) createCrawlJob(link string, pageType string, status string, sourceName string) RepoResult {
	fmt.Println(link)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checkCrawlJobExistsParams := sqlc.CheckCrawlJobExistsParams{
		Url:      link,
		Statuses: []string{CRAWLJOB_STATUS_WAITING, CRAWLJOB_STATUS_PICKED},
	}

	exists, err := r.Queries.CheckCrawlJobExists(ctx, checkCrawlJobExistsParams)
	if err != nil {
		return RepoResult{
			Err: err,
		}
	}

	err = nil

	if exists {
		return RepoResult{
			Err:   nil,
			Exist: true,
		}
	}

	createCrawlJobParams := sqlc.CreateCrawlJobParams{
		Url:        link,
		SourceName: sourceName,
		PageType:   pageType,
		Status:     status,
	}

	job, err := r.Queries.CreateCrawlJob(ctx, createCrawlJobParams)
	if err != nil {
		return RepoResult{
			Err: errors.New("failed to create crawl job: " + err.Error()),
		}
	}

	return RepoResult{
		Err:   nil,
		Exist: false,
		Job:   job,
	}
}

func (r JobRepository) UpdateCrawlJobStatus(jobID int64, newStatus string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := sqlc.UpdateCrawlJobStatusParams{
		Status: newStatus,
		JobID:  jobID,
	}

	return r.Queries.UpdateCrawlJobStatus(ctx, params)
}

func (r JobRepository) GetFirstWaitingCrawlJob() (sqlc.CrawlJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := CRAWLJOB_STATUS_WAITING
	return r.Queries.GetFirstCrawlJobByStatus(ctx, status)
}

func (r JobRepository) MakeOldCrawlJobsStatusFailed() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.Queries.ChangeAllCrawlJobsStatus(ctx, sqlc.ChangeAllCrawlJobsStatusParams{
		Statuses:  []string{CRAWLJOB_STATUS_WAITING, CRAWLJOB_STATUS_PICKED},
		NewStatus: CRAWLJOB_STATUS_FAILED,
	})
}