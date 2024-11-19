package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/helpers"
	"github.com/QBC8-Team7/MagicCrawler/internal/crawler/structs"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
)

type AdRepository struct {
	Queries *sqlc.Queries
	Logger  *logger.AppLogger
}

type AdQueryResult struct {
	Error error
	Exist bool
	AdId  int64
}

func (r AdRepository) CreateOrUpdateAd(crawledData structs.CrawledData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.FindAd(ctx, crawledData.PublisherAdKey)

	if result.Error != nil {
		return result.Error
	}

	if result.Exist {
		r.Logger.Infof(" | updating ad... | ad_id: %d | publisher_key_id: %s", result.AdId, crawledData.PublisherAdKey)
		updatingAdErr := r.UpdateAd(ctx, crawledData.PublisherAdKey, result.AdId, crawledData)
		if updatingAdErr != nil {
			return fmt.Errorf("error in updating ad: %s", updatingAdErr)
		}
		r.Logger.Infof(" | UPDATED | publisher_key_id: %s", crawledData.PublisherAdKey)
	} else {
		r.Logger.Infof(" | inserting ad... | ad_id: %d | publisher_key_id: %s", result.AdId, crawledData.PublisherAdKey)
		insertingAdErr := r.InsertAd(ctx, crawledData)
		if insertingAdErr != nil {
			return fmt.Errorf("error in inserting ad: %s", insertingAdErr)
		}
		r.Logger.Infof(" | INSERTED | publisher_key_id: %s", crawledData.PublisherAdKey)
	}

	return nil
}

func (r AdRepository) FindAd(ctx context.Context, publisherAdKey string) AdQueryResult {
	adID, err := r.Queries.GetAdByPublisherAdKey(ctx, publisherAdKey)
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
			AdId:  adID,
		}
	}
}

func (r AdRepository) InsertAd(ctx context.Context, crawledData structs.CrawledData) error {
	adParams, err := r.makeCreateAdParams(ctx, crawledData)
	if err != nil {
		return err
	}

	r.Logger.Infof(" | params prepared for inserting | publisher_key_id: %s", crawledData.PublisherAdKey)

	ad, err := r.Queries.CreateAd(ctx, adParams)
	if err != nil {
		return err
	}

	r.Logger.Infof(" | ad inserted | publisher_key_id: %s", crawledData.PublisherAdKey)

	_, err = r.InsertPrice(ctx, ad.ID, crawledData)
	if err != nil {
		return err
	}

	r.Logger.Infof(" | ad price inserted | publisher_key_id: %s", crawledData.PublisherAdKey)

	err = r.InsertPicture(ctx, ad.ID, crawledData)
	if err != nil {
		return err
	}

	r.Logger.Infof(" | ad image inserted | publisher_key_id: %s", crawledData.PublisherAdKey)

	return nil
}

func (r AdRepository) InsertPicture(ctx context.Context, adID int64, crawledData structs.CrawledData) error {
	if strings.TrimSpace(crawledData.ImageUrl) != "" {
		pictureParams := sqlc.CreateAdPictureParams{
			AdID: &adID,
			Url:  &crawledData.ImageUrl,
		}
		_, err := r.Queries.CreateAdPicture(ctx, pictureParams)
		return err
	}

	return nil
}

func (r AdRepository) UpdateAd(ctx context.Context, publisherAdKey string, adID int64, crawledData structs.CrawledData) error {
	adParams, err := r.makeUpdateAdParams(ctx, publisherAdKey, crawledData)
	if err != nil {
		return err
	}

	r.Logger.Infof(" | params prepared for updating | ad_id: %d | publisher_key_id: %s", adID, crawledData.PublisherAdKey)

	_, err = r.Queries.UpdateAd(ctx, adParams)
	if err != nil {
		return err
	}

	r.Logger.Infof(" | main properties updated | ad_id: %d | publisher_key_id: %s", adID, crawledData.PublisherAdKey)

	_, err = r.InsertPrice(ctx, adID, crawledData)
	if err != nil {
		return err
	}

	r.Logger.Infof(" | new price added | ad_id: %d | publisher_key_id: %s", adID, crawledData.PublisherAdKey)

	err = r.UpdatePicture(ctx, adID, crawledData)
	if err != nil {
		return err
	}

	r.Logger.Infof(" | ad picture updated | ad_id: %d | publisher_key_id: %s", adID, crawledData.PublisherAdKey)

	return nil
}

