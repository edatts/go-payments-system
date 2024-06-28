package main

import (
	"log"
	"os"

	"github.com/edatts/go-payment-system/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	pgxMigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {

	pgxCfg, err := pgx.ParseConfig(config.GetDBConfig().PostgresURL())
	if err != nil {
		log.Fatalf("faild parsing postgres url: %s", err)
	}

	pgxDriver, err := pgxMigrate.WithInstance(stdlib.OpenDB(*pgxCfg), &pgxMigrate.Config{})
	if err != nil {
		log.Fatalf("failed instantiating pgx driver: %s", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "pgx", pgxDriver)
	if err != nil {
		log.Fatalf("failed creating migration: %s", err)
	}

	switch os.Args[len(os.Args)-1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatalf("failed running migration: %s", err)
		}
	case "down":
		if err := m.Down(); err != nil {
			log.Fatalf("failed running migration: %s", err)
		}
	default:
		log.Fatal("invalid command, please provide either \"up\" or \"down\"")
	}

}
