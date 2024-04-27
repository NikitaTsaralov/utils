package postgres

import (
	"fmt"
	"time"

	"github.com/NikitaTsaralov/utils/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// New returns new Postgresql db instance with sqlx driver.
func New(c Config) *sqlx.DB {
	c.FillWithDefaults()

	sslMode := "disable"
	if c.SSLModeEnabled {
		sslMode = "enable"
	}

	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		c.Host,
		c.Port,
		c.User,
		c.DBName,
		sslMode,
		c.Password,
	)

	db, err := sqlx.Connect(c.Driver, dataSourceName)
	if err != nil {
		logger.Fatalf("can't connect to postgres: %v", err)
	}

	db.SetMaxOpenConns(c.Settings.MaxOpenConns)
	db.SetConnMaxLifetime(c.Settings.ConnMaxLifetime * time.Millisecond)
	db.SetMaxIdleConns(c.Settings.MaxIdleConns)
	db.SetConnMaxIdleTime(c.Settings.ConnMaxIdleTime * time.Millisecond)

	return db
}
