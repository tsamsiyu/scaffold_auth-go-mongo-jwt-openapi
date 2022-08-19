package signup

import (
	"context"
	"time"
)

type Scheduler struct {
	worker *Worker
}

func NewScheduler(w *Worker) *Scheduler {
	return &Scheduler{
		worker: w,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	stop := false

	go func() {
		<-ctx.Done()
		stop = true
	}()

	for !stop {
		if err := s.worker.Process(ctx); err != nil {
			return err
		}

		time.Sleep(time.Minute)
	}

	return nil
}
