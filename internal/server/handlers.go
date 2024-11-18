package server

import (
	"encoding/json"
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/pkg/watchlist"
	"net/http"
	"sort"
	"strconv"
	"time"

	myredis "github.com/QBC8-Team7/MagicCrawler/pkg/redis"

	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/labstack/echo/v4"
)

type jsonResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
}

type jsonListResponse struct {
	Success bool  `json:"success"`
	Message any   `json:"message"`
	Total   int64 `json:"total"`
}

func (s *Server) checkUserAccessToAd(userRole, userID string, adID int64) (bool, error) {
	// Check if they have an admin role
	if userRole == string(sqlc.UserRoleAdmin) || userRole == string(sqlc.UserRoleSuperAdmin) {
		return true, nil
	}

	userAds, err := s.db.GetUserAds(s.dbContext, userID)
	if err != nil {
		return false, fmt.Errorf("error checking user ads: %w", err)
	}

	// Check if the ad belongs to the user
	for _, id := range userAds {
		if id == adID {
			return true, nil
		}
	}

	// User neither owns the ad nor has privileges
	return false, nil
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: "ok",
	})
}

// Ad Group Handlers
func (s *Server) createAd(c echo.Context) error {
	type createAdWithPictureParam struct {
		sqlc.CreateAdParams
		PictureURL string `json:"pic_url"`
	}

	adParam := new(createAdWithPictureParam)
	if err := c.Bind(adParam); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid params",
		})
	}

	miniApID := int32(3)
	adParam.PublisherAdKey = "mini-app"
	adParam.PublisherID = &miniApID
	now := time.Now()
	adParam.PublishedAt = &now
	adParam.Url = nil

	ad, err := s.db.CreateAd(s.dbContext, adParam.CreateAdParams)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "invalid params",
		})
	}

	tgID, ok := c.Get("UserID").(string)
	if !ok || tgID == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	createUserAdParam := &sqlc.CreateUserAdParams{
		UserID: tgID,
		AdID:   ad.ID,
	}

	err = s.db.CreateUserAd(s.dbContext, *createUserAdParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal server while assigning new ad to user: %v", err),
		})
	}

	createPictureParam := sqlc.CreateAdPictureParams{
		AdID: &ad.ID,
		Url:  &adParam.PictureURL,
	}

	_, _ = s.db.CreateAdPicture(s.dbContext, createPictureParam)

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: ad.ID,
	})
}

func (s *Server) deleteAdByID(c echo.Context) error {
	adID, err := strconv.ParseInt(c.Param("adID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}

	userID, ok := c.Get("UserID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	userRole, ok := c.Get("UserRole").(string)
	if !ok || userRole == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user role is not set",
		})
	}

	hasAccess, err := s.checkUserAccessToAd(userRole, userID, adID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("error checking user access: %v", err),
		})
	}

	if !hasAccess {
		return c.JSON(http.StatusForbidden, jsonResponse{
			Success: false,
			Message: "you do not have permission or this ad does not exist",
		})
	}

	err = s.db.DeleteAd(s.dbContext, &adID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while deleting ad: %v", err),
		})
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: "ad deleted successfully",
	})
}

func (s *Server) getAdById(c echo.Context) error {
	adIDStr := c.Param("adID")
	adID, err := strconv.ParseInt(adIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}
	userID, ok := c.Get("UserID").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}
	params := sqlc.GetAdByIDParams{
		ID:     adID,
		UserID: userID,
	}

	ad, err := s.db.GetAdByID(s.dbContext, params)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonResponse{
			Success: false,
			Message: "ad not found",
		})
	}

	countStr, err := s.redis.Get(s.dbContext, myredis.CollectionAdCount, adIDStr)

	if err != nil {
		err := s.redis.Set(s.dbContext, myredis.CollectionAdCount, adIDStr, "1")
		if err != nil {
			s.logger.Warnf("can not set counter for ad %s: %v", adIDStr, err)
		} else {
			s.logger.Infof("counter set for Ad %s: 1", adIDStr)
		}
	} else {
		count, _ := strconv.Atoi(countStr)
		count = count + 1
		err := s.redis.Set(s.dbContext, myredis.CollectionAdCount, adIDStr, strconv.Itoa(count))
		if err != nil {
			s.logger.Warnf("can not increment counter for ad %s: %v", adIDStr, err)
		} else {
			s.logger.Infof("counter increment for Ad %s: %d", adIDStr, count)
		}
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: ad,
	})
}

