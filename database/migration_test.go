package database

import (
	"path/filepath"
	"testing"
)

func TestNewMigration(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		shouldErr bool
	}{
		{
			"Valid",
			filepath.Join(migrationsAssetsDir, "000_create_migrations.sql"),
			false,
		},
		{
			"Invalid: directory",
			filepath.Join("migr", "000_create_migrations.sql"),
			true,
		},
		{
			"Invalid: empty directory",
			filepath.Join("", "000_create_migrations.sql"),
			true,
		},
		{
			"Invalid: file name",
			filepath.Join(migrationsAssetsDir, "00_create_migrations.sql"),
			true,
		},
		{
			"Invalid: file name",
			filepath.Join(migrationsAssetsDir, "000.sql"),
			true,
		},
		{
			"Invalid: empty file name",
			filepath.Join(migrationsAssetsDir, ""),
			true,
		},
		{
			"Invalid: file extension",
			filepath.Join(migrationsAssetsDir, "000_create_migrations.sq"),
			true,
		},
		{
			"Invalid: empty file extension",
			filepath.Join(migrationsAssetsDir, "000_create_migrations"),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := newMigration(test.filePath)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"newMigration(%q), error=%v, shouldErr=%v",
					test.filePath, err, test.shouldErr,
				)
			}
		})
	}
}

func TestGetMigrations(t *testing.T) {
	if _, err := getMigrations(); err != nil {
		t.Error(err)
	}
}
