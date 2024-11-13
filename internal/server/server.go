package server

import (
	"context"
	"fmt"
	"log"

	"github.com/QBC8-Team7/MagicCrawler/internal/middleware"
	"github.com/labstack/echo/v4"
	echoMiddlewares "github.com/labstack/echo/v4/middleware"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
)

type Server struct {
	router    *echo.Echo
	logger    *logger.AppLogger
	cfg       *config.Config
	db        *sqlc.Queries
	dbContext context.Context
}

func NewServer(dbCtx context.Context, cfg *config.Config, db *sqlc.Queries) (*Server, error) {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, "")

	e := echo.New()

	s := &Server{
		router:    e,
		logger:    appLogger,
		cfg:       cfg,
		db:        db,
		dbContext: dbCtx,
	}

	registerRoutes(e, s)
	return s, nil
}

func (s *Server) Run() error {
	certFile := "/root/cert.crt"
	keyFile := "/root/private.key"

	s.router.Use(echoMiddlewares.CORSWithConfig(echoMiddlewares.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	s.router.Use(middleware.EchoRequestLogger(s.logger))
	s.router.Use(middleware.EchoAuthentication(s.dbContext, s.db))
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	log.Println(addr)

	if s.cfg.Server.Mode == "development" {
		return s.router.Start(addr)
	}
	return s.router.StartTLS(":443", certFile, keyFile)

}
