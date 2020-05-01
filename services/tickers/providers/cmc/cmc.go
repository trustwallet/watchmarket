package cmc

import (
	ticker "github.com/trustwallet/watchmarket/services/tickers"
	"strings"
)

const (
	id = "cmc"
)

type Parser struct {
	ID, currency string
	client       Client
}

func InitParser(api, key, currency string) Parser {
	m := &Parser{
		ID:       id,
		currency: currency,
		client:   NewClient(api, key),
	}
	return *m
}

func (m *Parser) GetData() (ticker.Tickers, error) {
	prices, err := m.client.GetData(m.currency)
	if err != nil {
		return nil, err
	}
	return normalizeTickers(prices, m.ID, m.currency), nil
}

func normalizeTicker(price Data, provider, currency string) ticker.Tickers {
	var (
		tokenId string
		tickers ticker.Tickers

		coinName = price.Symbol
		coinType = ticker.Coin
	)

	if price.Platform != nil {
		tokenId = strings.ToLower(price.Platform.TokenAddress)
		coinType = ticker.Token
		coinName = price.Platform.Symbol
		if len(tokenId) == 0 {
			tokenId = price.Symbol
		}
	}

	tickers = append(tickers, &ticker.Ticker{
		CoinName: coinName,
		CoinType: coinType,
		TokenId:  tokenId,
		Price: ticker.Price{
			Value:     price.Quote.USD.Price,
			Change24h: price.Quote.USD.PercentChange24h,
			Currency:  currency,
			Provider:  provider,
		},
		LastUpdate: price.LastUpdated,
	})

	return tickers

}

func normalizeTickers(prices CoinPrices, provider, currency string) ticker.Tickers {
	var tickers ticker.Tickers
	for _, price := range prices.Data {
		t := normalizeTicker(price, provider, currency)
		tickers = append(tickers, t...)
	}
	return tickers
}
