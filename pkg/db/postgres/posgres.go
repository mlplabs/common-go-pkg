package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	_defaultMaxPoolSize           = 1
	_defaultMaxOpenConnections    = 3
	_defaultMaxConnectionIdleTime = 2 * time.Minute
)

type Config struct {
	PoolMax  int
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

// Connect - создаёт новое подключение к БД postgreSQL.
func Connect(cfg *Config) (*sql.DB, error) {
	postgresDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB)

	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		return nil, fmt.Errorf("connectDb. Failed to open db: %w", err)
	}

	//db.SetConnMaxIdleTime(_defaultMaxConnectionIdleTime)
	//db.SetMaxOpenConns(_defaultMaxOpenConnections)
	//db.SetMaxIdleConns(_defaultMaxPoolSize)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connectDb. Failed to ping db: %w", err)
	}

	return db, nil
}
