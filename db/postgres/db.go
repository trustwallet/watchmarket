package postgres

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db/models"
	"go.elastic.co/apm/module/apmgorm"
	_ "go.elastic.co/apm/module/apmgorm/dialects/postgres"
	"time"
)

type Instance struct {
	Gorm *gorm.DB
}

func New(uri string, apm, logMode bool) (*Instance, error) {
	var (
		g   *gorm.DB
		err error
	)

	if apm {
		g, err = apmgorm.Open("postgres", uri)
		if err != nil {
			return nil, err
		}
	} else {
		g, err = gorm.Open("postgres", uri)
		if err != nil {
			return nil, err
		}
	}

	g.LogMode(logMode)

	g.AutoMigrate(
		&models.Rate{},
		&models.Ticker{},
	)

	i := &Instance{Gorm: g}

	return i, nil
}

func FatalWorker(timeout time.Duration, i Instance) {
	logger.Info("Run PG RestoreConnectionWorker")
	for {
		dbWriteErr := i.Gorm.DB().Ping()
		if dbWriteErr != nil {
			panic("PG is not available now")
		}
		time.Sleep(timeout)
	}
}
