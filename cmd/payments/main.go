package main

import (
	"log"

	"github.com/edatts/go-payment-system/pkg/api"
	"github.com/edatts/go-payment-system/pkg/api/payments"
	"github.com/edatts/go-payment-system/pkg/config"
	"github.com/edatts/go-payment-system/pkg/db"
	"github.com/edatts/go-payment-system/pkg/store"
)

func main() {

	dbCfg := config.GetDBConfig()

	// dbConn, err := db.NewPostgresStorage(dbCfg)
	// if err != nil {
	// 	log.Fatalf("db error: %s", err)
	// }

	dbConnPool, err := db.NewPostgresStorage(dbCfg)
	if err != nil {
		log.Fatalf("failed instantiating postgres storage: %s", err)
	}

	server := api.NewServer(":4001", payments.NewHandler(store.NewStore(dbConnPool)).Init())
	if err := server.Run(); err != nil {
		log.Fatalf("auth server error: %s", err)
	}

}
