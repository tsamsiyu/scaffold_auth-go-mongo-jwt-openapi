package signup

import (
	"context"
	"time"

	"go.uber.org/zap"

	userStore "apart-deal-api/pkg/store/user"
)

type NotificationWorker struct {
	logger   *zap.Logger
	handler  *NotificationHandler
	userRepo userStore.UserRepository
}

func NewNotificationWorker(
	userRepo userStore.UserRepository,
	handler *NotificationHandler,
	logger *zap.Logger,
) *NotificationWorker {
	return &NotificationWorker{
		userRepo: userRepo,
		handler:  handler,
		logger:   logger,
	}
}

func (w *NotificationWorker) Process(ctx context.Context) error {
	users, err := w.userRepo.FindAllNotNotifiedSignUpRequests(ctx)
	if err != nil {
		return err
	}

	w.logger.With(zap.Int("count", len(users))).Info("Found sign up req for sending notifications")

	for _, user := range users {
		if err := w.processItem(ctx, &user); err != nil {
			return err
		}
	}

	return nil
}

func (w *NotificationWorker) processItem(ctx context.Context, user *userStore.User) error {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := w.handler.Handle(childCtx, user); err != nil {
		return err
	}

	return nil
}
