package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func healthCheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func (s *Server) rootHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{"status": "very OK"})
}
