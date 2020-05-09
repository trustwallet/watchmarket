package coinmarketcap

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/services/tickers"
	"strings"
)

const (
	id = "coinmarketcap"
)

type Provider struct {
	ID, currency string
	client       Client
}

func InitProvider(proApi, assetsApi, key, currency string) Provider {
	return Provider{ID: id, currency: currency, client: NewClient(proApi, assetsApi, key)}
}

func (p Provider) GetData() (tickers.Tickers, error) {
	prices, err := p.client.FetchPrices(p.currency)
	if err != nil {
		return nil, err
	}

	coinsMap, err := p.client.FetchCoinMap()
	if err != nil {
		return nil, err
	}

	return normalizeTickers(prices, coinsMap, p.ID, p.currency), nil
}

func normalizeTickers(prices CoinPrices, coinsMap []CoinMap, provider, currency string) tickers.Tickers {
	var tickersList tickers.Tickers
	for _, price := range prices.Data {
		t := normalizeTicker(price, coinsMap, provider, currency)
		tickersList = append(tickersList, t...)
	}
	return tickersList
}

func normalizeTicker(price Data, coinsMap []CoinMap, provider, currency string) tickers.Tickers {
	var (
		tokenId       string
		tickersList   tickers.Tickers
		coinName      = price.Symbol
		coinType      = tickers.Coin
		emptyPlatform = Platform{}
	)

	if price.Platform != emptyPlatform {
		tokenId = strings.ToLower(price.Platform.TokenAddress)
		coinType = tickers.Token
		coinName = price.Platform.Symbol
		if len(tokenId) == 0 {
			tokenId = price.Symbol
		}
	}

	mappedCmcCoins, err := findCoin(coinsMap, price.Id)
	if err != nil {
		tickersList = append(tickersList, tickers.Ticker{
			CoinName: coinName,
			CoinType: coinType,
			TokenId:  tokenId,
			Price: tickers.Price{
				Value:     price.Quote.USD.Price,
				Change24h: price.Quote.USD.PercentChange24h,
				Currency:  currency,
				Provider:  provider,
			},
			LastUpdate: price.LastUpdated,
		})
		return tickersList
	}
	for _, mappedCmcCoin := range mappedCmcCoins {
		coinName = mappedCmcCoin.Coin.Symbol
		if mappedCmcCoin.CoinType == tickers.Coin {
			tokenId = ""
		} else if len(mappedCmcCoin.TokenId) > 0 {
			tokenId = strings.ToLower(mappedCmcCoin.TokenId)
		}
		tickersList = append(tickersList, tickers.Ticker{
			Coin:     mappedCmcCoin.Coin.ID,
			CoinName: coinName,
			CoinType: mappedCmcCoin.CoinType,
			TokenId:  tokenId,
			Price: tickers.Price{
				Value:     price.Quote.USD.Price,
				Change24h: price.Quote.USD.PercentChange24h,
				Currency:  currency,
				Provider:  provider,
			},
			LastUpdate: price.LastUpdated,
		})
	}

	return tickersList
}

func findCoin(rawCoins []CoinMap, coinId uint) ([]CoinResult, error) {
	coinMap := make(map[uint][]CoinMap)
	for _, rawCoin := range rawCoins {
		_, ok := coinMap[rawCoin.Id]
		if !ok {
			coinMap[rawCoin.Id] = make([]CoinMap, 0)
		}
		coinMap[rawCoin.Id] = append(coinMap[rawCoin.Id], rawCoin)
	}

	mappedCoins, ok := coinMap[coinId]
	if !ok {
		return nil, errors.E("findCoin coinId notFound", errors.Params{"coin_ID": coinId})
	}

	result := make([]CoinResult, 0)
	for _, mappedCoin := range mappedCoins {
		atlasCoin, ok := coin.Coins[mappedCoin.Coin]
		if !ok {
			continue
		}
		result = append(result, CoinResult{Coin: atlasCoin, Id: mappedCoin.Id, TokenId: mappedCoin.TokenId, CoinType: tickers.CoinType(mappedCoin.Type)})
	}
	return result, nil
}
