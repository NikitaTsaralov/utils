package postgres

import (
	"time"
)

type Config struct {
	Host           string `validate:"default=localhost"`
	Port           int    `validate:"default=5432"`
	User           string `validate:"default=root"`
	Password       string `validate:"required"`
	DBName         string `validate:"default=postgres"`
	SSLModeEnabled bool
	Driver         string `validate:"default=pgx"`
	Settings       Settings
}

type Settings struct {
	MaxOpenConns    int           `validate:"min=1,default=25"`
	ConnMaxLifetime time.Duration `validate:"min=1,default=5m"`
	MaxIdleConns    int           `validate:"min=1,default=25"`
	ConnMaxIdleTime time.Duration `validate:"min=1,default=5m"`
}
