package server

import (
	"apart-deal-api/pkg/api/aspects"
	"apart-deal-api/pkg/api/handlers/auth"
	"apart-deal-api/pkg/config"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func NewServer(logger *zap.Logger, cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = aspects.NewErrorHandler(logger)
	e.Use(aspects.NewLoggingMiddleware(logger, &aspects.LoggingMiddlewareConfig{
		IncludeRequestBodies:  cfg.IsDebug,
		IncludeResponseBodies: cfg.IsDebug,
	}))
	e.Use(aspects.NewTracingMiddleware())

	return e
}

func NewAuthRouteGroup(e *echo.Echo) auth.RouteGroup {
	return e.Group("/api/v1/auth")
}

func RegisterRoutes(
	e *echo.Echo,
	authGroup auth.RouteGroup,
	signUpHandler *auth.SignUpHandler,
	signUpConfirmHandler *auth.SignUpConfirmHandler,
	signInHandler *auth.SignInHandler,
) {
	e.GET("ready", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	auth.RegisterSignUpRoute(authGroup, signUpHandler)
	auth.RegisterSignUpConfirmRoute(authGroup, signUpConfirmHandler)
	auth.RegisterSignInRoute(authGroup, signInHandler)
}
