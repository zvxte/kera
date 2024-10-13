package database

import (
	"testing"
)

func TestNewMigrationInvalid(t *testing.T) {
	expectedErrorMessage := "Expected an error, but got nil"

	_, err := newMigration("")
	if err == nil {
		t.Error(expectedErrorMessage)
	}

	_, err = newMigration("00.sql")
	if err == nil {
		t.Error(expectedErrorMessage)
	}

	_, err = newMigration("000.txt")
	if err == nil {
		t.Error(expectedErrorMessage)
	}

	_, err = newMigration("000")
	if err == nil {
		t.Error(expectedErrorMessage)
	}
}

func TestNewMigrationValid(t *testing.T) {
	_, err := newMigration("000.sql")
	if err != nil {
		t.Error(err)
	}

	_, err = newMigration("999_migration.sql")
	if err != nil {
		t.Error(err)
	}
}
