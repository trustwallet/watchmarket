package db

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
	"github.com/trustwallet/watchmarket/db/models"
	"reflect"
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

	i := &Instance{Gorm: g}

	return i, nil
}

// postgress.BulkInsert(Instance.Gorm, []models.Ticker{...})
func BulkInsert(db *gorm.DB, dbModels interface{}) error {
	interfaceSlice, err := getInterfaceSlice(dbModels)
	if err != nil {
		return err
	}
	batchList := getInterfaceSliceBatch(interfaceSlice, batchCount)
	for _, batch := range batchList {
		err := gormbulk.BulkInsert(db, batch, len(batch))
		if err != nil {
			return err
		}
	}
	return nil
}
func getInterfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}

func getInterfaceSliceBatch(values []interface{}, sizeUint uint) [][]interface{} {
	size := int(sizeUint)
	resultLength := (len(values) + size - 1) / size
	result := make([][]interface{}, resultLength)
	lo, hi := 0, size
	for i := range result {
		if hi > len(values) {
			hi = len(values)
		}
		result[i] = values[lo:hi:hi]
		lo, hi = hi, hi+size
	}
	return result
}
