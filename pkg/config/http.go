package config

import "time"

type HTTP struct {
	Host         string        `env:"HTTP_HOST" envDefault:"localhost"`
	Port         string        `env:"HTTP_PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
}
