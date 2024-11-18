package database

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNewSqlDatabase(t *testing.T) {
	dataSourceName := os.Getenv("DSN")
	if dataSourceName == "" {
		t.Skip("skipping: DSN is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name           string
		driverName     string
		dataSourceName string
		shouldErr      bool
	}{
		{
			"Valid",
			PostgresDriverName,
			dataSourceName,
			false,
		},
		{
			"Invalid: empty driver name",
			"",
			dataSourceName,
			true,
		},
		{
			"Invalid: empty data source name",
			PostgresDriverName,
			"",
			true,
		},
		{
			"Invalid: driver name",
			"sqlite",
			dataSourceName,
			true,
		},
		{
			"Invalid: data source name",
			PostgresDriverName,
			"postgres://username:password@localhost:5432/database",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewSqlDatabase(
				ctx, test.driverName, test.dataSourceName,
			)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"NewSqlDatabase(%v, %q, %q), error=%v, shouldErr=%v",
					ctx, test.driverName, test.dataSourceName, err, test.shouldErr,
				)
			}
		})
	}
}

func TestSqlDatabase(t *testing.T) {
	dataSourceName := os.Getenv("DSN")
	if dataSourceName == "" {
		t.Skip("skipping: DSN is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlDatabase, err := NewSqlDatabase(ctx, PostgresDriverName, dataSourceName)
	if err != nil {
		t.Error(err)
	}

	t.Run("Setup", func(t *testing.T) {
		if err := sqlDatabase.Setup(ctx); err != nil {
			t.Error(err)
		}
	})

	t.Run("getDatabaseMigrationVersion", func(t *testing.T) {
		if _, err := sqlDatabase.getDatabaseMigrationVersion(ctx); err != nil {
			t.Error(err)
		}
	})

	t.Run("Teardown", func(t *testing.T) {
		if err := sqlDatabase.Teardown(ctx); err != nil {
			t.Error(err)
		}
	})
}
