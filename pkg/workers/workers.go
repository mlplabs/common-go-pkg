package workers

import (
	"context"
)

type Worker interface {
	Do(ctx context.Context)
}

type Workers struct {
	workers []Worker
}

func NewWorkers(workers ...Worker) *Workers {
	return &Workers{workers: workers}
}

func (w *Workers) Start(ctx context.Context) {
	for _, work := range w.workers {
		go work.Do(ctx)
	}
}
