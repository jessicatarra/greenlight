package database

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jessicatarra/greenlight/assets"
	_ "github.com/lib/pq"

	"time"
)

const defaultTimeout = 3 * time.Second

func New(dsn string, maxOpenConns int, maxIdleConns int, maxIdleTime string, automigrate bool) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	if automigrate {
		iofsDriver, err := iofs.New(assets.MigrationsFiles, "migrations")
		if err != nil {
			return nil, err
		}

		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, dsn)
		if err != nil {
			return nil, err
		}

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, err
		}
	}

	return db, nil
}
