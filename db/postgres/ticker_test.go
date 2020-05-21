package postgres

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/db/models"
	"regexp"
	"testing"
)

func TestInstance_AddTickers(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()
	mock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO "tickers"`)).WithArgs(float64(60), uint(60), "60", "60", sqlmock.AnyArg(), "60", "60", "60", sqlmock.AnyArg(), float64(60)).WillReturnResult(sqlmock.NewResult(1, 1))
	i := Instance{Gorm: db}

	assert.Nil(t, i.AddTickers([]models.Ticker{{
		Coin:      60,
		CoinName:  "60",
		CoinType:  "60",
		TokenId:   "60",
		Change24h: 60,
		Currency:  "60",
		Provider:  "60",
		Value:     60,
	},
	}))
}

func setupDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when sqlmock", err)
	}

	d, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	d.LogMode(true)
	return d, mock
}
