package database

import (
	"context"
	"database/sql"

	"namaya/oscarsthegrouch/model"

	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func ConnectDb(ctx context.Context) (*sql.DB, error) {
	dbClient, err := sql.Open("sqlite3", "./oscarsthegrouch.db")
	if err != nil {
		return nil, err
	}

	driver, err := sqlite.WithInstance(dbClient, &sqlite.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://data/migrations", "sqlite", driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return dbClient, nil
}

// TODO: each game needs a separate master ballot
func GetMasterBallot(ctx context.Context) (*model.Ballot, error) {
	return nil, nil
}
