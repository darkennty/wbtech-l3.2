package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"WBTech_L3.2/internal/api/handler"
	"WBTech_L3.2/internal/api/server"
	"WBTech_L3.2/internal/cache"
	"WBTech_L3.2/internal/config"
	"WBTech_L3.2/internal/repository"
	"WBTech_L3.2/internal/service"
	"github.com/wb-go/wbf/zlog"
)

func Run() {
	_ = os.Setenv("TZ", "UTC")

	zlog.InitConsole()
	logger := zlog.Logger
	cfg := config.Load()

	db, err := repository.NewPostgresDB(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init database")
	}
	defer func() {
		_ = db.Master.Close()
		for _, s := range db.Slaves {
			_ = s.Close()
		}
	}()

	redisClient := cache.NewRedisClient(cfg)
	defer func() {
		if redisClient != nil {
			_ = redisClient.Close()
		}
	}()

	linkCache := cache.NewLinkCache(redisClient)
	repo := repository.NewRepository(db)
	services := service.NewService(repo, linkCache)

	srv := new(server.Server)
	handlers := handler.NewHandler(services, logger)

	go func() {
		logger.Info().Str("addr", cfg.HTTPAddr).Msg("starting http server")
		if err = srv.Run(cfg.HTTPAddr, handlers.InitRoutes()); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logger.Fatal().Err(err).Msg("http server error")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	logger.Info().Msg("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("http server shutdown error")
	}
}
