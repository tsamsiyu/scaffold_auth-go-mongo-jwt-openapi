package signup

import (
	"context"
	"fmt"
	"time"

	"apart-deal-api/pkg/mail"

	"go.uber.org/zap"

	userStore "apart-deal-api/pkg/store/user"
)

type NotificationHandler struct {
	mailer   mail.Mailer
	userRepo userStore.UserRepository
	logger   *zap.Logger
}

func NewNotificationHandler(
	mailer mail.Mailer,
	userRepo userStore.UserRepository,
	logger *zap.Logger,
) *NotificationHandler {
	return &NotificationHandler{
		mailer:   mailer,
		logger:   logger,
		userRepo: userRepo,
	}
}

func (h *NotificationHandler) Handle(ctx context.Context, user *userStore.User) error {
	h.logger.
		With(zap.String("email", user.Email)).
		Info(fmt.Sprintf("NotificationHandler is starting"))

	if err := h.sendNotification(ctx, user); err != nil {
		return err
	}

	if err := h.userRepo.SaveNotifiedSignUpReqTime(ctx, user.UID, time.Now()); err != nil {
		return err
	}

	return nil
}

func (h *NotificationHandler) sendNotification(ctx context.Context, user *userStore.User) error {
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
