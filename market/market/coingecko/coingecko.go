package coingecko

import (
	"github.com/trustwallet/watchmarket/market/clients/coingecko"
	"github.com/trustwallet/watchmarket/market/market"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
)

const (
	id               = "coingecko"
	minimalMarketCap = 1
)

type Market struct {
	client *coingecko.Client
	cache  *coingecko.Cache
	market.Market
}

func InitMarket(api, updateTime string) market.Provider {
	m := &Market{
		client: coingecko.NewClient(api),
		Market: market.Market{
			Id:         id,
			UpdateTime: updateTime,
		},
	}
	return m
}

func (m *Market) GetData() (watchmarket.Tickers, error) {
	var tickers = make(watchmarket.Tickers, 0)
	coins, err := m.client.FetchCoinsList()
	if err != nil {
		return tickers, err
	}
	m.cache = coingecko.NewCache(coins)

	rates := m.client.FetchLatestRates(coins, watchmarket.DefaultCurrency)
	tickers = m.normalizeTickers(rates, m.GetId())
	return tickers, nil
}

func (m *Market) normalizeTicker(price coingecko.CoinPrice, provider string) watchmarket.Tickers {
	var tickers = make(watchmarket.Tickers, 0)
	tokenId := ""
	coinName := strings.ToUpper(price.Symbol)
	coinType := watchmarket.TypeCoin

	coins, err := m.cache.GetCoinsById(price.Id)
	if err != nil {
		ticker := createTicker(price, coinType, coinName, tokenId, provider)
		tickers = append(tickers, &ticker)
		return tickers
	}

	for _, cg := range coins {
		coinName = strings.ToUpper(cg.Symbol)
		if cg.CoinType == watchmarket.TypeCoin {
			tokenId = ""
		} else if len(cg.TokenId) > 0 {
			tokenId = cg.TokenId
		}

		ticker := createTicker(price, cg.CoinType, coinName, tokenId, provider)
		tickers = append(tickers, &ticker)
	}
	return tickers
}

func createTicker(price coingecko.CoinPrice, coinType watchmarket.CoinType, coinName, tokenId, provider string) watchmarket.Ticker {
	var ticker = watchmarket.Ticker{
		CoinName: coinName,
		CoinType: coinType,
		TokenId:  tokenId,
		Price: watchmarket.TickerPrice{
			Value:     0,
			Change24h: 0,
			Currency:  watchmarket.DefaultCurrency,
			Provider:  provider,
		},
		LastUpdate: price.LastUpdated,
	}

	if isRespectableMarketCap(price.MarketCap) {
		ticker.Price.Change24h = price.PriceChangePercentage24h
		ticker.Price.Value = price.CurrentPrice
	}
	return ticker
}

func isRespectableMarketCap(targetMarketCap float64) bool {
	return targetMarketCap >= minimalMarketCap
}

func (m *Market) normalizeTickers(prices coingecko.CoinPrices, provider string) watchmarket.Tickers {
	var tickers = make(watchmarket.Tickers, 0)
	for _, price := range prices {
		t := m.normalizeTicker(price, provider)
		tickers = append(tickers, t...)
	}
	return tickers
}
