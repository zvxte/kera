package database

import (
	"fmt"
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
