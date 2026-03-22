package repository

import (
	"context"

	"WBTech_L3.2/internal/model"
	"github.com/wb-go/wbf/dbpg"
)

type Url interface {
	CreateShortUrl(ctx context.Context, shortUrl, longUrl string) error
	GetLongUrl(ctx context.Context, shortUrl string) (string, error)
	SaveStats(ctx context.Context, shortUrl, userAgent string) error
	GetStats(ctx context.Context, shortUrl string) (model.Stat, error)
	GetAggregatedStats(ctx context.Context, shortUrl, aggregateBy string) ([]model.Stat, error)
}

type Repository struct {
	Url
}

func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{
		Url: NewUrlRepository(db),
	}
}
