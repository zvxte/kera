package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type migration struct {
	version  uint16
	fileName string
}

func newMigration(fileName string) (migration, error) {
	if len(fileName) < 3 || filepath.Ext(strings.ToLower(fileName)) != ".sql" {
		return migration{}, fmt.Errorf("invalid migration file name: %q", fileName)
	}

	version, err := strconv.ParseUint(fileName[:3], 10, 16)
	if err != nil {
		return migration{}, fmt.Errorf("failed to parse migration version: %q", fileName)
	}

	return migration{version: uint16(version), fileName: fileName}, nil
}

// getMigrations returns all migrations found in a given directory path sorted by file name.
func getMigrations(dirPath string) ([]migration, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrations := make([]migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		migration, err := newMigration(entry.Name())
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}
	return migrations, nil
}
