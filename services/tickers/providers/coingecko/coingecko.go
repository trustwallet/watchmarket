package coingecko

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	ticker "github.com/trustwallet/watchmarket/services/tickers"
	"strings"
)

const (
	id         = "coingecko"
	bucketSize = 500
)

type Provider struct {
	ID, currency string
	client       Client
}

func InitProvider(api, currency string) Provider {
	return Provider{
		ID:       id,
		currency: currency,
		client:   NewClient(api),
	}
}

func (m *Provider) GetData() (ticker.Tickers, error) {
	var tickers = make(ticker.Tickers, 0)
	coins, err := m.client.FetchCoins()
	if err != nil {
		return tickers, err
	}

	rates := m.client.FetchRates(coins, m.currency, bucketSize)
	tickers = m.normalizeTickers(rates, coins, m.ID, m.currency)
	return tickers, nil
}

func (m *Provider) normalizeTickers(prices CoinPrices, coins Coins, provider, currency string) ticker.Tickers {
	var (
		tickers    = make(ticker.Tickers, 0)
		cgCoinsMap = createCgCoinsMap(coins)
	)

	for _, price := range prices {
		t := m.normalizeTicker(price, cgCoinsMap, provider, currency)
		tickers = append(tickers, t...)
	}
	return tickers
}

func (m *Provider) normalizeTicker(price CoinPrice, coinsMap map[string][]CoinResult, provider, currency string) ticker.Tickers {
	var (
		tickers  = make(ticker.Tickers, 0)
		tokenId  = ""
		coinName = strings.ToUpper(price.Symbol)
		coinType = ticker.Coin
	)

	coins, err := getCgCoinsById(coinsMap, price.Id)
	if err != nil {
		t := createTicker(price, coinType, coinName, tokenId, provider, currency)
		tickers = append(tickers, t)
		return tickers
	}

	for _, cg := range coins {
		coinName = strings.ToUpper(cg.Symbol)
		if cg.CoinType == ticker.Coin {
			tokenId = ""
		} else if len(cg.TokenId) > 0 {
			tokenId = cg.TokenId
		}

		t := createTicker(price, cg.CoinType, coinName, tokenId, provider, currency)
		tickers = append(tickers, t)
	}
	return tickers
}

func getCgCoinsById(coinsMap map[string][]CoinResult, id string) ([]CoinResult, error) {
	coins, ok := coinsMap[id]
	if !ok {
		return nil, errors.E("No coin found by id", errors.Params{"id": id})
	}
	return coins, nil
}

func createCgCoinsMap(coins Coins) map[string][]CoinResult {
	var (
		coinsMap   = getCoinsMap(coins)
		cgCoinsMap = make(map[string][]CoinResult, 0)
	)

	for _, coin := range coins {
		for platform, addr := range coin.Platforms {
			platformCoin, ok := coinsMap[platform]
			if !ok {
				continue
			}

			_, ok = cgCoinsMap[coin.Id]
			if !ok {
				cgCoinsMap[coin.Id] = make([]CoinResult, 0)
			}

			cgCoinsMap[coin.Id] = append(cgCoinsMap[coin.Id], CoinResult{
				Symbol:   platformCoin.Symbol,
				TokenId:  strings.ToLower(addr),
				CoinType: ticker.Token,
			})
		}
	}

	return cgCoinsMap
}

func getCoinsMap(coins Coins) map[string]Coin {
	coinsMap := make(map[string]Coin)
	for _, coin := range coins {
		coinsMap[coin.Id] = coin
	}
	return coinsMap
}

func createTicker(price CoinPrice, coinType ticker.CoinType, coinName, tokenId, provider, currency string) ticker.Ticker {
	var t = ticker.Ticker{
		CoinName: coinName,
		CoinType: coinType,
		TokenId:  tokenId,
		Price: ticker.Price{
			Value:     price.CurrentPrice,
			Change24h: price.PriceChangePercentage24h,
			Currency:  currency,
			Provider:  provider,
		},
		LastUpdate: price.LastUpdated,
	}

	t.Price.Change24h = price.PriceChangePercentage24h
	t.Price.Value = price.CurrentPrice

	return t
}
