package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/db/models"
)

type Instance struct {
	Gorm *gorm.DB
}

func New(uri string) (*Instance, error) {
	g, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	g.AutoMigrate(
		&models.Subscription{},
	)

	i := &Instance{Gorm: g}

	return i, nil
}

func ConvertToError(errorsList []error) error {
	var errorsInfo string
	for _, e := range errorsList {
		errorsInfo += " " + e.Error()
	}
	return errors.E(errorsInfo)
}