func (r AdRepository) FindPublisherId(ctx context.Context, publisherName string) (sqlc.Publisher, error) {
	return r.Queries.GetPublisherByName(ctx, &publisherName)
}

func (r AdRepository) makeCreateAdParams(ctx context.Context, crawledData structs.CrawledData) (sqlc.CreateAdParams, error) {
	publisher, err := r.FindPublisherId(ctx, crawledData.SourceName)
	if err != nil {
		return sqlc.CreateAdParams{}, err
	}

	int32Meterage := int32(crawledData.Meterage)

	int32RoomsCount := int32(crawledData.RoomsCount)

	year, _ := strconv.Atoi(crawledData.Year)
	int32Year := int32(year)

	int32FloorNumber := int32(crawledData.FloorNumber)

	latitude, _ := strconv.ParseFloat(crawledData.Lat, 64)
	longitude, _ := strconv.ParseFloat(crawledData.Lon, 64)

	params := sqlc.CreateAdParams{
		PublisherAdKey: crawledData.PublisherAdKey,
		PublisherID:    &publisher.ID,
		Category:       crawledData.AdCategory,
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
		HasParking:     &crawledData.HasParking,
	}

	if strings.TrimSpace(crawledData.Author) != "" {
		params.Author = &crawledData.Author
	}

	publishedAt, err := helpers.PersianToMiladi(crawledData.PublishedAt)
	if err == nil {
		params.PublishedAt = &publishedAt
	}

	if latitude != 0 && longitude != 0 {
		params.Lat = &latitude
		params.Lng = &longitude
	}

	int32TotalFloors := int32(crawledData.TotalFloors)
	if int32TotalFloors != 0 {
		params.TotalFloors = &int32TotalFloors
	}

	return params, nil
}

func (r AdRepository) makeUpdateAdParams(ctx context.Context, publisheAdKey string, crawledData structs.CrawledData) (sqlc.UpdateAdParams, error) {
	createAdParams, err := r.makeCreateAdParams(ctx, crawledData)
	if err != nil {
		return sqlc.UpdateAdParams{}, err
	}

	updateAdParams := sqlc.UpdateAdParams{
		PublisherAdKey: &publisheAdKey,
		PublisherID:    createAdParams.PublisherID,
		PublishedAt:    createAdParams.PublishedAt,
		Category:       createAdParams.Category,
		Author:         createAdParams.Author,
		Url:            createAdParams.Url,
		Title:          createAdParams.Title,
		Description:    createAdParams.Description,
		City:           createAdParams.City,
		Neighborhood:   createAdParams.Neighborhood,
		HouseType:      createAdParams.HouseType,
		Meterage:       createAdParams.Meterage,
		RoomsCount:     createAdParams.RoomsCount,
		Year:           createAdParams.Year,
		Floor:          createAdParams.Floor,
		TotalFloors:    createAdParams.TotalFloors,
		HasWarehouse:   createAdParams.HasWarehouse,
		HasElevator:    createAdParams.HasElevator,
		HasParking:     createAdParams.HasParking,
		Lat:            createAdParams.Lat,
		Lng:            createAdParams.Lng,
	}

	return updateAdParams, err
}

func (r AdRepository) UpdatePicture(ctx context.Context, adID int64, crawledData structs.CrawledData) error {
	// TODO - use transaction here
	err := r.Queries.DeleteAllPicturesOfAd(ctx, &adID)
	if err != nil {
		return err
	}

	err = r.InsertPicture(ctx, adID, crawledData)
	if err != nil {
		return err
	}

	return nil
}

func (r AdRepository) InsertPrice(ctx context.Context, adID int64, crawledData structs.CrawledData) (sqlc.Price, error) {
	return r.Queries.CreatePrice(ctx, r.makePriceParam(adID, crawledData))
}

func (r AdRepository) makePriceParam(adID int64, crawledData structs.CrawledData) sqlc.CreatePriceParams {
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

	params.HasPrice = &hasPrice

	return params
}
