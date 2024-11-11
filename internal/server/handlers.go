package server

import (
	"github.com/labstack/echo/v4"
)

// MapRoutes for mapping all routes
func (s *Server) MapHandlers(e *echo.Echo) error {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "very OK"})
	})
	return nil
}