func (s *Server) getAllAds(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")
	if limitStr == "" {
		limitStr = "10"
	}
	if offsetStr == "" {
		offsetStr = "0"
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid limit",
		})
	}
	offsetInt, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid offset",
		})
	}
	limit := int32(limitInt)
	offset := int32(offsetInt)

	getAllAdsParam := sqlc.GetAllAdsParams{
		Limit:  &limit,
		Offset: &offset,
	}

	ads, err := s.db.GetAllAds(s.dbContext, getAllAdsParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while getting all ads: %v", err),
		})
	}

	allAdsCount, err := s.db.CountAds(s.dbContext)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while counting all ads: %v", err),
		})
	}

	if len(ads) == 0 {
		ads = []sqlc.Ad{}
	}

	return c.JSON(http.StatusOK, jsonListResponse{
		Success: true,
		Message: ads,
		Total:   allAdsCount,
	})
}

func (s *Server) searchAds(c echo.Context) error {
	filterParam := new(sqlc.FilterAdsParams)

	if err := c.Bind(&filterParam); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("invalid params: %v", err),
		})
	}

	if filterParam.Author != nil {
		authorQuery := "%" + (*filterParam.Author) + "%"
		filterParam.Author = &authorQuery
	}

	if filterParam.Limit == nil {
		limit := int32(10)
		filterParam.Limit = &limit
	} else if *filterParam.Limit <= 0 {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid limit",
		})
	}

	if filterParam.Offset == nil {
		offset := int32(0)
		filterParam.Offset = &offset
	} else if *filterParam.Offset < 0 {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid offset",
		})
	}

	userLimit := *filterParam.Limit
	userOffset := *filterParam.Offset

	ads, err := s.db.FilterAds(s.dbContext, *filterParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("invalid params: %v", err),
		})
	}

	infiniteLimit, zeroOffset := int32(1000), int32(0)
	filterParam.Limit, filterParam.Offset = &infiniteLimit, &zeroOffset

	userId, ok := c.Get("UserID").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	filterBytes, err := json.Marshal(filterParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("can not marshal filter: %v", err),
		})
	}

	err = s.redis.Set(s.dbContext, myredis.CollectionFilter, userId, filterBytes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while writing filter in redis: %v", err),
		})
	}

	allDesiredAds, err := s.db.FilterAds(s.dbContext, *filterParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("error while searching ads: %v", err),
		})
	}
	total := int64(len(allDesiredAds))

	if len(ads) == 0 {
		ads = []sqlc.Ad{}
	}

	allAdIDs := make([]int64, len(allDesiredAds))
	for i, ad := range allDesiredAds {
		allAdIDs[i] = ad.ID
	}

	responseBytes, err := json.Marshal(allAdIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("can not marshal response: %v", err),
		})
	}

	err = s.redis.Set(s.dbContext, myredis.CollectionFilterResponse, userId, responseBytes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while writing response in redis: %v", err),
		})
	}

	var minPrice, maxPrice *int64
	if minPriceStr := c.QueryParam("minPrice"); minPriceStr != "" {
		price, err := strconv.ParseInt(minPriceStr, 10, 64)
		if err == nil {
			minPrice = &price
		}
	}
	if maxPriceStr := c.QueryParam("maxPrice"); maxPriceStr != "" {
		price, err := strconv.ParseInt(maxPriceStr, 10, 64)
		if err == nil {
			maxPrice = &price
		}
	}

	if minPrice == nil && maxPrice == nil {
		return c.JSON(http.StatusOK, jsonListResponse{
			Success: true,
			Message: ads,
			Total:   total,
		})
	}

	var category string
	if filterParam.Category != nil {
		category = *filterParam.Category
	}

	if category == "" {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "category must be specified when applying price filters",
		})
	}

	var filteredAds []sqlc.Ad

	switch category {
	case string(sqlc.AdCategoryBuy):
		filteredAds, err = s.db.FilterAdsPriceBuy(s.dbContext, sqlc.FilterAdsPriceBuyParams{
			AdIds:    allAdIDs,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
		})
	case string(sqlc.AdCategoryRent):
		filteredAds, err = s.db.FilterAdsPriceRent(s.dbContext, sqlc.FilterAdsPriceRentParams{
			AdIds:    allAdIDs,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
		})
	case string(sqlc.AdCategoryMortgage):
		filteredAds, err = s.db.FilterAdsPriceMortgage(s.dbContext, sqlc.FilterAdsPriceMortgageParams{
			AdIds:    allAdIDs,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
		})
	default:
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid category for price filtering",
		})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("error filtering ads by price: %v", err),
		})
	}

	if len(filteredAds) == 0 {
		filteredAds = []sqlc.Ad{}
	}

	total = int64(len(filteredAds))

	start := userOffset
	end := userOffset + userLimit

	if start > int32(len(filteredAds)) {
		filteredAds = []sqlc.Ad{}
	} else if end > int32(len(filteredAds)) {
		end = int32(len(filteredAds))
	}

	return c.JSON(http.StatusOK, jsonListResponse{
		Success: true,
		Message: filteredAds[start:end],
		Total:   total,
	})
}

