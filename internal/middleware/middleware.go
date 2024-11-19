package middleware

import (
	"context"
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	myredis "github.com/QBC8-Team7/MagicCrawler/pkg/redis"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
)

type jsonResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
}

func WithRequestLogger(logger *logger.AppLogger) echo.MiddlewareFunc {
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

func WithAuthentication(ctx context.Context, db *sqlc.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ignoredPaths := map[string]bool{
				"/healthz": true,
			}

			if _, ok := ignoredPaths[c.Path()]; ok {
				return next(c)
			}

			userTgID := c.Request().Header.Get("Authorization")
			if userTgID == "" {
				return c.JSON(http.StatusUnauthorized, jsonResponse{
					Success: false,
					Message: "authorization header required",
				})
			}

			var user sqlc.User
			user, err := db.GetUserByTGID(ctx, userTgID)
			if err != nil {
				var role sqlc.NullUserRole
				_ = role.Scan("simple")

				period := int32(0)

				param := sqlc.CreateUserParams{TgID: userTgID, Role: role, WatchlistPeriod: &period}
				user, err = db.CreateUser(ctx, param)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, jsonResponse{
						Success: false,
						Message: err,
					})
				}
			}

			c.Set("UserRole", string(user.Role.UserRole))
			c.Set("UserID", user.TgID)

			return next(c)
		}
	}
}

type RateLimiterConfig struct {
	Limit  int
	Window time.Duration
}

func WithRateLimiter(ctx context.Context, redisClient *myredis.RedisClient, config RateLimiterConfig, logger *logger.AppLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get("UserID").(string)
			if !ok || userID == "" {
				logger.Warnf("UserID not found in context")
				return c.JSON(http.StatusUnauthorized, jsonResponse{
					Success: false,
					Message: "unauthorized: UserID required",
				})
			}

			key := fmt.Sprintf("%s:%s", myredis.CollectionRateLimit, userID)

			currentCount, err := redisClient.Client.Incr(ctx, key).Result()
			if err != nil {
				logger.Errorf("Failed to increment Redis key: %v", err)
				return c.JSON(http.StatusInternalServerError, jsonResponse{
					Success: false,
					Message: "internal server error",
				})
			}

			if currentCount == 1 {
				err = redisClient.Client.Expire(ctx, key, config.Window).Err()
				if err != nil {
					logger.Errorf("Failed to set TTL for Redis key: %v", err)
					return c.JSON(http.StatusInternalServerError, jsonResponse{
						Success: false,
						Message: "internal server error",
					})
				}
			}

			if currentCount > int64(config.Limit) {
				ttl, _ := redisClient.Client.TTL(ctx, key).Result()
				logger.Warnf("Rate limit exceeded for user %s", userID)
				return c.JSON(http.StatusTooManyRequests, jsonResponse{
					Success: false,
					Message: fmt.Sprintf("rate limit exceeded, retry after %s", ttl.String()),
				})
			}

			return next(c)
		}
	}
}
