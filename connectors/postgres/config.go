package postgres

import (
	"time"
)

type Config struct {
	Host           string
	Port           int
	User           string
	Password       string `validate:"required"`
	DBName         string
	SSLModeEnabled bool
	Driver         string
	Settings       Settings
}

// FillWithDefaults - sets default values for empty fields
func (c *Config) FillWithDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}

	if c.Port == 0 {
		c.Port = 5432
	}

	if c.User == "" {
		c.User = "root"
	}

	if c.DBName == "" {
		c.DBName = "postgres"
	}

	if c.Driver == "" {
		c.Driver = "pgx"
	}

	c.Settings.FillWithDefaults()
}

type Settings struct {
	MaxOpenConns    int           `validate:"min=1"`
	ConnMaxLifetime time.Duration `validate:"min=1"`
	MaxIdleConns    int           `validate:"min=1"`
	ConnMaxIdleTime time.Duration `validate:"min=1"`
}

func (s *Settings) FillWithDefaults() {
	if s.MaxOpenConns == 0 {
		s.MaxOpenConns = 25
	}

	if s.ConnMaxLifetime == 0 {
		s.ConnMaxLifetime = 300000
	}

	if s.MaxIdleConns == 0 {
		s.MaxIdleConns = 25
	}

	if s.ConnMaxIdleTime == 0 {
		s.ConnMaxIdleTime = 300000
	}
}
