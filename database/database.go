package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const PostgresDriverName = "pgx"

type SqlDatabase struct {
	DB *sql.DB
}

// NewSqlDatabase returns a pointer to new SqlDatabase instance.
// This function checks if supported driver name is provided,
// and pings the database to validate given data source name.
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

// Setup sets up database migrations.
func (sd *SqlDatabase) Setup(ctx context.Context) error {
	migrations, err := getMigrations()
	if err != nil {
		return err
	}

	// Migrations table exists if SqlDatabase.getDatabaseMigrationVersion() method does not return an error.
	// In this case we want to exclude migrations with lower or equal version.

	latestMigrationVersion := migrations[len(migrations)-1].version
	databaseMigrationVersion, err := sd.getDatabaseMigrationVersion(ctx)
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

	// Now we can execute remaining migrations in a single transaction.

	tx, err := sd.DB.BeginTx(ctx, nil)
	defer tx.Rollback()

	for _, migration := range migrations {
		_, err = tx.ExecContext(ctx, migration.query)
		if err != nil {
			return fmt.Errorf("failed to execute %q: %w", migration.query, err)
		}
	}

	query := `
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
		return fmt.Errorf("failed to commit migrations transaction: %w", err)
	}

	return nil
}

func (sd *SqlDatabase) getDatabaseMigrationVersion(ctx context.Context) (uint16, error) {
	query := `
		SELECT version FROM migrations;
	`
	row := sd.DB.QueryRowContext(ctx, query)

	var databaseMigrationVersion uint16
	err := row.Scan(&databaseMigrationVersion)
	if err != nil {
		return databaseMigrationVersion, fmt.Errorf("failed to get migrations version: %w", err)
	}
	return databaseMigrationVersion, nil
}

// Teardown drops database migrations.
func (sd *SqlDatabase) Teardown(ctx context.Context) error {
	query := `
	DROP SCHEMA public CASCADE;
	CREATE SCHEMA public;
	`
	_, err := sd.DB.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to teardown: %w", err)
	}

	query = `
	SELECT tablename
	FROM pg_tables
	WHERE schemaname = 'public';
	`
	err = sd.DB.QueryRowContext(ctx, query).Scan()
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to teardown: %w", err)
	}

	return nil
}
