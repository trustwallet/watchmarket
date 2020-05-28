package postgres

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/trustwallet/watchmarket/db/models"
	"go.elastic.co/apm/module/apmgorm"
	_ "go.elastic.co/apm/module/apmgorm/dialects/postgres"
)

type Instance struct {
	Gorm *gorm.DB
}

func New(uri, env string, logMode bool) (*Instance, error) {
	var (
		g   *gorm.DB
		err error
	)

	if env == "prod" {
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
