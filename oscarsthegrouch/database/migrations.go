package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDb(ctx context.Context) (*sql.DB, error) {
	dbClient, err := sql.Open("sqlite3", "./oscarsthegrouch.db")
	if err != nil {
		return nil, err
	}

	err = RunMigrations(ctx, dbClient)
	if err != nil {
		return nil, err
	}

	return dbClient, nil
}

func RunMigrations(ctx context.Context, dbClient *sql.DB) error {
	migrationsFolderName := "./data/migrations"

	entries, err := os.ReadDir(migrationsFolderName)
	if err != nil {
		return err
	}

	migrationsStart, err := getMigrationStart(dbClient)
	if err != nil {
		return err
	}

	for i, migrationFile := range entries[migrationsStart:] {
		migrationQuery, err := os.ReadFile(migrationsFolderName + "/" + migrationFile.Name())
		if err != nil {
			return err
		}

		result, err := dbClient.ExecContext(ctx, fmt.Sprintf(`
            BEGIN TRANSACTION;
            %s
            UPDATE schema_migrations SET version = ?;
        `, migrationQuery), migrationsStart+i+1)

		if err != nil {
			return err
		}

		rowsUpdated, err := result.RowsAffected()
		if err != nil {
			return err
		}

		slog.InfoContext(ctx, "Migration %s applied successfully. %s records updated.", migrationFile.Name(), rowsUpdated)
	}

	return nil
}

func getMigrationStart(conn *sql.DB) (int, error) {
	query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='schema_migrations';"

	var count int
	err := conn.QueryRow(query).Scan(&count)
	if err != nil {
		return -1, err
	}

	if count == 0 {
		return 0, nil
	}

	query = "SELECT version FROM schema_migrations;"

	var migrationNumber int
	err = conn.QueryRow(query).Scan(&migrationNumber)
	if err != nil {
		return -1, err
	}

	return migrationNumber, nil
}
