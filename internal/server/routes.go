package server

import (
	"github.com/labstack/echo/v4"
)

func registerRoutes(e *echo.Echo, s *Server) {
	s.router.GET("/healthz", healthCheck)

	adGroup := e.Group("/ad")
	adGroup.DELETE("/:adID", s.deleteAdByID)
	adGroup.GET("/:adID", s.getAdById)
	adGroup.GET("", s.getAllAds)
	adGroup.POST("", s.createAd)

	userGroup := e.Group("/user")
	userGroup.GET("/ad", s.getUsersAds)
	userGroup.POST("/favorite/:adID", s.createUserFavoriteAd)
	userGroup.DELETE("/favorite/:adID", s.deleteUserFavoriteAd)
	userGroup.GET("/favorite", s.getUserFavoriteAds)
	userGroup.GET("", s.getUserInfo)

	priceGroup := e.Group("/price")
	priceGroup.GET("/:adID/all", s.getAdsAllPrices)
	priceGroup.GET("/:adID", s.getAdsLatestPrice)
	priceGroup.POST("", s.setPriceOnAd)

	pictureGroup := e.Group("/picture")
	pictureGroup.GET("/:adID", s.getAdPictures)
	pictureGroup.DELETE("/:pictureID", s.deletePicture)
	pictureGroup.POST("", s.createAdPicture)

}
