package dependencies

import (
	"context"
	"time"

	"github.com/Netflix/go-env"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

type RedisConfig struct {
	RedisUri  string `env:"REDIS_URI,required=true"`
	RedisDb   int    `env:"REDIS_DB,required=true"`
	RedisPass string `env:"REDIS_PASS"`
}

func NewRedisConfig() (*RedisConfig, error) {
	var cfg RedisConfig

	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewRedisClient(cfg *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUri,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDb,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cmd := client.Ping(ctx)

	if err := cmd.Err(); err != nil {
		return nil, err
	}

	return client, nil
}

var RedisModule = fx.Module(
	"Redis",
	fx.Provide(
		NewRedisConfig,
		NewRedisClient,
	),
	fx.Invoke(func(lc fx.Lifecycle, client *redis.Client) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				_ = client.Shutdown(ctx)

				return nil
			},
		})
	}),
)
