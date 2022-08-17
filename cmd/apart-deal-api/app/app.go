package app

import (
	"apart-deal-api/dependencies"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger) error {
	app := fx.New(
		fx.Supply(logger),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		dependencies.ConfigModule,
		dependencies.RedisModule,
		dependencies.DbModule,
		dependencies.SmtpModule,
		dependencies.RepositoryModule,
		dependencies.AuthServicesModule,
	)

	app.Run()

	if err := app.Err(); err != nil {
		logger.Fatal(err.Error())
	}

	return nil
}
