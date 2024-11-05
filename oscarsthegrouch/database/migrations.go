package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func RunMigrations() error {
	conn, err := sql.Open("sqlite3", "./oscarsthegrouch.db")
	if err != nil {
		return err
	}
	defer conn.Close()

	migrationsFolderName := "./data/migrations"

	entries, err := os.ReadDir(migrationsFolderName)
	if err != nil {
		return err
	}

	migrationsStart, err := getMigrationStart(conn)
	if err != nil {
		return err
	}

	for i, migrationFile := range entries[migrationsStart:] {
		conn.Exec("BEGIN TRANSACTION")
		migration, err := os.ReadFile(migrationsFolderName + "/" + migrationFile.Name())
		if err != nil {
			return err
		}
		conn.Exec(string(migration))
		conn.Exec("COMMIT")

		conn.Exec("UPDATE schema_migrations SET version = ?;", migrationsStart+i+1)
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
