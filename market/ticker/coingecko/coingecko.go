package coingecko

import (
	"github.com/trustwallet/watchmarket/market/clients/coingecko"
	"github.com/trustwallet/watchmarket/market/ticker"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
)

const (
	id                   = "coingecko"
	minimalTradingVolume = 5000
)

type Market struct {
	client *coingecko.Client
	cache  *coingecko.Cache
	ticker.Market
}

func InitMarket(api, updateTime string) ticker.TickerProvider {
	m := &Market{
		client: coingecko.NewClient(api),
		Market: ticker.Market{
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
		t := createTicker(price, coinType, coinName, tokenId, provider)
		tickers = append(tickers, &t)
		return tickers
	}

	for _, cg := range coins {
		coinName = strings.ToUpper(cg.Symbol)
		if cg.CoinType == watchmarket.TypeCoin {
			tokenId = ""
		} else if len(cg.TokenId) > 0 {
			tokenId = cg.TokenId
		}

		t := createTicker(price, cg.CoinType, coinName, tokenId, provider)
		tickers = append(tickers, &t)
	}
	return tickers
}

func createTicker(price coingecko.CoinPrice, coinType watchmarket.CoinType, coinName, tokenId, provider string) watchmarket.Ticker {
	var t = watchmarket.Ticker{
		CoinName: coinName,
		CoinType: coinType,
		TokenId:  tokenId,
		Price: watchmarket.TickerPrice{
			Value:     price.CurrentPrice,
			Change24h: price.PriceChangePercentage24h,
			Currency:  watchmarket.DefaultCurrency,
			Provider:  provider,
		},
		LastUpdate: price.LastUpdated,
	}

	if isRespectableTradingVolume(price.TotalVolume) {
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

func (m *Market) normalizeTickers(prices coingecko.CoinPrices, provider string) watchmarket.Tickers {
	var tickers = make(watchmarket.Tickers, 0)
	for _, price := range prices {
		t := m.normalizeTicker(price, provider)
		tickers = append(tickers, t...)
	}
	return tickers
}
