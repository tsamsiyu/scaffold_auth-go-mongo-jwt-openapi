package dependencies

import (
	authHandlers "apart-deal-api/pkg/api/handlers/auth"
	"apart-deal-api/pkg/api/server"
	"apart-deal-api/pkg/auth"
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Netflix/go-env"
	"go.uber.org/fx"
)

type ApiRunFn func(ctx context.Context) error

type ApiConfig struct {
	Port int `env:"API_PORT,required=true"`
}

func NewApiConfig() (*ApiConfig, error) {
	var cfg ApiConfig

	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewApiRunFn(e *echo.Echo, logger *zap.Logger, shutdowner fx.Shutdowner, apiCfg *ApiConfig) ApiRunFn {
	return func(ctx context.Context) error {
		logger.Info(fmt.Sprintf("Starting API on port %d", apiCfg.Port))
		errCh := make(chan error)

		go func() {
			defer close(errCh)

			uri := fmt.Sprintf(":%d", apiCfg.Port)
			if err := e.Start(uri); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return errors.Errorf("Could not manage to start API server during fx startup time")
		case err := <-errCh:
			if err != nil {
				return errors.Errorf("API server failed during startup with error: %s", err)
			}

			return errors.Errorf("API server stopped during startup")
		case <-time.After(time.Second * 2):
			logger.Info("API considered running")

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
}

var ApiModule = fx.Module(
	"API",
	fx.Provide(
		NewApiConfig,
		NewApiRunFn,
		server.NewServer,
		authHandlers.NewAuthHandler,
		server.RegisterRoutes,
		auth.NewTokenStore,
	),
	fx.Invoke(func(lc fx.Lifecycle, fn ApiRunFn, e *echo.Echo) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return fn(ctx)
			},
			OnStop: func(ctx context.Context) error {
				_ = e.Shutdown(ctx)

				return nil
			},
		})
	}),
)
