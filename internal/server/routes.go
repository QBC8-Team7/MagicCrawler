package server

import (
	"github.com/labstack/echo/v4"
)

func registerRoutes(e *echo.Echo, s *Server) {
	apiV1 := e.Group("/api/v1")

	apiV1.GET("/healthz", healthCheck)

	adGroup := apiV1.Group("/ad")
	adGroup.DELETE("/:adID", s.deleteAdByID)
	adGroup.GET("/search", s.searchAds)
	adGroup.GET("/populars", s.getPopularAds)
	adGroup.GET("/:adID", s.getAdById)
	adGroup.GET("", s.getAllAds)
	adGroup.POST("", s.createAd)

	userGroup := apiV1.Group("/user")
	userGroup.GET("/ad", s.getUsersAds)
	userGroup.POST("/favorite/:adID", s.createUserFavoriteAd)
	userGroup.DELETE("/favorite/:adID", s.deleteUserFavoriteAd)
	userGroup.GET("/favorite", s.getUserFavoriteAds)
	userGroup.GET("", s.getUserInfo)

	priceGroup := apiV1.Group("/price")
	priceGroup.GET("/:adID/all", s.getAdsAllPrices)
	priceGroup.GET("/:adID", s.getAdsLatestPrice)
	priceGroup.POST("", s.setPriceOnAd)

	pictureGroup := apiV1.Group("/picture")
	pictureGroup.GET("/:adID", s.getAdPictures)
	pictureGroup.DELETE("/:pictureID", s.deletePicture)
	pictureGroup.POST("", s.createAdPicture)
}
