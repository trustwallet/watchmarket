package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/trustwallet/watchmarket/db/models"
	"strings"
	"time"
)

const (
	batchLimit           = 3000
	rawBulkTickersInsert = `INSERT INTO tickers(updated_at,created_at,coin,coin_name,coin_type,token_id,change24h,currency,provider,value,last_updated,volume,market_cap,show_option) VALUES %s ON CONFLICT ON CONSTRAINT tickers_pkey DO UPDATE SET value = excluded.value, change24h = excluded.change24h, updated_at = excluded.updated_at, last_updated = excluded.last_updated, volume = excluded.volume, market_cap = excluded.market_cap`
)

func (i *Instance) AddTickers(tickers []models.Ticker) error {
	batch := toTickersBatch(normalizeTickers(tickers), batchLimit)
	for _, b := range batch {
		err := bulkCreateTicker(i.Gorm, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func toTickersBatch(tickers []models.Ticker, sizeUint uint) [][]models.Ticker {
	size := int(sizeUint)
	resultLength := (len(tickers) + size - 1) / size
	result := make([][]models.Ticker, resultLength)
	lo, hi := 0, size
	for i := range result {
		if hi > len(tickers) {
			hi = len(tickers)
		}
		result[i] = tickers[lo:hi:hi]
		lo, hi = hi, hi+size
	}
	return result
}

func normalizeTickers(tickers []models.Ticker) []models.Ticker {
	normalizedTickers := make([]models.Ticker, 0, len(tickers))
	for _, t := range tickers {
		if !isBadTicker(t.Coin, t.CoinName, t.CoinType, t.TokenId, t.Currency, t.Provider, t.Value, t.Change24h, tickers) {
			normalizedTickers = append(normalizedTickers, t)
		}
	}
	return toUniqueTickers(normalizedTickers)
}

func toUniqueTickers(sample []models.Ticker) []models.Ticker {
	var unique []models.Ticker
sampleLoop:
	for _, v := range sample {
		for i, u := range unique {
			if v == u {
				unique[i] = v
				continue sampleLoop
			}
		}
		unique = append(unique, v)
	}
	return unique
}

func isBadTicker(coin uint, coinName, coinType, tokenId, currency, provider string, value, change24 float64, tickers []models.Ticker) bool {
	for _, t := range tickers {
		if t.Coin == coin &&
			t.CoinName == coinName &&
			t.CoinType == coinType &&
			t.TokenId == tokenId &&
			t.Currency == currency &&
			t.Provider == provider && (t.Value != value || t.Change24h != change24) {
			return true
		}
	}
	return false
}

func bulkCreateTicker(db *gorm.DB, dataList []models.Ticker) error {
	var (
		valueStrings []string
		valueArgs    []interface{}
	)

	for _, d := range dataList {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, d.Coin)
		valueArgs = append(valueArgs, d.CoinName)
		valueArgs = append(valueArgs, d.CoinType)
		valueArgs = append(valueArgs, d.TokenId)
		valueArgs = append(valueArgs, d.Change24h)
		valueArgs = append(valueArgs, d.Currency)
		valueArgs = append(valueArgs, d.Provider)
		valueArgs = append(valueArgs, d.Value)
		valueArgs = append(valueArgs, d.LastUpdated)
		valueArgs = append(valueArgs, d.Volume)
		valueArgs = append(valueArgs, d.MarketCap)
		valueArgs = append(valueArgs, d.ShowOption)
	}

	smt := fmt.Sprintf(rawBulkTickersInsert, strings.Join(valueStrings, ","))

	if err := db.Exec(smt, valueArgs...).Error; err != nil {
		return err
	}

	return nil
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
