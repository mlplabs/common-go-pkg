package redis

type Config struct {
	Host     string `env:"REDIS_HOST" envDefault:"127.0.0.1"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	Database int    `env:"REDIS_DB" envDefault:"0"`
}
