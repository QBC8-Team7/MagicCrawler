package server

import "github.com/labstack/echo/v4"

func registerRoutes(e *echo.Echo, s *Server) {
	e.GET("/healthz", healthCheckHandler)

	e.POST("/ad", s.createAdHandler)
}
