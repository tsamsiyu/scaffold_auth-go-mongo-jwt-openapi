package dependencies

import (
	"apart-deal-api/pkg/worker/signup"
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"go.uber.org/fx"
)

var WorkerModule = fx.Module(
	"Worker",
	fx.Provide(
		signup.NewHandler,
		signup.NewWorker,
		signup.NewScheduler,
	),
	fx.Invoke(func(
		lc fx.Lifecycle,
		signupScheduler *signup.Scheduler,
		shutdowner fx.Shutdowner,
		logger *zap.Logger,
	) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := signupScheduler.Start(context.Background()); err != nil {
						logger.Error(errors.WithStack(err).Error())

						_ = shutdowner.Shutdown()
					}
				}()

				return nil
			},
		})
	}),
)
