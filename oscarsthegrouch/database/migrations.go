package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func RunMigrations() error {
	migrationsFolderName := "./data/migrations"

	entries, err := os.ReadDir(migrationsFolderName)
	if err != nil {
		return err
	}

	conn, err := sql.Open("sqlite3", "./oscarsthegrouch.db")
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, migrationFile := range entries {
		conn.Exec("BEGIN TRANSACTION")
		migrationFile, err := os.ReadFile(migrationsFolderName + "/" + migrationFile.Name())
		if err != nil {
			return err
		}
		conn.Exec(string(migrationFile))
		conn.Exec("COMMIT")
	}

	return nil
}
