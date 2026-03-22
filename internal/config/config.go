package config

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	HTTPAddr       string
	DatabaseDSN    string
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	BaseRetryDelay time.Duration
	MaxRetryCount  int
}

func Load() Config {
	c := config.New()

	_ = c.LoadEnvFiles(".env")
	c.EnableEnv("APP")

	c.SetDefault("http.addr", "8080")
	c.SetDefault("db.dsn", "postgres://postgres:postgres@localhost:5432/delayed_notifier?sslmode=disable")
	c.SetDefault("redis.addr", "localhost:6379")
	c.SetDefault("redis.password", "")
	c.SetDefault("redis.db", 0)
	c.SetDefault("retry.base_delay", time.Minute)
	c.SetDefault("retry.max", 5)

	dsn := getDatabaseDSN(c)

	return Config{
		HTTPAddr:       c.GetString("http.addr"),
		DatabaseDSN:    dsn,
		RedisAddr:      c.GetString("redis.addr"),
		RedisPassword:  c.GetString("redis.password"),
		RedisDB:        c.GetInt("redis.db"),
		BaseRetryDelay: c.GetDuration("retry.base_delay"),
		MaxRetryCount:  c.GetInt("retry.max"),
	}
}

func getDatabaseDSN(c *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.GetString("postgres.user"),
		c.GetString("postgres.pass"),
		c.GetString("postgres.host"),
		c.GetString("postgres.port"),
		c.GetString("postgres.db"),
		c.GetString("postgres.ssl.mode"))
}
