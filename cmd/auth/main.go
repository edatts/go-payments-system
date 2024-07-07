package main

import (
	"log"

	"github.com/edatts/go-payment-system/pkg/api"
	"github.com/edatts/go-payment-system/pkg/api/auth"
	"github.com/edatts/go-payment-system/pkg/config"
	"github.com/edatts/go-payment-system/pkg/db"
)

func main() {

	// connCfg, err := pgx.ParseConfig("")
	// if err != nil {
	// 	log.Fatalf("failed parsing db connection address: %s", err)
	// }

	if config.Envs.JWT_SECRET == "" {
		panic("No JWT secret provided, base64 encoded JWT secret must be provided.")
	}

	dbCfg := config.GetDBConfig()

	// TODO: Replace this will a connection pool later.
	dbConn, err := db.NewPostgresStorage(dbCfg)
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

	server := api.NewServer(":4000", auth.NewHandler(auth.NewStore(dbConn)))
	if err := server.Run(); err != nil {
		log.Fatalf("auth server error: %s", err)
	}

}
