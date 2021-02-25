package postgres

import (
	"github.com/trustwallet/watchmarket/services/controllers"
	"strconv"
	"strings"

	"github.com/trustwallet/watchmarket/db/models"
	"gorm.io/gorm/clause"
)

func (i *Instance) AddTickers(tickers []models.Ticker) error {
	for _, b := range normalizeTickers(tickers) {
		err := i.Gorm.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "coin",
				},
				{
					Name: "coin_name",
				},
				{
					Name: "coin_type",
				},
				{
					Name: "token_id",
				},
				{
					Name: "currency",
				},
				{
					Name: "provider",
				},
			},
			DoUpdates: clause.AssignmentColumns([]string{"value", "change24h", "volume", "total_supply", "circulating_supply", "market_cap", "last_updated", "updated_at"}),
		}).Create(&b).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func normalizeTickers(tickers []models.Ticker) []models.Ticker {
	tickersMap := make(map[string]models.Ticker)
	for _, ticker := range tickers {
		key := strconv.Itoa(int(ticker.Coin)) +
			ticker.CoinName + ticker.CoinType +
			ticker.TokenId + ticker.Currency +
			ticker.Provider
		if _, ok := tickersMap[key]; !ok {
			tickersMap[key] = ticker
		}
	}
	result := make([]models.Ticker, 0, len(tickersMap))
	for _, ticker := range tickersMap {
		result = append(result, ticker)
	}
	return result
}

func (i *Instance) GetTickers(assets []controllers.Asset) ([]models.Ticker, error) {
	var ticker []models.Ticker
	db := i.Gorm
	for _, asset := range assets {
		db = db.Or("coin = ? AND token_id = ?", asset.CoinId, strings.ToLower(asset.TokenId))
	}
	if err := db.Find(&ticker).Error; err != nil {
		return nil, err
	}
	return ticker, nil
}

func (i *Instance) GetAllTickers() ([]models.Ticker, error) {
	var tickers []models.Ticker
	if err := i.Gorm.Find(&tickers).Error; err != nil {
		return nil, err
	}
	return tickers, nil
}
