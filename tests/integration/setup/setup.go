// +build integration

package setup

import (
	"github.com/trustwallet/watchmarket/db/postgres"
	"log"
)

func RunPgContainer() *postgres.Instance {
	dbConn, err := runPgContainerAndInitConnection()
	if err != nil {
		log.Fatal(err)
	}
	return dbConn
}

func StopPgContainer() {
	if err := stopPgContainer(); err != nil {
		log.Fatal(err)
	}
}
