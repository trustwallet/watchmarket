package binancedex

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/services/markets"
	"strconv"
	"time"
)

var (
	id       = "binancedex"
	BNBAsset = coin.Binance().Symbol
)

func (p Provider) GetTickers() (markets.Tickers, error) {
	prices, err := p.client.fetchPrices()
	if err != nil {
		return nil, err
	}
	return normalizeTickers(prices, p.ID), nil
}

func normalizeTickers(prices []CoinPrice, provider string) markets.Tickers {
	tickersList := make(markets.Tickers, 0)
	for _, price := range prices {
		t, err := normalizeTicker(price, provider)
		if err != nil {
			logger.Error(err)
			continue
		}
		tickersList = append(tickersList, t)
	}
	return tickersList
}

func normalizeTicker(price CoinPrice, provider string) (markets.Ticker, error) {
	var t markets.Ticker

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

	t = markets.Ticker{
		Coin:     coin.BNB,
		CoinName: BNBAsset,
		CoinType: markets.Token,
		TokenId:  tokenId,
		Price: markets.Price{
			Value:     value,
			Change24h: value24h,
			Currency:  BNBAsset,
			Provider:  provider,
		},
		LastUpdate: time.Now(),
	}
	return t, nil
}
