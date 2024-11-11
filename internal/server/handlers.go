package server

import (
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type jsonResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}

func healthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: "ok",
	})
}

func (s *Server) createAdHandler(c echo.Context) error {
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
