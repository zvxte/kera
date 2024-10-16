package database

import (
	"path/filepath"
	"testing"
)

func TestNewMigrationInvalid(t *testing.T) {
	expectedErrorMessage := "Expected an error, but got nil"

	_, err := newMigration("")
	if err == nil {
		t.Error(expectedErrorMessage)
	}

	_, err = newMigration(filepath.Join("migrations", "00.sql"))
	if err == nil {
		t.Error(expectedErrorMessage)
	}

	_, err = newMigration(filepath.Join("migrations", "000.txt"))
	if err == nil {
		t.Error(expectedErrorMessage)
	}

	_, err = newMigration(filepath.Join("migrations", "000"))
	if err == nil {
		t.Error(expectedErrorMessage)
	}
}

func TestNewMigrationValid(t *testing.T) {
	validFilePath := filepath.Join("migrations", "000_create_migrations.sql")
	_, err := newMigration(validFilePath)
	if err != nil {
		t.Error(err)
	}
}

func TestGetMigrations(t *testing.T) {
	_, err := getMigrations()
	if err != nil {
		t.Error(err)
	}
}
