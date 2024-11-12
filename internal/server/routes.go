package server

import "github.com/labstack/echo/v4"

func registerRoutes(e *echo.Echo, s *Server) {
	s.router.GET("/healthz", healthCheckHandler)

	adGroup := e.Group("/ad")
	priceGroup := e.Group("/price")
	userGroup := e.Group("/user")

	adGroup.DELETE("/:adID", s.deleteAdByID)
	adGroup.GET("/:adID", s.getAdById)
	adGroup.GET("", s.getAllAds)
	adGroup.POST("", s.createAd)

	priceGroup.GET("/:adID/all", s.getAdsAllPrices)
	priceGroup.GET("/:adID", s.getAdsLatestPrice)
	priceGroup.POST("", s.setPriceOnAd)

	userGroup.GET("/ad", s.getUsersAds)
	//userGroup.GET("/favorite", s.getUsersFavoriteAdsHandler)
}
