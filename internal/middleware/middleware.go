package middleware

import (
	"context"
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	"gopkg.in/telebot.v4"
)

func EchoRequestLogger(logger *logger.AppLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			logger.Infof("Received %s request for %s from %s at %s",
				req.Method,
				req.URL.Path,
				req.RemoteAddr,
				time.Now().Format(time.RFC3339))

			err := next(c)

			logger.Infof("Responded with status %d for %s request to %s",
				res.Status,
				req.Method,
				req.URL.Path)

			return err
		}
	}
}

func EchoAuthentication(ctx context.Context, db *sqlc.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userTgID := c.Request().Header.Get("Authorization")
			if userTgID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
			}

			var user sqlc.User
			user, err := db.GetUserByTGID(ctx, userTgID)
			if err != nil {
				var role sqlc.NullUserRole
				_ = role.Scan("simple")

				var period int32
				period = 0

				param := sqlc.CreateUserParams{TgID: userTgID, Role: role, WatchlistPeriod: &period}
				user, err = db.CreateUser(ctx, param)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("can not create user: %w", err))
				}
			}

			c.Set("UserRole", user.Role.UserRole)

			return next(c)
		}
	}
}

func WithLogging(logger *logger.AppLogger) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			msg := c.Message()
			if msg != nil {
				logger.Infof("Received message from %s (%d): %s at %s",
					msg.Sender.Username,
					msg.Sender.ID,
					msg.Text,
					time.Now().Format(time.RFC3339))
			}
			return next(c)
		}
	}
}
