package service

import (
	"context"
	"errors"
	"math/rand"

	"WBTech_L3.2/internal/cache"
	"WBTech_L3.2/internal/model"
	"WBTech_L3.2/internal/repository"
)

type UrlService struct {
	repo  repository.Url
	cache cache.LinkCache
}

func NewUrlService(repo repository.Url, cache cache.LinkCache) *UrlService {
	return &UrlService{
		repo:  repo,
		cache: cache,
	}
}

func (s *UrlService) CreateShortUrl(ctx context.Context, longUrl, desiredShortUrl string) (string, error) {
	var shortUrl string
	if desiredShortUrl != "" {
		shortUrl = desiredShortUrl
		if _, err := s.repo.GetLongUrl(ctx, shortUrl); err == nil {
			return "", errors.New("this short_url already exists. try another")
		} else if !errors.Is(err, repository.ErrNotFound) {
			return "", err
		}
	} else {
		for {
			shortUrl = randomShortUrl(5)

			if _, err := s.repo.GetLongUrl(ctx, shortUrl); err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					break
				} else {
					return "", err
				}
			}
		}
	}

	if err := s.repo.CreateShortUrl(ctx, shortUrl, longUrl); err != nil {
		return "", err
	}
	_ = s.cache.Set(ctx, shortUrl, longUrl)

	return shortUrl, nil
}

func (s *UrlService) GetLongUrl(ctx context.Context, shortUrl string) (string, error) {
	if n, err := s.cache.Get(ctx, shortUrl); err == nil && shortUrl != "" {
		return n, nil
	}

	res, err := s.repo.GetLongUrl(ctx, shortUrl)
	if err != nil {
		return "", err
	}

	_ = s.cache.Set(ctx, shortUrl, res)
	return res, nil
}

func (s *UrlService) SaveStats(ctx context.Context, shortUrl, userAgent string) error {
	return s.repo.SaveStats(ctx, shortUrl, userAgent)
}

func (s *UrlService) GetStats(ctx context.Context, shortUrl string) (model.Stat, error) {
	return s.repo.GetStats(ctx, shortUrl)
}

func (s *UrlService) GetAggregatedStats(ctx context.Context, shortUrl, aggregateBy string) ([]model.Stat, error) {
	return s.repo.GetAggregatedStats(ctx, shortUrl, aggregateBy)
}

func randomShortUrl(n int) string {
	var letters = []rune("abcdefghjklmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
