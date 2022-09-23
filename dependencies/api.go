package dependencies

import (
	"context"
	"fmt"
	"strings"
	"time"

	"apart-deal-api/pkg/api/auth"
	"apart-deal-api/pkg/api/server"
	"apart-deal-api/pkg/store/user"

	"github.com/Netflix/go-env"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	authHandlers "apart-deal-api/pkg/api/handlers/auth"
)

type ApiRunFn func(ctx context.Context) error

type ApiConfig struct {
	Port         int    `env:"API_PORT,required=true"`
	AllowOrigins string `env:"ALLOW_ORIGINS"`
	TokenSecret  string `env:"JWT_SECRET,required=true"`
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
				if !strings.Contains(err.Error(), "Server closed") {
					errCh <- err
				}
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

func NewAuthenticationService(cfg *ApiConfig, userRepo user.UserRepository) *auth.AuthenticationService {
	return auth.NewAuthenticationService(cfg.TokenSecret, userRepo)
}

var ApiModule = fx.Module(
	"API",
	fx.Provide(
		NewApiConfig,
		NewApiRunFn,
		server.NewServer,
		server.NewAuthRouteGroup,
		NewAuthenticationService,
		authHandlers.NewSignUpHandler,
		authHandlers.NewSignUpConfirmHandler,
		authHandlers.NewSignInHandler,
	),
	fx.Invoke(func(cfg *ApiConfig, e *echo.Echo) {
		if cfg.AllowOrigins == "" {
			return
		}

		origins := strings.Split(cfg.AllowOrigins, ",")
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: origins,
		}))
	}),
	fx.Invoke(server.RegisterRoutes),
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
