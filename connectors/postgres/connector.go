package postgres

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Start returns new Postgresql db instance with sqlx driver.
func Start(c Config) (*sqlx.DB, error) {
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
		return nil, err
	}

	db.SetMaxOpenConns(c.Settings.MaxOpenConns)
	db.SetConnMaxLifetime(c.Settings.ConnMaxLifetime * time.Millisecond)
	db.SetMaxIdleConns(c.Settings.MaxIdleConns)
	db.SetConnMaxIdleTime(c.Settings.ConnMaxIdleTime * time.Millisecond)

	return db, nil
}
