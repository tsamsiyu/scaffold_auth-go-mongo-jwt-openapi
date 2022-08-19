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
	//e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
	//	Timeout: time.Second * 5,
	//}))

	return e
}

func RegisterRoutes(
	e *echo.Echo,
	signUpHandler *auth.SignUpHandler,
	signUpConfirmHandler *auth.SignUpConfirmHandler,
) {

	e.GET("ready", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	group := e.Group("/api/v1")

	group.POST("/auth/sign-up", signUpHandler.Handle)
	group.POST("/auth/sign-up-confirm", signUpConfirmHandler.Handle)

	//r.Add("POST", "auth/sign-up", authHandler.SignUp)
	//r.Add("POST", "auth/confirm-sign-up", authHandler.ConfirmSignUp)
	//r.Add("POST", "auth/sign-in", authHandler.SignIn)
	//r.Add("POST", "auth/sign-out", authHandler.SignOut)
}
