package repository

import (
	"context"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

func NewPostgresDB(ctx context.Context, dsn string) (*dbpg.DB, error) {
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(dsn, nil, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = db.Master.PingContext(ctx); err != nil {
		_ = db.Master.Close()
		for _, s := range db.Slaves {
			_ = s.Close()
		}
		return nil, err
	}

	return db, nil
}
