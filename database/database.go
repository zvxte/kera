package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	Postgres          = "pgx"
	ConnectionTimeout = 5 * time.Second
)

type SqlDatabase struct {
	DB *sql.DB
}

func NewSqlDatabase(driverName, dataSourceName string) (*SqlDatabase, error) {
	switch driverName {
	case Postgres:
	// Supported
	default:
		return nil, fmt.Errorf("unsupported driver name %q", driverName)
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &SqlDatabase{DB: db}, nil
}
