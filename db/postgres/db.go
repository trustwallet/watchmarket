package postgres

import (
	"github.com/trustwallet/watchmarket/db/models"
	"gorm.io/gorm/logger"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Instance struct {
	Gorm *gorm.DB
}

func New(url string, logMode bool) (*Instance, error) {
	var cfg *gorm.Config
	if logMode {
		cfg = &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	} else {
		cfg = &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	}

	db, err := gorm.Open(postgres.Open(url), cfg)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Rate{},
		&models.Ticker{},
	)
	if err != nil {
		log.Error(err)
	}

	i := &Instance{Gorm: db}

	return i, nil
}

func FatalWorker(timeout time.Duration, i Instance) {
	log.Info("Run PG RestoreConnectionWorker")
	for {
		db, err := i.Gorm.DB()
		if err != nil {
			panic("PG is not available now")
		}

		dbWriteErr := db.Ping()
		if dbWriteErr != nil {
			panic("PG is not available now")
		}
		time.Sleep(timeout)
	}
}
