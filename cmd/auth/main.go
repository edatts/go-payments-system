package main

import (
	"log"

	"github.com/edatts/go-payment-system/pkg/api"
	"github.com/edatts/go-payment-system/pkg/api/auth"
	"github.com/edatts/go-payment-system/pkg/config"
	"github.com/edatts/go-payment-system/pkg/db"
	"github.com/edatts/go-payment-system/pkg/store"
)

func main() {
	if config.Envs.JWT_SECRET == "" {
		panic("No JWT secret provided, base64 encoded JWT secret must be provided.")
	}

	dbCfg := config.GetDBConfig()

	dbPool, err := db.NewPostgresStorage(dbCfg)
	if err != nil {
		log.Fatalf("db error: %s", err)
	}

	// TODO: Implement graceful shutdown for api servers.
	internalServer := api.NewServer(":4444", auth.NewInternalHandler())
	go func() {
		if err := internalServer.Run(); err != nil {
			log.Fatalf("internal auth server error: %s", err)
		}
	}()

	server := api.NewServer(":4000", auth.NewHandler(store.NewUserStore(dbPool)))
	if err := server.Run(); err != nil {
		log.Fatalf("auth server error: %s", err)
	}
}
