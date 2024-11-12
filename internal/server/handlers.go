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
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}

func (s *Server) checkUserAccessToAd(userRole, userID string, adID int64) (bool, error) {
	userAds, err := s.db.GetUserAds(s.dbContext, &userID)
	if err != nil {
		return false, fmt.Errorf("error checking user ads: %w", err)
	}

	// Check if the ad belongs to the user
	for _, id := range userAds {
		if id != nil && *id == adID {
			return true, nil
		}
	}

	// Check if they have an admin role
	if userRole == string(sqlc.UserRoleAdmin) || userRole == string(sqlc.UserRoleSuperAdmin) {
		return true, nil
	}

	// User neither owns the ad nor has privileges
	return false, nil
}

func healthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: "ok",
	})
}

func (s *Server) createAd(c echo.Context) error {
	adParam := new(sqlc.CreateAdParams)
	if err := c.Bind(adParam); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("cannot parse the body into ad: %v", err),
		})
	}

	adParam.PublisherAdKey = "mini-app"
	adParam.PublisherID = nil
	adParam.PublishedAt = time.Now()
	adParam.Url = nil

	ad, err := s.db.CreateAd(s.dbContext, *adParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("internal error while creating new ad: %v", err),
		})
	}

	tgID := c.Request().Header.Get("Authorization")
	if tgID == "" {
		return c.JSON(http.StatusUnauthorized, jsonResponse{
			Success: false,
			Message: "cannot get user id from header",
		})
	}

	createUserAdParam := &sqlc.CreateUserAdParams{
		UserID: &tgID,
		AdID:   &ad.ID,
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
	if !ok {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	userRole, ok := c.Get("UserRole").(string)
	if !ok {
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

func (s *Server) setPriceOnAd(c echo.Context) error {
	priceParam := new(sqlc.CreatePriceParams)

	if err := c.Bind(priceParam); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse{
			Success: false,
			Message: fmt.Sprintf("cannot parse the body into price: %v", err),
		})
	}

	userID, ok := c.Get("UserID").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, jsonResponse{
			Success: false,
			Message: "user ID is not set",
		})
	}

	userRole, ok := c.Get("UserRole").(string)
	if !ok {
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
			Message: "ad not found",
		})
	}

	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: latestPrice,
	})
}
