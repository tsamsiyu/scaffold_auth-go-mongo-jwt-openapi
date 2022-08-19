package signup

import (
	"context"
	"time"

	userStore "apart-deal-api/pkg/store/user"
)

type Worker struct {
	handler  *Handler
	userRepo userStore.UserRepository
}

func NewWorker(
	userRepo userStore.UserRepository,
	handler *Handler,
) *Worker {
	return &Worker{
		userRepo: userRepo,
		handler:  handler,
	}
}

func (w *Worker) Process(ctx context.Context) error {
	users, err := w.userRepo.FindAllNotNotifiedSignUpRequests(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if err := w.processItem(ctx, &user); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) processItem(ctx context.Context, user *userStore.User) error {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := w.handler.Handle(childCtx, user); err != nil {
		return err
	}

	return nil
}
