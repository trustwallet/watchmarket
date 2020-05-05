package binancedex

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	ticker "github.com/trustwallet/watchmarket/services/tickers"
	"strconv"
	"time"
)

var (
	id       = "binancedex"
	BNBAsset = coin.Binance().Symbol
)

type Parser struct {
	ID, currency string
	client       Client
}

func InitParser(api, currency string) Parser {
	m := Parser{
		ID:       id,
		currency: currency,
		client:   NewClient(api),
	}
	return m
}

func (p Parser) GetData() (ticker.Tickers, error) {
	prices, err := p.client.GetPrices()
	if err != nil {
		return nil, err
	}
	return normalizeTickers(prices, p.ID, p.currency), nil
}

func normalizeTickers(prices []CoinPrice, provider, currency string) ticker.Tickers {
	tickers := make(ticker.Tickers, 0)
	for _, price := range prices {
		t, err := normalizeTicker(price, provider, currency)
		if err != nil {
			continue
		}
		tickers = append(tickers, t)
	}
	return tickers
}

func normalizeTicker(price CoinPrice, provider, currency string) (ticker.Ticker, error) {
	var t ticker.Ticker

	if price.QuoteAssetName != BNBAsset && price.BaseAssetName != BNBAsset {
		return t, errors.E("invalid quote/base asset",
			errors.Params{"Symbol": price.BaseAssetName, "QuoteAsset": price.QuoteAssetName})
	}

	value, err := strconv.ParseFloat(price.LastPrice, 64)
	if err != nil {
		return t, errors.E(err, "normalizeTicker parse value error",
			errors.Params{"LastPrice": price.LastPrice, "Symbol": price.BaseAssetName})
	}

	value24h, err := strconv.ParseFloat(price.PriceChangePercent, 64)
	if err != nil {
		return t, errors.E(err, "normalizeTicker parse value24h error",
			errors.Params{"PriceChange": price.PriceChangePercent, "Symbol": price.BaseAssetName})
	}

	tokenId := price.BaseAssetName
	if tokenId == BNBAsset {
		tokenId = price.QuoteAssetName
		value = 1.0 / value
	}

	t = ticker.Ticker{
		CoinName: BNBAsset,
		CoinType: ticker.Token,
		TokenId:  tokenId,
		Price: ticker.Price{
			Value:     value,
			Change24h: value24h,
			Currency:  currency,
			Provider:  provider,
		},
		LastUpdate: time.Now(),
	}
	return t, nil
}
