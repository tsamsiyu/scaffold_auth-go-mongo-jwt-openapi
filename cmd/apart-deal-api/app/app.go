package app

import (
	"context"
	"fmt"
	"time"

	"apart-deal-api/pkg/api/server"
	"apart-deal-api/pkg/auth"
	"apart-deal-api/pkg/config"
	"apart-deal-api/pkg/store/user"

	"github.com/Netflix/go-env"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	authHandlers "apart-deal-api/pkg/api/handlers/auth"
	authDomain "apart-deal-api/pkg/domain/auth"
	appMongo "apart-deal-api/pkg/mongo"
	appRedis "apart-deal-api/pkg/redis"
)

type Envs struct {
	Port        int    `env:"API_PORT,required=true"`
	DbUri       string `env:"MONGO_URI,required=true"`
	DbName      string `env:"MONGO_DOMAIN_DB,required=true"`
	DbCaCrt     string `env:"MONGO_CA_CRT"`
	DbClientCrt string `env:"MONGO_CLIENT_CRT"`
	DbClientKey string `env:"MONGO_CLIENT_KEY"`
	RedisUri    string `env:"REDIS_URI,required=true"`
	RedisDb     int    `env:"REDIS_DB,required=true"`
	RedisPass   string `env:"REDIS_PASS"`
}

func NewEnvs() (*Envs, error) {
	var cfg Envs

	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Run(envs *Envs, logger *zap.Logger) error {
	appCfg := &config.Config{
		IsDebug: logger.Core().Enabled(zapcore.DebugLevel),
	}

	app := fx.New(
		fx.Supply(appCfg),
		fx.Supply(logger),
		fx.Supply(envs),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Provide(server.NewServer),
		fx.Provide(func() (*mongo.Client, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			var tls *appMongo.TLSConfig
			if envs.DbCaCrt != "" {
				tls = &appMongo.TLSConfig{
					CaCrt:     envs.DbCaCrt,
					ClientKey: envs.DbClientKey,
					ClientCrt: envs.DbClientCrt,
				}
			}

			return appMongo.NewClient(ctx, &appMongo.DbRef{
				URI: envs.DbUri,
				TLS: tls,
			})
		}),
		fx.Provide(func() *redis.Client {
			return appRedis.NewClient(envs.RedisUri, envs.RedisDb, envs.RedisPass)
		}),
		fx.Provide(auth.NewTokenStore),
		fx.Provide(func(client *mongo.Client, envs *Envs) appMongo.MainDB {
			return appMongo.ProvideDatabase(client, envs.DbName)
		}),
		fx.Provide(user.NewUserRepository),
		fx.Provide(authHandlers.NewAuthHandler),
		fx.Provide(authDomain.NewSignUpService),
		fx.Provide(authDomain.NewConfirmSignUpService),
		fx.Invoke(server.RegisterRoutes),
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
	envs *Envs,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := runAPI(e, envs, logger, shutdowner); err != nil {
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

func runAPI(e *echo.Echo, envs *Envs, logger *zap.Logger, shutdowner fx.Shutdowner) error {
	logger.Info(fmt.Sprintf("Starting API on port %d", envs.Port))
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		uri := fmt.Sprintf(":%d", envs.Port)
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
		logger.Info("API considered started")

		break
	}

	go func() {
		err := <-errCh
		if err != nil {
			logger.Error(fmt.Sprintf("API stopped with error: %s", err))
			_ = shutdowner.Shutdown()
		}
	}()

	return nil
}
