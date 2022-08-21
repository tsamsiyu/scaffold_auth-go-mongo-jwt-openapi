package dependencies

import (
	"context"
	"time"

	"apart-deal-api/pkg/worker/signup"

	"go.uber.org/fx"

	pkgScheduler "apart-deal-api/pkg/worker/scheduler"
)

var WorkerModule = fx.Module(
	"NotificationWorker",
	fx.Provide(
		signup.NewNotificationHandler,
		signup.NewNotificationWorker,
		signup.NewObsoleteReqWorker,
		pkgScheduler.NewScheduler,
	),
	fx.Invoke(func(
		scheduler *pkgScheduler.Scheduler,
		notificationWorker *signup.NotificationWorker,
		obsoleteReqWorker *signup.ObsoleteReqWorker,
	) {
		scheduler.Register(notificationWorker, time.Minute)
		scheduler.Register(obsoleteReqWorker, time.Minute)
	}),
	fx.Invoke(func(lc fx.Lifecycle, scheduler *pkgScheduler.Scheduler) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				scheduler.Start(context.Background())

				return nil
			},
		})
	}),
)
