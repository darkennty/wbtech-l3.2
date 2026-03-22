package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"WBTech_L3.2/internal/model"
	"github.com/wb-go/wbf/dbpg"
)

type UrlPostgresRepository struct {
	db *dbpg.DB
}

var ErrNotFound = errors.New("no rows in result set")

func NewUrlRepository(db *dbpg.DB) *UrlPostgresRepository {
	return &UrlPostgresRepository{db: db}
}

func (r *UrlPostgresRepository) CreateShortUrl(ctx context.Context, shortUrl, longUrl string) error {
	const q = `
INSERT INTO urls (
	short_url, long_url
) VALUES ($1,$2);
`
	_, err := r.db.ExecContext(ctx, q,
		shortUrl, longUrl,
	)
	return err
}

func (r *UrlPostgresRepository) GetLongUrl(ctx context.Context, shortUrl string) (string, error) {
	const q = `
SELECT long_url FROM urls WHERE short_url = $1;
`
	row := r.db.QueryRowContext(ctx, q, shortUrl)

	var longUrl string
	if err := row.Scan(&longUrl); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		} else {
			return "", err
		}
	}

	return longUrl, nil
}

func (r *UrlPostgresRepository) SaveStats(ctx context.Context, shortUrl, userAgent string) error {
	const q = `
INSERT INTO stats (
	url, time, user_agent
) VALUES ($1,$2,$3);
`

	currentTime := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, q,
		shortUrl, currentTime, userAgent,
	)

	return err
}

func (r *UrlPostgresRepository) GetStats(ctx context.Context, shortUrl string) (model.Stat, error) {
	q := `
SELECT id, url, "time", user_agent FROM stats WHERE url = $1;
`
	rows, err := r.db.QueryContext(ctx, q, shortUrl)
	if err != nil {
		return model.Stat{}, err
	}
	defer rows.Close()

	var clicks []model.Click
	for rows.Next() {
		var click model.Click
		if err = rows.Scan(&click.ID, &click.Url, &click.Time, &click.UserAgent); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return model.Stat{}, ErrNotFound
			} else {
				return model.Stat{}, err
			}
		}
		clicks = append(clicks, click)
	}

	var stat model.Stat
	stat.ClicksTotal = len(clicks)
	stat.Clicks = clicks

	return stat, nil
}

func (r *UrlPostgresRepository) GetAggregatedStats(ctx context.Context, shortUrl, aggregateBy string) ([]model.Stat, error) {
	q := `
SELECT id, url, "time", user_agent FROM stats WHERE url = $1;
`
	rows, err := r.db.QueryContext(ctx, q, shortUrl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aggregatedClicks := make(map[string][]model.Click)
	for rows.Next() {
		var click model.Click
		if err = rows.Scan(&click.ID, &click.Url, &click.Time, &click.UserAgent); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			} else {
				return nil, err
			}
		}

		if aggregateBy == "day" {
			datetime, err := time.Parse("2006-01-02T15:04:05Z07", click.Time)
			if err != nil {
				return nil, err
			}
			day := datetime.Format("2006-01-02")
			aggregatedClicks[day] = append(aggregatedClicks[day], click)
		} else if aggregateBy == "month" {
			datetime, err := time.Parse("2006-01-02T15:04:05Z07", click.Time)
			if err != nil {
				return nil, err
			}
			month := datetime.Format("2006-01")
			aggregatedClicks[month] = append(aggregatedClicks[month], click)
		} else if aggregateBy == "user_agent" || aggregateBy == "useragent" {
			aggregatedClicks[click.UserAgent] = append(aggregatedClicks[click.UserAgent], click)
		}
	}

	var stats []model.Stat
	for _, clicks := range aggregatedClicks {
		stats = append(stats, model.Stat{
			ClicksTotal: len(clicks),
			Clicks:      clicks,
		})
	}

	return stats, nil
}
