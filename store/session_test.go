package store

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/zvxte/kera/database"
)

func TestNewSqlSessionStore(t *testing.T) {
	dataSourceName := os.Getenv("DSN")
	if dataSourceName == "" {
		t.Skip("skipping: DSN is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.NewSqlDatabase(ctx, database.PostgresDriverName, dataSourceName)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name      string
		db        *sql.DB
		shouldErr bool
	}{
		{"Valid", db.DB, false},
		{"Invalid: nil sql.DB pointer", nil, true},
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
