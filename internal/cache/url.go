package cache

import (
	"context"
	"errors"
	"time"

	"github.com/wb-go/wbf/redis"
)

type LinkCache interface {
	Set(ctx context.Context, shortUrl, longUrl string) error
	Get(ctx context.Context, shortUrl string) (string, error)
	Delete(ctx context.Context, shortUrl string) error
}

type redisLinkCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewLinkCache(client *redis.Client) LinkCache {
	if client == nil {
		return &noopLinkCache{}
	}
	return &redisLinkCache{client: client, ttl: 10 * time.Minute}
}

func (c *redisLinkCache) Set(ctx context.Context, shortUrl, longUrl string) error {
	return c.client.SetWithExpiration(ctx, cacheKey(shortUrl), longUrl, c.ttl)
}

func (c *redisLinkCache) Get(ctx context.Context, shortUrl string) (string, error) {
	res, err := c.client.Get(ctx, cacheKey(shortUrl))
	if err != nil {
		if errors.Is(err, redis.NoMatches) {
			return "", nil
		}
		return "", err
	}

	return res, nil
}

func (c *redisLinkCache) Delete(ctx context.Context, shortUrl string) error {
	return c.client.Del(ctx, cacheKey(shortUrl))
}

type noopLinkCache struct{}

func (n *noopLinkCache) Set(ctx context.Context, shortUrl, longUrl string) error { return nil }
func (n *noopLinkCache) Get(ctx context.Context, shortUrl string) (string, error) {
	return "", nil
}
func (n *noopLinkCache) Delete(ctx context.Context, shortUrl string) error { return nil }

func cacheKey(id string) string {
	return "shortener:" + id
}
