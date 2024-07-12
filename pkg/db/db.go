package db

import (
	"context"
	"fmt"

	"github.com/edatts/go-payment-system/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresStorage(dbCfg config.DBConfig) (*pgxpool.Pool, error) {
	// postgresql://user:secret@host:5432/dbName
	dbUrl := dbCfg.PostgresURL()

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return &pgxpool.Pool{}, fmt.Errorf("failed connecting to pg instance: %w", err)
	}

	return pool, nil
}
