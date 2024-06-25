package db

import (
	"context"
	"fmt"

	"github.com/edatts/go-payment-system/pkg/config"
	"github.com/jackc/pgx/v5"
)

// func NewPostgresStorage(cfg *pgx.ConnConfig) (*pgx.Conn, error) {
func NewPostgresStorage(dbCfg config.DBConfig) (*pgx.Conn, error) {
	// postgresql://user:secret@host:5432/dbName
	dbUrl := dbCfg.PostgresURL()

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return &pgx.Conn{}, fmt.Errorf("failed connecting to pg instance: %w", err)
	}

	return conn, nil
}
