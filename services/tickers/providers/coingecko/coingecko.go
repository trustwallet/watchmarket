package coingecko

import (
	"github.com/trustwallet/watchmarket/services/clients/coingecko"
	ticker "github.com/trustwallet/watchmarket/services/tickers"
	"strings"
)

const (
	id                   = "coingecko"
	minimalTradingVolume = 5000
	minimalMarketCap     = 5000
	bucketSize           = 500
)

type Parser struct {
	ID, currency string
	client       Client
}

func InitMarket(api, currency string) Parser {
	p := &Parser{
		ID:       id,
		currency: currency,
		client:   *NewClient(api),
	}
	return *p
}

func (m *Parser) GetData() (ticker.Tickers, error) {
	var tickers = make(ticker.Tickers, 0)
	coins, err := m.client.FetchCoinsList()
	if err != nil {
		return tickers, err
	}

	rates := m.client.FetchLatestRates(coins, m.currency, bucketSize)
	tickers = m.normalizeTickers(rates, m.ID)
	return tickers, nil
}

func (m *Parser) normalizeTicker(price CoinPrice, provider string) ticker.Tickers {
	var tickers = make(ticker.Tickers, 0)
	tokenId := ""
	coinName := strings.ToUpper(price.Symbol)
	coinType := ticker.Coin

	coins, err := m.cache.GetCoinsById(price.Id)
	if err != nil {
		t := createTicker(price, coinType, coinName, tokenId, provider)
		tickers = append(tickers, &t)
		return tickers
	}

	for _, cg := range coins {
		coinName = strings.ToUpper(cg.Symbol)
		if cg.CoinType == ticker.Coin {
			tokenId = ""
		} else if len(cg.TokenId) > 0 {
			tokenId = cg.TokenId
		}

		t := createTicker(price, cg.CoinType, coinName, tokenId, provider)
		tickers = append(tickers, &t)
	}
	return tickers
}

func createTicker(price CoinPrice, coinType ticker.CoinType, coinName, tokenId, provider string) ticker.Ticker {
	var t = ticker.Ticker{
		CoinName: coinName,
		CoinType: coinType,
		TokenId:  tokenId,
		Price: ticker.Price{
			Value:     price.CurrentPrice,
			Change24h: price.PriceChangePercentage24h,
			Currency:  watchmarket.DefaultCurrency,
			Provider:  provider,
		},
		LastUpdate: price.LastUpdated,
	}

	if isRespectableTradingVolume(price.TotalVolume) && isRespectableMarketCap(price.MarketCap) {
		t.Price.Change24h = price.PriceChangePercentage24h
		t.Price.Value = price.CurrentPrice
	} else {
		t.Price.Change24h = 0
		t.Price.Value = 0
	}

	return t
}

func isRespectableTradingVolume(targetTradingVolume float64) bool {
	return targetTradingVolume >= minimalTradingVolume
}

func isRespectableMarketCap(targetMarketCap float64) bool {
	return targetMarketCap >= minimalMarketCap
}

func (m *Parser) normalizeTickers(prices CoinPrices, provider string) ticker.Tickers {
	var tickers = make(ticker.Tickers, 0)
	for _, price := range prices {
		t := m.normalizeTicker(price, provider)
		tickers = append(tickers, t...)
	}
	return tickers
}
