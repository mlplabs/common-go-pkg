package workers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type WorkerMuConfig struct {
	Name           string        `env:"WORKER_NAME"`
	Enabled        bool          `env:"WORKER_ENABLED" envDefault:"false"`
	LockKey        string        `env:"WORKER_LOCK_KEY"`
	UniqueId       string        `env:"WORKER_UNIQUE_ID" envDefault:""`
	Interval       time.Duration `env:"WORKER_INTERVAL" envDefault:"300s"`
	AutoReleaseTTL time.Duration `env:"WORKER_RELEASE_TTL" envDefault:"200s"`
	LockTimeout    time.Duration `env:"WORKER_LOCK_TIMEOUT" envDefault:"3s"`
}

type WorkerMu struct {
	cfg *WorkerMuConfig
	rc  *redis.Client
	Worker
}

func NewWorkerMu(cfg *WorkerMuConfig, rc *redis.Client) WorkerMu {
	if cfg != nil && cfg.UniqueId == "" {
		cfg.UniqueId = uuid.New().String()
	}
	return WorkerMu{
		cfg: cfg,
		rc:  rc,
	}
}

func (w *WorkerMu) GetCfg() *WorkerMuConfig {
	return w.cfg
}

func (w *WorkerMu) Lock(ctx context.Context) (bool, error) {
	rCtx, cancel := context.WithTimeout(ctx, w.cfg.LockTimeout)
	defer cancel()

	lockSuccess, err := w.rc.SetNX(rCtx, w.cfg.LockKey, w.cfg.UniqueId, w.cfg.AutoReleaseTTL).Result()
	if err != nil {
		return false, fmt.Errorf("%s - redis.SetNX: %w", w.cfg.Name, err)
	}

	return lockSuccess, nil
}

func (w *WorkerMu) Release(ctx context.Context) (bool, error) {
	val, err := w.rc.Get(ctx, w.cfg.LockKey).Result()
	if err != nil {
		return false, fmt.Errorf("%s - redis.Get: %w", w.cfg.Name, err)
	}
	if val == w.cfg.UniqueId {
		res, err := w.rc.Del(ctx, w.cfg.LockKey).Result()
		if err != nil {
			return res == 1, fmt.Errorf("%s - redis.Del: %w", w.cfg.Name, err)
		}
		return res == 1, nil
	}
	return false, fmt.Errorf("%s - redis.Del: wrong key owner", w.cfg.Name)
}
