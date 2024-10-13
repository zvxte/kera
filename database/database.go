package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	PostgresDriverName = "pgx"
)

type SqlDatabase struct {
	DB *sql.DB
}

func NewSqlDatabase(ctx context.Context, driverName string, dataSourceName string) (*SqlDatabase, error) {
	switch driverName {
	case PostgresDriverName:
	// Supported
	default:
		return nil, fmt.Errorf("unsupported driver name: %q", driverName)
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &SqlDatabase{DB: db}, nil
}

func (sqlDatabase *SqlDatabase) Setup(ctx context.Context) error {
	migrationsDirPath := "migrations"

	entries, err := os.ReadDir(migrationsDirPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrations := make([]migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		migration, err := newMigration(entry.Name())
		if err != nil {
			return err
		}
		migrations = append(migrations, migration)
	}
	latestMigrationVersion := migrations[len(migrations)-1].version

	query := `
		SELECT version FROM migrations;
	`
	row := sqlDatabase.DB.QueryRowContext(ctx, query)

	var databaseMigrationVersion uint16
	err = row.Scan(&databaseMigrationVersion)
	if err == nil {
		if databaseMigrationVersion == latestMigrationVersion {
			return nil
		}
		if databaseMigrationVersion > latestMigrationVersion {
			return fmt.Errorf("database migration version %d is greater than latest migration version %d", databaseMigrationVersion, latestMigrationVersion)
		}

		for i, migration := range migrations {
			if migration.version >= databaseMigrationVersion {
				migrations = migrations[i+1:]
				break
			}
		}
	}

	tx, err := sqlDatabase.DB.BeginTx(ctx, nil)
	defer tx.Rollback()

	for _, migration := range migrations {
		content, err := os.ReadFile(fmt.Sprint(migrationsDirPath, "/", migration.fileName))
		if err != nil {
			return fmt.Errorf("failed to read %q: %w", migration.fileName, err)
		}
		query := string(content)

		_, err = tx.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to execute %q: %w", migration.fileName, err)
		}
	}

	query = `
		INSERT INTO migrations (id, version)
		VALUES (0, $1)
		ON CONFLICT (id)
		DO UPDATE SET version = $1;
	`
	_, err = tx.ExecContext(ctx, query, latestMigrationVersion)
	if err != nil {
		return fmt.Errorf("failed to update database migration version: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
