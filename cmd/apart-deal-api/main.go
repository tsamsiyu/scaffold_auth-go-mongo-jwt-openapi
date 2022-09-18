package main

import (
	"apart-deal-api/dependencies"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	logger := dependencies.LoggerFromEnv()

	app := fx.New(
		fx.Supply(logger),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		dependencies.ConfigModule,
		dependencies.DbModule,
		dependencies.SmtpModule,
		dependencies.RepositoryModule,
		dependencies.AuthServicesModule,
		dependencies.ApiModule,
	)

	app.Run()

	if err := app.Err(); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("API stopped")
}
