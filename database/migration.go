package database

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

const migrationsAssetsDir = "migrations"

type migration struct {
	version uint16
	query   string
}

func newMigration(filePath string) (migration, error) {
	fileName := filepath.Base(filePath)
	if len(fileName) < 3 || filepath.Ext(strings.ToLower(fileName)) != ".sql" {
		return migration{}, fmt.Errorf("invalid migration file name: %q", fileName)
	}

	version, err := strconv.ParseUint(fileName[:3], 10, 16)
	if err != nil {
		return migration{}, fmt.Errorf("failed to parse migration version: %q", fileName)
	}

	content, err := assets.ReadFile(filePath)
	if err != nil {
		return migration{}, fmt.Errorf("failed to read %q: %w", fileName, err)
	}

	return migration{version: uint16(version), query: string(content)}, nil
}

// getMigrations returns all migrations found in a given directory path sorted by file name.
func getMigrations() ([]migration, error) {
	entries, err := assets.ReadDir(migrationsAssetsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrations := make([]migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		entryPath := filepath.Join(migrationsAssetsDir, entry.Name())
		migration, err := newMigration(entryPath)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}
