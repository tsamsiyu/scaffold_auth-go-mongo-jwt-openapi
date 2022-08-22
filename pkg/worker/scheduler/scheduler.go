package scheduler

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Worker interface {
	Process(ctx context.Context) error
}

type workerInterval struct {
	contract Worker
	interval time.Duration
}

type Scheduler struct {
	logger          *zap.Logger
	workerIntervals []workerInterval
}

func NewScheduler(logger *zap.Logger) *Scheduler {
	return &Scheduler{
		logger:          logger,
		workerIntervals: make([]workerInterval, 0),
	}
}

func (s *Scheduler) Register(w Worker, interval time.Duration) {
	s.workerIntervals = append(s.workerIntervals, workerInterval{
		contract: w,
		interval: interval,
	})
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, wi := range s.workerIntervals {
		go func(wi workerInterval) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if err := wi.contract.Process(ctx); err != nil {
						s.logger.Error(errors.WithStack(err).Error())
					}
				}

				time.Sleep(wi.interval)
			}
		}(wi)
	}
}
