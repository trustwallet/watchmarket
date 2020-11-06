// +build integration

package setup

import (
	"fmt"
	"github.com/ory/dockertest"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/db/postgres"

	"gorm.io/gorm"

	"log"
)

const (
	pgUser = "user"
	pgPass = "pass"
	pgDB   = "watchmarket"
)

var (
	pgResource     *dockertest.Resource
	pgContainerENV = []string{
		"POSTGRES_USER=" + pgUser,
		"POSTGRES_PASSWORD=" + pgPass,
		"POSTGRES_DB=" + pgDB,
	}

	tables = []interface{}{
		&models.Ticker{},
		&models.Rate{},
	}

	uri string
)

func runPgContainerAndInitConnection() (*postgres.Instance, error) {
	pool := runPgContainer()
	var (
		dbConn *postgres.Instance
		err    error
	)
	if err := pool.Retry(func() error {
		dbConn, err = postgres.New(uri, false)
		return err
	}); err != nil {
		return nil, err
	}
	autoMigrate(dbConn.Gorm)

	return dbConn, nil
}

func CleanupPgContainer(dbConn *gorm.DB) {
	if err := dbConn.Migrator().DropTable(tables...); err != nil {
		log.Fatal(err)
	}
	autoMigrate(dbConn)
}

func autoMigrate(dbConn *gorm.DB) {
	dbConn.AutoMigrate(tables...)
}

func stopPgContainer() error {
	return pgResource.Close()
}

func runPgContainer() *dockertest.Pool {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	pgResource, err = pool.Run("postgres", "latest", pgContainerENV)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	uri = fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		pgUser, pgPass, pgResource.GetPort("5432/tcp"), pgDB,
	)
	return pool
}
