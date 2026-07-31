package ntfy

import (
	"time"
)

type Config struct {
	BaseURL string        `env:"BASE_URL"`
	Token   string        `env:"TOKEN"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"3s"`
}
