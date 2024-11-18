package watchlist

import (
	"context"
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	"github.com/QBC8-Team7/MagicCrawler/pkg/notification"
	myredis "github.com/QBC8-Team7/MagicCrawler/pkg/redis"
	"sync"
	"time"
)

type Service struct {
	mu       sync.Mutex
	channels map[string]chan struct{}
	redis    *myredis.RedisClient
	ctx      context.Context
}

var (
	managerInstance *Service
	once            sync.Once
)

func GetService(ctx context.Context, redisClient *myredis.RedisClient) *Service {
	once.Do(func() {
		managerInstance = &Service{
			channels: make(map[string]chan struct{}),
			redis:    redisClient,
			ctx:      ctx,
		}
	})
	return managerInstance
}

func (wch *Service) StartWatch(botToken string, appLogger *logger.AppLogger, userID string, period int) {
	wch.mu.Lock()
	defer wch.mu.Unlock()

	if ch, exists := wch.channels[userID]; exists {
		close(ch)
		delete(wch.channels, userID)
	}

	cancelChan := make(chan struct{})
	wch.channels[userID] = cancelChan

	go func(userID string, period time.Duration, cancelChan chan struct{}) {
		select {
		case <-time.After(time.Duration(10) * time.Second):
			fmt.Printf("Watchlist period expired for user: %s\n", userID)

			// todo: check for new ads and notify user

			notificationService, err := notification.GetService(botToken, appLogger)
			if err != nil {
				appLogger.Warnf("Get notification notificationService error: %v", err)
			}
			err = notificationService.SendMessage(userID, "STH NEW ADDED TO TOUR ASS")
			if err != nil {
				appLogger.Warnf("Send notification message error: %v", err)
			}
		case <-cancelChan:
			return
		}
	}(userID, time.Duration(period)*time.Minute, cancelChan)
}

func (wch *Service) StopWatch(userID string) {
	wch.mu.Lock()
	defer wch.mu.Unlock()

	if ch, exists := wch.channels[userID]; exists {
		close(ch)
		delete(wch.channels, userID)
	}
}

func (wch *Service) StopAll() {
	wch.mu.Lock()
	defer wch.mu.Unlock()

	for userID, ch := range wch.channels {
		close(ch)
		delete(wch.channels, userID)
		fmt.Printf("Watchlist stopped for user: %s\n", userID)
	}
}
