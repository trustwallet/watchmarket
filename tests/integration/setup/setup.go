package setup

import (
	"github.com/trustwallet/watchmarket/db"
	"log"
)

func RunPgContainer() *db.Instance {
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
