package redisdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

const (
	CollectionAdCount        string = "ads"
	CollectionFilter         string = "filter"
	CollectionFilterResponse string = "filter_response"
	CollectionRateLimit      string = "rate_limit"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(ctx context.Context, host, port, password string, db int) *RedisClient {
	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisClient{Client: client, Ctx: ctx}
}

func (r *RedisClient) Set(ctx context.Context, collection string, key string, value interface{}) error {
	prefixedKey := fmt.Sprintf("%s:%s", collection, key)
	return r.Client.Set(ctx, prefixedKey, value, 0).Err()
}

func (r *RedisClient) Get(ctx context.Context, collection string, key string) (string, error) {
	prefixedKey := fmt.Sprintf("%s:%s", collection, key)
	return r.Client.Get(ctx, prefixedKey).Result()
}

func (r *RedisClient) Delete(ctx context.Context, collection string, key string) error {
	prefixedKey := fmt.Sprintf("%s:%s", collection, key)
	return r.Client.Del(ctx, prefixedKey).Err()
}

func (r *RedisClient) GetAllFromCollection(ctx context.Context, collection string, count int64) (map[string]string, error) {
	pattern := fmt.Sprintf("%s:*", collection)
	result := make(map[string]string)

	iter := r.Client.Scan(ctx, 0, pattern, count).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		value, err := r.Client.Get(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get value for key %s: %w", key, err)
		}

		trimmedKey := key[len(collection)+1:]
		result[trimmedKey] = value
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate keys: %w", err)
	}

	return result, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
