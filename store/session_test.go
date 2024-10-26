package store

import (
	"database/sql"
	"os"
	"testing"

	"github.com/zvxte/kera/database"
)

func TestNewSqlSessionStore(t *testing.T) {
	dataSourceName := os.Getenv("DSN")
	if dataSourceName == "" {
		t.Skip("skipping: DSN is not set")
	}

	db, err := sql.Open(database.PostgresDriverName, dataSourceName)
	if err != nil {
		t.Errorf("failed to open database: %v", err)
	}

	tests := []struct {
		name      string
		db        *sql.DB
		shouldErr bool
	}{
		{"Valid", db, false},
		{"Invalid nil", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewSqlSessionStore(test.db)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"NewSqlSessionStore(%v), error=%v, shouldErr=%v",
					test.db, err, test.shouldErr,
				)
			}
		})
	}
}
