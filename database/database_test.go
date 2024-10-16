package database

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNewSqlDatabaseInvalid(t *testing.T) {
	empty := ""
	expectedErrorMessage := "Expected an error, but got nil"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := NewSqlDatabase(ctx, empty, empty)
	if err == nil {
		t.Error(expectedErrorMessage)
	}

	_, err = NewSqlDatabase(ctx, PostgresDriverName, empty)
	if err == nil {
		t.Error(expectedErrorMessage)
	}
}

func TestNewSqlDatabaseValid(t *testing.T) {
	dataSourceName := os.Getenv("DSN")
	if dataSourceName == "" {
		t.Skip("Skipping: DSN is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := NewSqlDatabase(ctx, PostgresDriverName, dataSourceName)
	if err != nil {
		t.Error(err)
	}
}

func TestSqlDatabaseSetup(t *testing.T) {
	dataSourceName := os.Getenv("DSN")
	if dataSourceName == "" {
		t.Skip("Skipping: DSN is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlDatabase, err := NewSqlDatabase(ctx, PostgresDriverName, dataSourceName)
	if err != nil {
		t.Error(err)
	}

	err = sqlDatabase.Setup(ctx)
	if err != nil {
		t.Error(err)
	}
}
