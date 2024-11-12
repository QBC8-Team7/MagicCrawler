package server

import "github.com/labstack/echo/v4"

func registerRoutes(e *echo.Echo, s *Server) {
	s.router.GET("/healthz", healthCheckHandler)

	adGroup := e.Group("/ad")
	priceGroup := e.Group("/price")

	adGroup.DELETE("/:adID", s.deleteAdByID)
	adGroup.POST("", s.createAd)
	adGroup.GET("/:adID", s.getAdById)
	adGroup.GET("", s.getAllAds)

	priceGroup.POST("", s.setPriceOnAd)
	priceGroup.GET("/:adID", s.getAdsLatestPrice)
	//priceGroup.GET("/:adID/all", s.getAdsAllPricesHandler)

	userGroup := e.Group("/user")
	userGroup.GET("/ad", s.getUsersAds)
	//userGroup.GET("/favorite", s.getUsersFavoriteAdsHandler)
}
