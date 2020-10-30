package binancedex

import (
	"context"
	"strconv"
	"strings"
	"time"

	"errors"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

var (
	id       = "binancedex"
	BNBAsset = coin.Binance().Symbol
)

func (p Provider) GetTickers(ctx context.Context) (watchmarket.Tickers, error) {
	prices, err := p.client.fetchPrices(ctx)
	if err != nil {
		return nil, err
	}
	return normalizeTickers(prices, p.id), nil
}

func normalizeTickers(prices []CoinPrice, provider string) watchmarket.Tickers {
	tickersList := make(watchmarket.Tickers, 0)
	for _, price := range prices {
		t, err := normalizeTicker(price, provider)
		if err != nil {
			continue
		}
		tickersList = append(tickersList, t)
	}
	return tickersList
}

func normalizeTicker(price CoinPrice, provider string) (watchmarket.Ticker, error) {
	var t watchmarket.Ticker

	if price.QuoteAssetName != BNBAsset && price.BaseAssetName != BNBAsset {
		return t, errors.New("invalid quote/base asset")
	}

	value, err := strconv.ParseFloat(price.LastPrice, 64)
	if err != nil {
		return t, errors.New(err.Error() + " normalizeTicker parse value error")
	}

	value24h, err := strconv.ParseFloat(price.PriceChangePercent, 64)
	if err != nil {
		return t, errors.New(err.Error() + " normalizeTicker parse value24h error")
	}

	tokenId := price.BaseAssetName
	if tokenId == BNBAsset {
		tokenId = price.QuoteAssetName
		value = 1.0 / value
	}

	volume, err := strconv.ParseFloat(price.Volume, 32)
	if err != nil {
		volume = 0
	}

	t = watchmarket.Ticker{
		Coin:     coin.BNB,
		CoinName: BNBAsset,
		CoinType: watchmarket.Token,
		TokenId:  strings.ToLower(tokenId),
		Price: watchmarket.Price{
			Value:     value,
			Change24h: value24h,
			Currency:  BNBAsset,
			Provider:  provider,
		},
		LastUpdate: time.Now(),
		Volume:     volume,
	}
	return t, nil
}
