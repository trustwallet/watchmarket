package postgres

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/trustwallet/watchmarket/db/models"
)

type Instance struct {
	Gorm *gorm.DB
}

const batchCount = 3000

func New(uri string) (*Instance, error) {
	g, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	g.AutoMigrate(
		&models.Rate{},
		&models.Ticker{},
	)
	//g.LogMode(true)

	i := &Instance{Gorm: g}

	return i, nil
}
