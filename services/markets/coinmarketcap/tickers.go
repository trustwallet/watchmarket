package coinmarketcap

import (
	"context"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
)

func (p Provider) GetTickers(ctx context.Context) (watchmarket.Tickers, error) {
	prices, err := p.client.fetchPrices(p.currency, ctx)
	if err != nil {
		return nil, err
	}

	return normalizeTickers(prices, p.Cm, p.id, p.currency), nil
}

func normalizeTickers(prices CoinPrices, coinsMap []CoinMap, provider, currency string) watchmarket.Tickers {
	var tickersList watchmarket.Tickers
	for _, price := range prices.Data {
		t := normalizeTicker(price, coinsMap, provider, currency)
		tickersList = append(tickersList, t...)
	}
	return tickersList
}

func normalizeTicker(price Data, coinsMap []CoinMap, provider, currency string) watchmarket.Tickers {
	var (
		tokenId       string
		tickersList   watchmarket.Tickers
		coinName      = price.Symbol
		coinType      = watchmarket.Coin
		emptyPlatform = Platform{}
	)

	if price.Platform != emptyPlatform {
		tokenId = strings.ToLower(price.Platform.TokenAddress)
		coinType = watchmarket.Token
		coinName = price.Platform.Symbol
		if len(tokenId) == 0 {
			tokenId = price.Symbol
		}
	}

	mappedCmcCoins, err := findCoin(coinsMap, price.Id)
	if err != nil {
		tickersList = append(tickersList, watchmarket.Ticker{
			Coin:     watchmarket.UnknownCoinID,
			CoinName: coinName,
			CoinType: coinType,
			TokenId:  tokenId,
			Price: watchmarket.Price{
				Value:     price.Quote.USD.Price,
				Change24h: price.Quote.USD.PercentChange24h,
				Currency:  currency,
				Provider:  provider,
			},
			LastUpdate: price.LastUpdated,
			Volume:     price.Quote.USD.Volume,
			MarketCap:  price.Quote.USD.MarketCap,
		})
		return tickersList
	}
	for _, mappedCmcCoin := range mappedCmcCoins {
		coinName = mappedCmcCoin.Coin.Symbol
		if mappedCmcCoin.CoinType == watchmarket.Coin {
			tokenId = ""
		} else if len(mappedCmcCoin.TokenId) > 0 {
			tokenId = strings.ToLower(mappedCmcCoin.TokenId)
		}
		tickersList = append(tickersList, watchmarket.Ticker{
			Coin:     mappedCmcCoin.Coin.ID,
			CoinName: coinName,
			CoinType: mappedCmcCoin.CoinType,
			TokenId:  strings.ToLower(tokenId),
			Price: watchmarket.Price{
				Value:     price.Quote.USD.Price,
				Change24h: price.Quote.USD.PercentChange24h,
				Currency:  currency,
				Provider:  provider,
			},
			LastUpdate: price.LastUpdated,
			Volume:     price.Quote.USD.Volume,
			MarketCap:  price.Quote.USD.MarketCap,
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
		result = append(result, CoinResult{Coin: atlasCoin, Id: mappedCoin.Id, TokenId: mappedCoin.TokenId, CoinType: watchmarket.CoinType(mappedCoin.Type)})
	}
	return result, nil
}
