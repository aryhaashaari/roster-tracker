package cachex

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Define callback, when returning error
type Closure func(bytes []byte) error

const (
	cacheNil string = `redis: nil`
)

// AgentCache contract
type Cacher interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val any, duration time.Duration) error
	Delete(ctx context.Context, key ...string) error
}

type cache struct {
	rds             redis.Cmdable
	retentionSecond time.Duration
}

// NewAgentCache creates new agent redis client
func NewCache(redis redis.Cmdable) Cacher {
	return &cache{
		rds: redis,
	}
}

func (c *cache) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	cmd := c.rds.Set(ctx, key, val, exp)
	return cmd.Err()
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	cmd := c.rds.Get(ctx, key)
	b, e := cmd.Bytes()

	if e == redis.Nil {
		return b, nil
	}

	return b, e
}

func (c *cache) Delete(ctx context.Context, key ...string) error {
	cmd := c.rds.Del(ctx, key...)
	return cmd.Err()
}

func (c *cache) Sort(ctx context.Context, key string, offset, count int64, order string) ([]string, error) {
	return c.rds.Sort(ctx, key, &redis.Sort{Offset: offset, Count: count, Order: order}).Result()
}

func (c *cache) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.rds.ZRange(ctx, key, start, stop).Result()
}
