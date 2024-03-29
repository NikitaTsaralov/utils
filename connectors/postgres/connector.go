package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// New returns new Postgresql db instance with sqlx driver.
func New(c Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		c.Host,
		c.Port,
		c.User,
		c.DBName,
		c.SSLMode,
		c.Password,
	)

	db, err := sqlx.Connect(c.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.Settings.MaxOpenConns)
	db.SetConnMaxLifetime(c.Settings.ConnMaxLifetime * time.Second)
	db.SetMaxIdleConns(c.Settings.MaxIdleConns)
	db.SetConnMaxIdleTime(c.Settings.ConnMaxIdleTime * time.Second)

	return db, nil
}
