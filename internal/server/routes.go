package server

import "github.com/labstack/echo/v4"

func registerRoutes(e *echo.Echo, s *Server) {
	e.GET("/", s.rootHandler)
	e.GET("healthz", healthCheckHandler)
}
