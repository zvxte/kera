package database

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNewSqlDatabase(t *testing.T) {
	empty := ""
	expectedErrorMessage := "expected an error, but got nil"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := NewSqlDatabase(ctx, empty, empty); err == nil {
		t.Error(expectedErrorMessage)
	}

	if _, err := NewSqlDatabase(ctx, PostgresDriverName, empty); err == nil {
		t.Error(expectedErrorMessage)
	}

	dataSourceName := os.Getenv("DSN")
	if dataSourceName == "" {
		t.Skip("skipping: DSN is not set")
	}

	if _, err := NewSqlDatabase(ctx, PostgresDriverName, dataSourceName); err != nil {
		t.Error(err)
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
