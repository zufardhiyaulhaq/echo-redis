package redis_client

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zufardhiyaulhaq/echo-redis/pkg/settings"
)

var ctx = context.Background()

type RedisClientInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Close() error
}

type RedisClient struct {
	Client RedisClientInterface
}

func (r RedisClient) Set(key string, value string) error {
	status := r.Client.Set(ctx, key, value, 1*time.Hour)

	return status.Err()
}

func (r RedisClient) Get(key string) (string, error) {
	out := r.Client.Get(ctx, key)

	if out.Err() != nil {
		return "", out.Err()
	}

	return out.Val(), nil
}

func (r RedisClient) Close() error {
	return r.Client.Close()
}

func NewCluster(settings settings.Settings) RedisClient {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       settings.RedisHosts,
		MaxRetries:  settings.RedisRetry,
		IdleTimeout: -1,
	})

	return RedisClient{
		Client: client,
	}
}

func New(settings settings.Settings) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:        settings.RedisHosts[0],
		MaxRetries:  settings.RedisRetry,
		IdleTimeout: -1,
	})

	return RedisClient{
		Client: client,
	}
}
