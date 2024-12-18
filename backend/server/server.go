package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zvxte/kera/database"
	"github.com/zvxte/kera/server/handler"
	"github.com/zvxte/kera/store/habitstore"
	"github.com/zvxte/kera/store/sessionstore"
	"github.com/zvxte/kera/store/userstore"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() (*Server, error) {
	logger := log.Default()

	dataSourceName := os.Getenv("DSN")
	if dataSourceName == "" {
		return nil, errors.New("failed to create Server: DSN is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlDatabase, err := database.NewSqlDatabase(
		ctx, database.PostgresDriverName, dataSourceName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Server: %w", err)
	}
	err = sqlDatabase.Setup(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Server: %w", err)
	}

	userStore, err := userstore.NewSql(sqlDatabase.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create Server: %w", err)
	}

	sessionStore, err := sessionstore.NewSql(sqlDatabase.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create Server: %w", err)
	}

	habitStore, err := habitstore.NewSql(sqlDatabase.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create Server: %w", err)
	}

	authMux := handler.NewAuthMux(userStore, sessionStore, logger)
	meMux := handler.NewMeMux(userStore, sessionStore, logger)
	habitsMux := handler.NewHabitsMux(habitStore, userStore, logger)

	mux := http.NewServeMux()
	mux.Handle("/auth/", http.StripPrefix("/auth", authMux))
	mux.Handle("/me/", handler.SessionMiddleware(
		http.StripPrefix("/me", meMux), sessionStore),
	)
	mux.Handle("/habits/", handler.SessionMiddleware(
		http.StripPrefix("/habits", habitsMux), sessionStore),
	)
	return &Server{mux: mux}, nil
}

func (server *Server) Run(address string) {
	http.ListenAndServe(address, server.mux)
}