func (s *Server) getPopularAds(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		if c.QueryParam("limit") != "" {
			return c.JSON(http.StatusBadRequest, jsonResponse{
				Success: false,
				Message: fmt.Sprintf("invalid limit: %v", err),
			})
		}
		limit = 5
	}

	adsPopularity, err := s.redis.GetAllFromCollection(s.dbContext, myredis.CollectionAdCount, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("error fetching ads: %v", err),
		})
	}

	type adPopularityPair struct {
		ID         int64
		Popularity int
	}

	var ads []adPopularityPair
	for idStr, popStr := range adsPopularity {
		id, err1 := strconv.ParseInt(idStr, 10, 64)
		pop, err2 := strconv.Atoi(popStr)
		if err1 == nil && err2 == nil {
			ads = append(ads, adPopularityPair{ID: id, Popularity: pop})
		}
	}

	sort.Slice(ads, func(i, j int) bool {
		return ads[i].Popularity > ads[j].Popularity
	})
	if len(ads) > limit {
		ads = ads[:limit]
	}

	adIDs := make([]int64, len(ads))
	for i, ad := range ads {
		adIDs[i] = ad.ID
	}

	dbAds, err := s.db.GetAdsByIds(s.dbContext, adIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("error fetching ads from database: %v", err),
		})
	}

	adMap := make(map[int64]sqlc.Ad, len(dbAds))
	for _, ad := range dbAds {
		adMap[ad.ID] = ad
	}

	popularAds := make([]sqlc.Ad, 0, len(ads))
	for _, ad := range ads {
		if dbAd, exists := adMap[ad.ID]; exists {
			popularAds = append(popularAds, dbAd)
		}
	}

	return c.JSON(http.StatusOK, jsonListResponse{
		Success: true,
		Message: popularAds,
		Total:   int64(limit),
	})
}

// Price Group Handlers
func (s *Server) setPriceOnAd(c echo.Context) error {
	priceParam := new(sqlc.CreatePriceParams)

	if err := c.Bind(priceParam); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("cannot parse the body into price: %v", err),
		})
	}

	userID, ok := c.Get("UserID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	userRole, ok := c.Get("UserRole").(string)
	if !ok || userRole == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user role is not set",
		})
	}

	hasAccess, err := s.checkUserAccessToAd(userRole, userID, priceParam.AdID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("error checking user access: %v", err),
		})
	}

	if !hasAccess {
		return c.JSON(http.StatusForbidden, jsonResponse{
			Success: false,
			Message: "you do not have permission or this ad does not exist",
		})
	}

	price, err := s.db.CreatePrice(s.dbContext, *priceParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while setting price on ad: %v", err),
		})
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: price,
	})
}

func (s *Server) getAdsLatestPrice(c echo.Context) error {
	adID, err := strconv.ParseInt(c.Param("adID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}

	latestPrice, err := s.db.GetLatestPriceByAdID(s.dbContext, adID)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonResponse{
			Success: false,
			Message: "no prices found or ad does not exist",
		})
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: latestPrice,
	})
}

func (s *Server) getAdsAllPrices(c echo.Context) error {
	adID, err := strconv.ParseInt(c.Param("adID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}

	prices, err := s.db.GetAllPricesByAdID(s.dbContext, adID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while getting all prices: %v", err),
		})
	}

	if len(prices) == 0 {
		prices = []sqlc.Price{}
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: prices,
	})
}

// Picture Groups Handlers
func (s *Server) getAdPictures(c echo.Context) error {
	adID, err := strconv.ParseInt(c.Param("adID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}

	userID, ok := c.Get("UserID").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}
	params := sqlc.GetAdByIDParams{
		ID:     adID,
		UserID: userID,
	}

	_, err = s.db.GetAdByID(s.dbContext, params)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonResponse{
			Success: false,
			Message: "ad not found",
		})
	}

	adPictures, err := s.db.GetAdPictures(s.dbContext, &adID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while getting ad adPictures: %v", err),
		})
	}

	if len(adPictures) == 0 {
		adPictures = []sqlc.AdPicture{}
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: adPictures,
	})
}

func (s *Server) deletePicture(c echo.Context) error {
	pictureID, err := strconv.ParseInt(c.Param("pictureID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid picture ID",
		})
	}

	_, err = s.db.GetPictureByID(s.dbContext, pictureID)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonResponse{
			Success: false,
			Message: "picture not found",
		})
	}

	err = s.db.DeletePictureByID(s.dbContext, pictureID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while deleting ad picture: %v", err),
		})
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: "ad picture deleted successfully",
	})
}

