package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovrivero/rinha-backend/internal/config"
)

var Connection *pgxpool.Pool

func Connect() error {
	var err error
	config, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return err
	}
	Connection, err = pgxpool.NewWithConfig(context.Background(), config)

	return err
}

func Close() {
	if Connection == nil {
		return
	}
	Connection.Close()
}
