package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/helpers"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
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

type CrawlJobRepository struct {
	Queries *sqlc.Queries
}

type RepoResult struct {
	Err   error
	Exist bool
	Job   sqlc.CrawlJob
}

func (repo CrawlJobRepository) CreateCrawlJobForSinglePageLinks(links []string, sourceName string) []error {
	// TODO - maybe we need to use transaction here to make sure if all links inserted successfully
	var errors []error
	for _, link := range links {
		result := repo.createCrawlJob(link, CRAWLJOB_TYPE_SINGLE, CRAWLJOB_STATUS_WAITING, sourceName)
		if result.Err != nil {
			errors = append(errors, result.Err)
		}
	}

	return errors
}

func (repo CrawlJobRepository) CreateCrawlJobArchivePageLink(link string, sourceName string) RepoResult {
	return repo.createCrawlJob(link, CRAWLJOB_TYPE_ARCHIVE, CRAWLJOB_STATUS_WAITING, sourceName)
}

func (repo CrawlJobRepository) createCrawlJob(link string, pageType string, status string, sourceName string) RepoResult {
	fmt.Println(link)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checkCrawlJobExistsParams := sqlc.CheckCrawlJobExistsParams{
		Url:      link,
		Statuses: []string{CRAWLJOB_STATUS_WAITING, CRAWLJOB_STATUS_PICKED},
	}

	exists, err := repo.Queries.CheckCrawlJobExists(ctx, checkCrawlJobExistsParams)
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

	job, err := repo.Queries.CreateCrawlJob(ctx, createCrawlJobParams)
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

func (cjr CrawlJobRepository) UpdateCrawlJobStatus(jobID int64, newStatus string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := sqlc.UpdateCrawlJobStatusParams{
		Status: newStatus,
		JobID:  jobID,
	}

	return cjr.Queries.UpdateCrawlJobStatus(ctx, params)
}

func (cjr CrawlJobRepository) GetFirstWaitingCrawlJob() (sqlc.CrawlJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := CRAWLJOB_STATUS_WAITING
	return cjr.Queries.GetFirstCrawlJobByStatus(ctx, status)
}

// TODO - make a separate repository
type AdQueryResult struct {
	Error error
	Exist bool
	AdId  int64
}

func (cjr CrawlJobRepository) CreateOrUpdateAd(crawledData structs.CrawledData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := cjr.FindAd(ctx, crawledData.AdId)

	if result.Error != nil {
		return result.Error
	}

	if result.Exist {
		cjr.UpdateAd(ctx, result.AdId, crawledData)
	} else {
		cjr.InsertAd(ctx, crawledData)
	}

	return nil
}

func (cjr CrawlJobRepository) FindAd(ctx context.Context, publisherAdKey string) AdQueryResult {
	adId, err := cjr.Queries.GetAdByPublisherAdKey(ctx, publisherAdKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AdQueryResult{
				Error: nil,
				Exist: false,
			}
		} else {
			return AdQueryResult{
				Error: err,
			}
		}
	} else {
		return AdQueryResult{
			Error: nil,
			Exist: true,
			AdId:  adId,
		}
	}
}

func (cjr CrawlJobRepository) InsertAd(ctx context.Context, crawledData structs.CrawledData) error {
	adParams, err := cjr.MakeCreateAdParams(ctx, crawledData)
	if err != nil {
		return err
	}

	// TODO - use transaction for these

	ad, err := cjr.Queries.CreateAd(ctx, adParams)
	if err != nil {
		return err
	}

	priceParams := cjr.MakePriceParam(ad.ID, crawledData)
	cjr.Queries.CreatePrice(ctx, priceParams)

	pictureParams := sqlc.CreateAdPictureParams{
		AdID: &ad.ID,
		Url:  &crawledData.ImageUrl,
	}
	cjr.Queries.CreateAdPicture(ctx, pictureParams)

	return nil
}

func (cjr CrawlJobRepository) UpdateAd(ctx context.Context, adId int64, crawledData structs.CrawledData) {
	// UPDATE PRICE
	// UPDATE IMAGE
}

func (cjr CrawlJobRepository) FindPublisherId(ctx context.Context, publisherName string) (sqlc.Publisher, error) {
	return cjr.Queries.GetPublisherByName(ctx, &publisherName)
}

func (cjr CrawlJobRepository) MakeCreateAdParams(ctx context.Context, crawledData structs.CrawledData) (sqlc.CreateAdParams, error) {
	publisher, err := cjr.FindPublisherId(ctx, crawledData.SourceName)
	if err != nil {
		return sqlc.CreateAdParams{}, err
	}

	publishedAt, err := helpers.PersianToMiladi(crawledData.PublishedAt)
	if err != nil {
		return sqlc.CreateAdParams{}, err
	}

	// TODO - get author
	author := ""

	int32Meterage := int32(crawledData.Meterage)

	int32RoomsCount := int32(crawledData.RoomsCount)

	year, _ := strconv.Atoi(crawledData.Year)
	int32Year := int32(year)

	int32FloorNumber := int32(crawledData.FloorNumber)

	latitude, _ := strconv.ParseFloat(crawledData.Lat, 64)
	longitude, _ := strconv.ParseFloat(crawledData.Lon, 64)

	return sqlc.CreateAdParams{
		PublisherAdKey: crawledData.AdId,
		PublisherID:    &publisher.ID,
		PublishedAt:    publishedAt,
		Category:       crawledData.AdCategory,
		Author:         &author,
		Url:            &crawledData.URL,
		Title:          &crawledData.Title,
		Description:    &crawledData.Description,
		City:           &crawledData.City,
		Neighborhood:   &crawledData.Neighborhood,
		HouseType:      crawledData.HouseType,
		Meterage:       &int32Meterage,
		RoomsCount:     &int32RoomsCount,
		Year:           &int32Year,
		Floor:          &int32FloorNumber,
		TotalFloors:    nil,
		HasWarehouse:   &crawledData.HasWarehouse,
		HasElevator:    &crawledData.HasElevator,
		HasParking:     nil,
		Lat:            &latitude,
		Lng:            &longitude,
	}, nil
}

func (cjr CrawlJobRepository) MakePriceParam(adID int64, crawledData structs.CrawledData) sqlc.CreatePriceParams {
	hasPrice := false

	params := sqlc.CreatePriceParams{
		AdID: adID,
	}

	totalPirce, err := strconv.Atoi(crawledData.TotalPrice)
	if err != nil {
		params.TotalPrice = nil
	} else {
		int64totalPrice := int64(totalPirce)
		params.TotalPrice = &int64totalPrice
		hasPrice = true
	}

	pricePerMeter, err := strconv.Atoi(crawledData.PricePerMeter)
	if err != nil {
		params.PricePerMeter = nil
	} else {
		int64pricePerMeter := int64(pricePerMeter)
		params.PricePerMeter = &int64pricePerMeter
		hasPrice = true
	}

	prePaidPrice, err := strconv.Atoi(crawledData.PrePaidPrice)
	if err != nil {
		params.Mortgage = nil
	} else {
		int64prePaidPrice := int64(prePaidPrice)
		params.Mortgage = &int64prePaidPrice
		hasPrice = true
	}

	monthlyRentPrice, err := strconv.Atoi(crawledData.MonthlyRentPrice)
	if err != nil {
		params.NormalPrice = nil
	} else {
		int64monthlyRentPrice := int64(monthlyRentPrice)
		params.NormalPrice = &int64monthlyRentPrice
		hasPrice = true
	}

	hasPrice = false

	params.HasPrice = &hasPrice

	return params
}
