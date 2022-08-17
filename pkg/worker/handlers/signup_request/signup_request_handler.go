package signup_request

import (
	"context"
	"fmt"
	"time"

	"apart-deal-api/pkg/mail"
	"apart-deal-api/pkg/store/user"

	"go.uber.org/zap"
)

type SignUpRequestHandler struct {
	userRepo   user.UserRepository
	logger     *zap.Logger
	mailClient *mail.Client
}

func NewSignUpRequestHandler(
	logger *zap.Logger,
	userRepo user.UserRepository,
	mailClient *mail.Client,
) *SignUpRequestHandler {
	return &SignUpRequestHandler{
		logger:     logger,
		mailClient: mailClient,
		userRepo:   userRepo,
	}
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

	if err := h.userRepo.SaveNotifiedSignUpReqTime(ctx, userModel.UID, time.Now()); err != nil {
		return err
	}

	return nil
}

func (h *SignUpRequestHandler) sendNotification(ctx context.Context, userModel *user.User) error {
	body := fmt.Sprintf(
		`Hello dear %s!
Here's your confirmation code: %s`,
		userModel.Name,
		userModel.SignUpReq.Code,
	)

	if err := h.mailClient.SendTxt(ctx, mail.SendPayload{
		To:   []string{userModel.Email},
		Body: body,
	}); err != nil {
		return err
	}

	return nil
}
