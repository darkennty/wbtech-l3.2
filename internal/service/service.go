package service

import (
	"context"

	"WBTech_L3.2/internal/cache"
	"WBTech_L3.2/internal/model"
	"WBTech_L3.2/internal/repository"
)

type Url interface {
	CreateShortUrl(ctx context.Context, longUrl, desiredShortUrl string) (string, error)
	GetLongUrl(ctx context.Context, shortUrl string) (string, error)
	SaveStats(ctx context.Context, shortUrl, userAgent string) error
	GetStats(ctx context.Context, shortUrl string) (model.Stat, error)
	GetAggregatedStats(ctx context.Context, shortUrl, aggregateBy string) ([]model.Stat, error)
}

type Service struct {
	Url
}

func NewService(repo *repository.Repository, linkCache cache.LinkCache) *Service {
	return &Service{
		Url: NewUrlService(repo.Url, linkCache),
	}
}
