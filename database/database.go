package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DriverName uint8

func (driverName DriverName) String() string {
	switch driverName {
	case PostgresDriverName:
		return "pgx"
	default:
		return ""
	}
}

const (
	InvalidDriverName DriverName = iota
	PostgresDriverName
)

const (
	ConnectionTimeout = 5 * time.Second
)

type SqlDatabase struct {
	driverName DriverName
	DB         *sql.DB
}

func NewSqlDatabase(driverName DriverName, dataSourceName string) (*SqlDatabase, error) {
	switch driverName {
	case PostgresDriverName:
	// Supported
	default:
		return nil, fmt.Errorf("unsupported driver name %q", driverName.String())
	}

	db, err := sql.Open(driverName.String(), dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &SqlDatabase{driverName: driverName, DB: db}, nil
}
