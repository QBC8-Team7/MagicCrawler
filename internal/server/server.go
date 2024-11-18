package server

import (
	"context"
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/pkg/watchlist"

	"github.com/QBC8-Team7/MagicCrawler/internal/middleware"
	"github.com/labstack/echo/v4"
	echoMiddlewares "github.com/labstack/echo/v4/middleware"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	myredis "github.com/QBC8-Team7/MagicCrawler/pkg/redis"
)

type Server struct {
	router    *echo.Echo
	logger    *logger.AppLogger
	cfg       *config.Config
	db        *sqlc.Queries
	redis     *myredis.RedisClient
	dbContext context.Context
}

func NewServer(dbCtx context.Context, cfg *config.Config, db *sqlc.Queries, redisClient *myredis.RedisClient) (*Server, error) {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path, cfg.Logger.SysPath)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, "")

	appLogger.StartSystemMetricsLogging()

	e := echo.New()

	s := &Server{
		router:    e,
		logger:    appLogger,
		cfg:       cfg,
		db:        db,
		dbContext: dbCtx,
		redis:     redisClient,
	}

	registerRoutes(e, s)
	return s, nil
}

func (s *Server) Run() error {
	defer watchlist.GetService(s.dbContext, s.redis).StopAll()

	s.router.Use(echoMiddlewares.CORSWithConfig(echoMiddlewares.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	s.router.Use(middleware.WithRequestLogger(s.logger))
	s.router.Use(middleware.WithAuthentication(s.dbContext, s.db))

	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)

	return s.router.Start(addr)

}
