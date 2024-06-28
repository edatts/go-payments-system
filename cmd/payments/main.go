package payments

import (
	"log"

	"github.com/edatts/go-payment-system/pkg/api"
	"github.com/edatts/go-payment-system/pkg/api/payments"
	"github.com/edatts/go-payment-system/pkg/config"
	"github.com/edatts/go-payment-system/pkg/db"
)

func main() {

	dbCfg := config.GetDBConfig()

	// TODO: Replace this will a connection pool later.
	dbConn, err := db.NewPostgresStorage(dbCfg)
	if err != nil {
		log.Fatalf("db error: %s", err)
	}

	server := api.NewServer("localhost:4000", payments.NewHandler(payments.NewStore(dbConn)))
	if err := server.Run(); err != nil {
		log.Fatalf("auth server error: %s", err)
	}

}
