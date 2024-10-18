package database

import (
	"path/filepath"
	"testing"
)

func TestNewMigration(t *testing.T) {
	expectedErrorMessage := "expected an error, but got nil"

	validFilePath := filepath.Join("migrations", "000_create_migrations.sql")
	if _, err := newMigration(validFilePath); err != nil {
		t.Error(err)
	}

	if _, err := newMigration(""); err == nil {
		t.Error(expectedErrorMessage)
	}

	if _, err := newMigration(filepath.Join("migrations", "00.sql")); err == nil {
		t.Error(expectedErrorMessage)
	}

	if _, err := newMigration(filepath.Join("migrations", "000.txt")); err == nil {
		t.Error(expectedErrorMessage)
	}

	if _, err := newMigration(filepath.Join("migrations", "000")); err == nil {
		t.Error(expectedErrorMessage)
	}
}

func TestGetMigrations(t *testing.T) {
	if _, err := getMigrations(); err != nil {
		t.Error(err)
	}
}
