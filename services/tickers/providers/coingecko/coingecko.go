package coingecko

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/services/tickers"
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
	return Provider{ID: id, currency: currency, client: NewClient(api)}
}

func (m Provider) GetData() (tickers.Tickers, error) {
	coins, err := m.client.fetchCoins()
	if err != nil {
		return tickers.Tickers{}, err
	}

	rates := m.client.fetchRates(coins, m.currency, bucketSize)
	tickersList := m.normalizeTickers(rates, coins, m.ID, m.currency)
	return tickersList, nil
}

func (m Provider) normalizeTickers(prices CoinPrices, coins Coins, provider, currency string) tickers.Tickers {
	var (
		tickersList = make(tickers.Tickers, 0)
		cgCoinsMap  = createCgCoinsMap(coins)
	)

	for _, price := range prices {
		t := m.normalizeTicker(price, cgCoinsMap, provider, currency)
		tickersList = append(tickersList, t...)
	}
	return tickersList
}

func (m Provider) normalizeTicker(price CoinPrice, coinsMap map[string][]CoinResult, provider, currency string) tickers.Tickers {
	var (
		tickersList = make(tickers.Tickers, 0)
		tokenId     = ""
		coinName    = strings.ToUpper(price.Symbol)
		coinType    = tickers.Coin
	)

	coins, err := getCgCoinsById(coinsMap, price.Id)
	if err != nil {
		t := createTicker(price, coinType, coinName, tokenId, provider, currency)
		tickersList = append(tickersList, t)
		return tickersList
	}

	for _, cg := range coins {
		coinName = strings.ToUpper(cg.Symbol)
		if cg.CoinType == tickers.Coin {
			tokenId = ""
		} else if len(cg.TokenId) > 0 {
			tokenId = cg.TokenId
		}

		t := createTicker(price, cg.CoinType, coinName, tokenId, provider, currency)
		tickersList = append(tickersList, t)
	}
	return tickersList
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
				CoinType: tickers.Token,
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

func createTicker(price CoinPrice, coinType tickers.CoinType, coinName, tokenId, provider, currency string) tickers.Ticker {
	return tickers.Ticker{
		CoinName: coinName,
		CoinType: coinType,
		TokenId:  tokenId,
		Price: tickers.Price{
			Value:     price.CurrentPrice,
			Change24h: price.PriceChangePercentage24h,
			Currency:  currency,
			Provider:  provider,
		},
		LastUpdate: price.LastUpdated,
	}
}