func (s *Server) createAdPicture(c echo.Context) error {
	createPicParam := new(sqlc.CreateAdPictureParams)

	if err := c.Bind(&createPicParam); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid params",
		})
	}
	userID, ok := c.Get("UserID").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}

	params := sqlc.GetAdByIDParams{
		ID:     *createPicParam.AdID,
		UserID: userID,
	}

	_, err := s.db.GetAdByID(s.dbContext, params)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonResponse{
			Success: false,
			Message: "ad not found",
		})
	}

	picture, err := s.db.CreateAdPicture(s.dbContext, *createPicParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while creating ad picture: %v", err),
		})
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: picture,
	})
}

// User Group Handlers
func (s *Server) getUsersAds(c echo.Context) error {
	userID, ok := c.Get("UserID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	adIDs, err := s.db.GetUserAds(s.dbContext, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while getting user ads: %v", err),
		})
	}

	ads, err := s.db.GetAdsByIds(s.dbContext, adIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while getting user ads: %v", err),
		})
	}

	if len(ads) == 0 {
		ads = []sqlc.Ad{}
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: ads,
	})
}

func (s *Server) createUserFavoriteAd(c echo.Context) error {
	adID, err := strconv.ParseInt(c.Param("adID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}

	userID, ok := c.Get("UserID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	userFavorites, err := s.db.GetUserFavoriteAds(s.dbContext, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while getting user favorites: %v", err),
		})
	}

	for _, favoriteID := range userFavorites {
		if favoriteID == adID {
			return c.JSON(http.StatusConflict, jsonResponse{
				Success: false,
				Message: "user favorite ad already exists",
			})
		}
	}

	params := sqlc.GetAdByIDParams{
		ID:     adID,
		UserID: userID,
	}

	_, err = s.db.GetAdByID(s.dbContext, params)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonResponse{
			Success: false,
			Message: "ad not found",
		})
	}

	createFavoriteParam := &sqlc.CreateUserFavoriteAdParams{
		UserID: userID,
		AdID:   adID,
	}

	err = s.db.CreateUserFavoriteAd(s.dbContext, *createFavoriteParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while creating user favorite ad: %v", err),
		})
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: "user favorite ad created",
	})
}

func (s *Server) getUserFavoriteAds(c echo.Context) error {
	userID, ok := c.Get("UserID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	adIDs, err := s.db.GetUserFavoriteAds(s.dbContext, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while getting user favorite ads: %v", err),
		})
	}

	ads, err := s.db.GetAdsByIds(s.dbContext, adIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while getting user ads: %v", err),
		})
	}

	if len(ads) == 0 {
		ads = []sqlc.Ad{}
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: ads,
	})
}

func (s *Server) deleteUserFavoriteAd(c echo.Context) error {
	adID, err := strconv.ParseInt(c.Param("adID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}

	userID, ok := c.Get("UserID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	userFavorites, err := s.db.GetUserFavoriteAds(s.dbContext, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while getting user favorites: %v", err),
		})
	}

	favoriteExists := false
	for _, favoriteID := range userFavorites {
		if favoriteID == adID {
			favoriteExists = true
		}
	}
	if !favoriteExists {
		return c.JSON(http.StatusNotFound, jsonResponse{
			Success: false,
			Message: "this ad is not one of user's favorite ads",
		})
	}

	deleteFavoriteParam := &sqlc.DeleteUserFavoriteAdParams{
		UserID: userID,
		AdID:   adID,
	}
	err = s.db.DeleteUserFavoriteAd(s.dbContext, *deleteFavoriteParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while deleting user favorite ad: %v", err),
		})
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: "user favorite ad deleted",
	})
}

func (s *Server) updateUserWatchListPeriod(c echo.Context) error {
	userID, ok := c.Get("UserID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	_, err := s.redis.Get(s.dbContext, myredis.CollectionFilter, userID)
	if err != nil {
		return c.JSON(http.StatusConflict, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("user must search and apply one filter first"),
		})
	}

	adParam := new(sqlc.UpdateUserParams)
	if err := c.Bind(adParam); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid params",
		})
	}
	adParam.TgID = userID

	if adParam.WatchlistPeriod == nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "WatchList is Required",
		})
	}
	_, err = s.db.UpdateUser(s.dbContext, *adParam)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonResponse{
			Success: false,
			Message: "user found",
		})
	}

	if *adParam.WatchlistPeriod == int32(0) {
		watchlist.GetService(s.dbContext, s.redis, s.db).StopWatch(userID)
	} else {
		watchlist.GetService(s.dbContext, s.redis, s.db).StartWatch(s.cfg.Bot.Token, s.logger, userID, int(*adParam.WatchlistPeriod))
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: "user watchlist updated",
	})
}

func (s *Server) getUserInfo(c echo.Context) error {
	userID, ok := c.Get("UserID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	userRole, ok := c.Get("UserRole").(string)
	if !ok || userRole == "" {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user role is not set",
		})
	}

	user, err := s.db.GetUserByTGID(s.dbContext, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while getting user info: %v", err),
		})
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: user,
	})
}
