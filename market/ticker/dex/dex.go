package dex

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/market/ticker"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"math/big"
	"net/url"
	"strconv"
	"time"
)

var (
	id       = "dex"
	BNBAsset = coin.Binance().Symbol
)

type Market struct {
	ticker.Market
	blockatlas.Request
}

func InitMarket(api string, updateTime string) ticker.TickerProvider {
	m := &Market{
		Market: ticker.Market{
			Id:         id,
			UpdateTime: updateTime,
		},
		Request: blockatlas.InitClient(api),
	}
	return m
}

func (m *Market) GetData() (watchmarket.Tickers, error) {
	var prices []*CoinPrice
	err := m.Get(&prices, "v1/ticker/24hr", url.Values{"limit": {"1000"}})
	if err != nil {
		return nil, err
	}
	rate, err := m.Storage.GetRate(BNBAsset)
	if err != nil {
		return nil, errors.E(err, "rate not found", errors.Params{"asset": BNBAsset})
	}
	result := normalizeTickers(prices, m.GetId())
	if rate.PercentChange24h != nil {
		rate.PercentChange24h.Mul(rate.PercentChange24h, big.NewFloat(-1))
	}
	result.ApplyRate(watchmarket.DefaultCurrency, 1/rate.Rate, rate.PercentChange24h)
	return result, nil
}

func normalizeTicker(price *CoinPrice, provider string) (*watchmarket.Ticker, error) {
	if price.QuoteAssetName != BNBAsset && price.BaseAssetName != BNBAsset {
		return nil, errors.E("invalid quote/base asset",
			errors.Params{"Symbol": price.BaseAssetName, "QuoteAsset": price.QuoteAssetName})
	}
	value, err := strconv.ParseFloat(price.LastPrice, 64)
	if err != nil {
		return nil, errors.E(err, "normalizeTicker parse value error",
			errors.Params{"LastPrice": price.LastPrice, "Symbol": price.BaseAssetName})
	}
	value24h, err := strconv.ParseFloat(price.PriceChangePercent, 64)
	if err != nil {
		return nil, errors.E(err, "normalizeTicker parse value24h error",
			errors.Params{"PriceChange": price.PriceChangePercent, "Symbol": price.BaseAssetName})
	}
	tokenId := price.BaseAssetName
	if tokenId == BNBAsset {
		tokenId = price.QuoteAssetName
		value = 1.0 / value
	}
	return &watchmarket.Ticker{
		CoinName: BNBAsset,
		CoinType: watchmarket.TypeToken,
		TokenId:  tokenId,
		Price: watchmarket.TickerPrice{
			Value:     value,
			Change24h: value24h,
			Currency:  "BNB",
			Provider:  provider,
		},
		LastUpdate: time.Now(),
	}, nil
}

func normalizeTickers(prices []*CoinPrice, provider string) (tickers watchmarket.Tickers) {
	for _, price := range prices {
		t, err := normalizeTicker(price, provider)
		if err != nil {
			continue
		}
		tickers = append(tickers, t)
	}
	return
}
