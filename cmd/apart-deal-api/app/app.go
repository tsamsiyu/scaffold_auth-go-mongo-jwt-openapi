package app

import (
	"context"
	"fmt"
	"go.uber.org/fx/fxevent"
	"time"

	"github.com/Netflix/go-env"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	appMongo "apart-deal-api/pkg/mongo"
	appRedis "apart-deal-api/pkg/redis"
)

type Config struct {
	Port      int    `env:"API_PORT,required=true"`
	DbUri     string `env:"MONGO_URI,required=true"`
	DbName    string `env:"MONGO_DOMAIN_DB,required=true"`
	RedisUri  string `env:"REDIS_URI,required=true"`
	RedisDb   int    `env:"REDIS_DB,required=true"`
	RedisPass string `env:"REDIS_PASS"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Run(runCfg *Config, logger *zap.Logger) error {
	app := fx.New(
		fx.Supply(logger),
		fx.Supply(runCfg),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Provide(func() *echo.Echo {
			e := echo.New()
			e.HideBanner = true
			e.HidePort = true

			return e
		}),
		fx.Provide(func() (*mongo.Client, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			return appMongo.NewClient(ctx, runCfg.DbUri)
		}),
		fx.Provide(func() *redis.Client {
			return appRedis.NewClient(runCfg.RedisUri, runCfg.RedisDb, runCfg.RedisPass)
		}),
		fx.Invoke(registerFxHooks),
	)

	app.Run()

	if err := app.Err(); err != nil {
		logger.Fatal(err.Error())
	}

	return nil
}

func registerFxHooks(
	lc fx.Lifecycle,
	e *echo.Echo,
	mem *redis.Client,
	db *mongo.Client,
	shutdowner fx.Shutdowner,
	logger *zap.Logger,
	runCfg *Config,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := runAPI(e, runCfg, logger, shutdowner); err != nil {
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping API")

			_ = e.Shutdown(ctx)
			_ = db.Disconnect(ctx)
			_ = mem.Shutdown(ctx)

			return nil
		},
	})
}

func runAPI(e *echo.Echo, runCfg *Config, logger *zap.Logger, shutdowner fx.Shutdowner) error {
	logger.Info(fmt.Sprintf("Starting API on port %d", runCfg.Port))
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		uri := fmt.Sprintf(":%d", runCfg.Port)
		if err := e.Start(uri); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return errors.Errorf("API server failed during startup with error: %s", err)
		}

		return errors.Errorf("API server stopped during startup")
	case <-time.After(time.Second * 5):
		break
	}

	go func() {
		select {
		case err := <-errCh:
			if err != nil {
				logger.Error(fmt.Sprintf("API stopped with error: %s", err))
				_ = shutdowner.Shutdown()
			}
		}
	}()

	return nil
}
