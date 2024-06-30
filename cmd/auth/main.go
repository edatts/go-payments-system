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

	dbCfg := config.GetDBConfig()

	// TODO: Replace this will a connection pool later.
	dbConn, err := db.NewPostgresStorage(dbCfg)
	if err != nil {
		log.Fatalf("db error: %s", err)
	}

	server := api.NewServer(":4000", auth.NewHandler(auth.NewStore(dbConn)))
	if err := server.Run(); err != nil {
		log.Fatalf("auth server error: %s", err)
	}

}
