package signup_request

import (
	"context"
	"time"

	"apart-deal-api/pkg/store/user"

	"go.uber.org/zap"
)

type SignUpRequestHandler struct {
	userRepo user.UserRepository
	logger   *zap.Logger
}

func (h *SignUpRequestHandler) Handle(ctx context.Context) error {
	users, err := h.userRepo.FindAllNotNotifiedSignUpRequests(ctx)
	if err != nil {
		return err
	}

	for _, userModel := range users {
		childCtx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		if err := h.handleModel(childCtx, &userModel); err != nil {
			return err
		}
	}

	return nil
}

func (h *SignUpRequestHandler) handleModel(ctx context.Context, userModel *user.User) error {
	if err := h.sendNotification(ctx, userModel); err != nil {
		return err
	}

	if err := h.userRepo.SaveNotifiedSignUpReqTime(ctx, userModel.UID); err != nil {
		return err
	}

	return nil
}

func (h *SignUpRequestHandler) sendNotification(ctx context.Context, userModel *user.User) error {
	return nil
}
