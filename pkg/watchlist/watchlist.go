package watchlist

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
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
	db       *sqlc.Queries
}

var (
	managerInstance *Service
	once            sync.Once
)

func GetService(ctx context.Context, redisClient *myredis.RedisClient, db *sqlc.Queries) *Service {
	once.Do(func() {
		managerInstance = &Service{
			channels: make(map[string]chan struct{}),
			redis:    redisClient,
			ctx:      ctx,
			db:       db,
		}
	})
	return managerInstance
}

func findDifference(adIds, prevRes []int64) []int64 {
	prevResMap := make(map[int64]struct{})
	for _, id := range prevRes {
		prevResMap[id] = struct{}{}
	}

	var diff []int64
	for _, id := range adIds {
		if _, exists := prevResMap[id]; !exists {
			diff = append(diff, id)
		}
	}

	return diff
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
		for {
			select {
			case <-time.After(period):
				appLogger.Infof("checking for new add for %s", userID)

				prevResStr, err := wch.redis.Get(wch.ctx, myredis.CollectionFilterResponse, userID)
				if err != nil {
					appLogger.Errorf("failed to check for previous response foruser  %s: %v", userID, err)
				}
				var prevRes []int64
				err = json.Unmarshal([]byte(prevResStr), &prevRes)
				if err != nil {
					appLogger.Errorf("failed to unmarshal previous response for user %s: %v", userID, err)
				}

				filterStr, err := wch.redis.Get(wch.ctx, myredis.CollectionFilter, userID)
				if err != nil {
					appLogger.Errorf("failed to check for filter for user %s: %v", userID, err)
				}
				var filter sqlc.FilterAdsParams
				err = json.Unmarshal([]byte(filterStr), &filter)
				if err != nil {
					appLogger.Errorf("failed to unmarshal filter for user %s: %v", userID, err)
				}

				newAds, err := wch.db.FilterAds(wch.ctx, filter)
				if err != nil {
					appLogger.Errorf("failed to filter for user %s: %v", userID, err)
				}

				adIds := make([]int64, len(newAds))
				for i, ad := range newAds {
					adIds[i] = ad.ID
				}

				notificationService, err := notification.GetService(botToken, appLogger)
				if err != nil {
					appLogger.Warnf("Get notification notificationService error: %v", err)
				}

				diff := findDifference(adIds, prevRes)
				if len(diff) > 0 {
					err = notificationService.SendMessage(userID, "Some new ads registered on your magic crawler!")
					if err != nil {
						appLogger.Warnf("Send notification message error: %v", err)
					}
				} else {
					err = notificationService.SendMessage(userID, "I checked... there is no new ad with your filters")
					if err != nil {
						appLogger.Warnf("Send notification message error: %v", err)
					}
				}

			case <-cancelChan:
				return
			}
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
