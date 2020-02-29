package coingecko

import (
	"github.com/trustwallet/watchmarket/market/clients/coingecko"
	"github.com/trustwallet/watchmarket/market/market"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
)

const (
	id               = "coingecko"
	minimalMarketCap = 0
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

func (m *Market) GetData() (result watchmarket.Tickers, err error) {
	coins, err := m.client.FetchCoinsList()
	if err != nil {
		return
	}
	m.cache = coingecko.NewCache(coins)

	rates := m.client.FetchLatestRates(coins, watchmarket.DefaultCurrency)
	result = m.normalizeTickers(rates, m.GetId())
	return
}

func (m *Market) normalizeTicker(price coingecko.CoinPrice, provider string) (tickers watchmarket.Tickers) {
	tokenId := ""
	coinName := strings.ToUpper(price.Symbol)
	coinType := watchmarket.TypeCoin

	cgCoins, err := m.cache.GetCoinsById(price.Id)
	if err != nil {
		ticker := createTicker(price, coinType, coinName, tokenId, provider)
		tickers = append(tickers, &ticker)
		return
	}

	for _, cg := range cgCoins {
		coinName = strings.ToUpper(cg.Symbol)
		if cg.CoinType == watchmarket.TypeCoin {
			tokenId = ""
		} else if len(cg.TokenId) > 0 {
			tokenId = cg.TokenId
		}

		ticker := createTicker(price, cg.CoinType, coinName, tokenId, provider)
		tickers = append(tickers, &ticker)
	}
	return
}

func createTicker(price coingecko.CoinPrice, coinType watchmarket.CoinType, coinName, tokenId, provider string) watchmarket.Ticker {
	if isRespectableMarketCap(price.MarketCap) {
		return watchmarket.Ticker{
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
	}
	return watchmarket.Ticker{
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
}

func isRespectableMarketCap(targetMarketCap float64) bool {
	return targetMarketCap >= minimalMarketCap
}

func (m *Market) normalizeTickers(prices coingecko.CoinPrices, provider string) (tickers watchmarket.Tickers) {
	for _, price := range prices {
		t := m.normalizeTicker(price, provider)
		tickers = append(tickers, t...)
	}
	return
}
