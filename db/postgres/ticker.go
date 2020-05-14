package postgres

import (
	"github.com/trustwallet/watchmarket/db/models"
)

func (i *Instance) AddTickers(tickers []models.Ticker) error {
	// TODO: Upsert
	db := i.Gorm.Set("gorm:insert_option", "ON CONFLICT (coin,coin_name,coin_type,token_id,currency,provider) DO UPDATE SET value = excluded.value, change24h = excluded.change24h")
	return BulkInsert(db, tickers)
}

func (i *Instance) GetTickersByQueries(tickerQueries []models.TickerQuery) ([]models.Ticker, error) {
	var ticker []models.Ticker
	db := i.Gorm
	for _, tq := range tickerQueries {
		db = db.Or("coin = ? AND token_id = ?", tq.Coin, tq.TokenId)
	}
	if err := db.Find(&ticker).Error; err != nil {
		return nil, err
	}
	return ticker, nil
}

func (i *Instance) GetTickers(coin uint, tokenId string) ([]models.Ticker, error) {
	var ticker []models.Ticker
	if err := i.Gorm.Where("coin = ? AND token_id = ?", coin, tokenId).
		Find(&ticker).Error; err != nil {
		return nil, err
	}
	return ticker, nil
}
