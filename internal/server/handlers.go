package server

import (
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type jsonResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
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
	adParam := new(sqlc.CreateAdParams)
	if err := c.Bind(adParam); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("invalid params"),
		})
	}

	adParam.PublisherAdKey = "mini-app"
	adParam.PublisherID = nil
	adParam.PublishedAt = time.Now()
	adParam.Url = nil

	ad, err := s.db.CreateAd(s.dbContext, *adParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("invalid params"),
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
	adID, err := strconv.ParseInt(c.Param("adID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: "invalid ad ID",
		})
	}

	ad, err := s.db.GetAdByID(s.dbContext, adID)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonResponse{
			Success: false,
			Message: "ad not found",
		})
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

	if len(ads) == 0 {
		ads = []sqlc.Ad{}
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: ads,
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

	_, err = s.db.GetAdByID(s.dbContext, adID)
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

	_, err := s.db.GetAdByID(s.dbContext, *createPicParam.AdID)
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

	_, err = s.db.GetAdByID(s.dbContext, adID)
	if err != nil {
		return c.JSON(http.StatusNotFound, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("ad not found"),
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
