package signup

import (
	"context"
	"time"

	"go.uber.org/zap"

	userStore "apart-deal-api/pkg/store/user"
)

const (
	SignUpReqExpiration = time.Minute * 15
)

type ObsoleteReqWorker struct {
	logger   *zap.Logger
	userRepo userStore.UserRepository
}

func NewObsoleteReqWorker(userRepo userStore.UserRepository, logger *zap.Logger) *ObsoleteReqWorker {
	return &ObsoleteReqWorker{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (w *ObsoleteReqWorker) Process(ctx context.Context) error {
	deleted, err := w.userRepo.DeleteAllPendingOlderThan(ctx, time.Now().Truncate(SignUpReqExpiration))
	if err != nil {
		return err
	}

	w.logger.With(zap.Int("count", deleted)).Info("Deleted sign up req")

	return nil
}
