package signup

import (
	"context"
	"fmt"
	"time"

	"apart-deal-api/pkg/mail"

	"go.uber.org/zap"

	userStore "apart-deal-api/pkg/store/user"
)

type Handler struct {
	mailer   mail.Mailer
	userRepo userStore.UserRepository
	logger   *zap.Logger
}

func NewHandler(
	mailer mail.Mailer,
	userRepo userStore.UserRepository,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		mailer:   mailer,
		logger:   logger,
		userRepo: userRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, user *userStore.User) error {
	h.logger.Info(fmt.Sprintf("Handler is starting for %s", user.UID))

	if err := h.sendNotification(ctx, user); err != nil {
		return err
	}

	if err := h.userRepo.SaveNotifiedSignUpReqTime(ctx, user.UID, time.Now()); err != nil {
		return err
	}

	return nil
}

func (h *Handler) sendNotification(ctx context.Context, user *userStore.User) error {
	body := fmt.Sprintf(
		`Hello dear %s!
Here's your confirmation code: %s`,
		user.Name,
		user.SignUpReq.Code,
	)

	if err := h.mailer.Send(ctx, mail.Letter{
		To:   []string{user.Email},
		Body: body,
	}); err != nil {
		return err
	}

	return nil
}
