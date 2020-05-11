package coingecko

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/services/markets"
	"strings"
)

func (m Provider) GetTickers() (markets.Tickers, error) {
	coins, err := m.client.fetchCoins()
	if err != nil {
		return markets.Tickers{}, err
	}

	rates := m.client.fetchRates(coins)
	tickersList := m.normalizeTickers(rates, coins, m.ID, m.currency)
	return tickersList, nil
}

func (m Provider) normalizeTickers(prices CoinPrices, coins Coins, provider, currency string) markets.Tickers {
	var (
		tickersList = make(markets.Tickers, 0)
		cgCoinsMap  = createCgCoinsMap(coins)
	)

	for _, price := range prices {
		_, ok := cgCoinsMap[price.Id]
		if !ok {
			continue
		}
		t := m.normalizeTicker(price, cgCoinsMap, provider, currency)
		tickersList = append(tickersList, t...)
	}
	return tickersList
}

func (m Provider) normalizeTicker(price CoinPrice, coinsMap map[string][]CoinResult, provider, currency string) markets.Tickers {
	var (
		tickersList = make(markets.Tickers, 0)
		tokenId     = ""
		coinName    = strings.ToUpper(price.Symbol)
		coinType    = markets.Coin
	)

	coins, err := getCgCoinsById(coinsMap, price.Id)
	if err != nil {
		t := createTicker(price, coinType, unknownCoinID, coinName, tokenId, provider, currency)
		tickersList = append(tickersList, t)
		return tickersList
	}

	for _, cg := range coins {
		coinName = strings.ToUpper(cg.Symbol)
		if cg.CoinType == markets.Coin {
			tokenId = ""
		} else if len(cg.TokenId) > 0 {
			tokenId = cg.TokenId
		}

		t := createTicker(price, cg.CoinType, cg.PotentialCoinID, coinName, tokenId, provider, currency)
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

	for _, c := range coins {
		if isBasicCoin(c.Symbol) {
			cr := CoinResult{
				Symbol:          c.Symbol,
				TokenId:         "",
				CoinType:        markets.Coin,
				PotentialCoinID: getCoinBySymbol(c.Symbol).ID,
			}
			cgCoinsMap[c.Id] = []CoinResult{cr}
			continue
		}

		for platform, addr := range c.Platforms {
			if len(platform) == 0 || len(addr) == 0 {
				continue
			}

			platformCoin, ok := coinsMap[platform]
			if !ok {
				continue
			}

			_, ok = cgCoinsMap[c.Id]
			if !ok {
				cgCoinsMap[c.Id] = make([]CoinResult, 0)
			}

			cr := CoinResult{
				Symbol:          platformCoin.Symbol,
				TokenId:         strings.ToLower(addr),
				CoinType:        markets.Token,
				PotentialCoinID: getCoinId(platform),
			}

			cgCoinsMap[c.Id] = []CoinResult{cr}
		}
	}

	return cgCoinsMap
}

func getCoinsMap(coins Coins) map[string]Coin {
	coinsMap := make(map[string]Coin)
	for _, c := range coins {
		coinsMap[c.Id] = c
	}
	return coinsMap
}

func createTicker(price CoinPrice, coinType markets.CoinType, coinID uint, coinName, tokenId, provider, currency string) markets.Ticker {
	return markets.Ticker{
		Coin:     coinID,
		CoinName: coinName,
		CoinType: coinType,
		TokenId:  tokenId,
		Price: markets.Price{
			Value:     price.CurrentPrice,
			Change24h: price.PriceChangePercentage24h,
			Currency:  currency,
			Provider:  provider,
		},
		LastUpdate: price.LastUpdated,
	}
}

func getCoinId(platformName string) uint {
	switch strings.ToLower(platformName) {
	case "binancecoin":
		return coin.Binance().ID
	case "bitcoin-cash":
		return coin.Bitcoincash().ID
	case "ethereum-classic":
		return coin.Classic().ID
	case strings.ToLower(coin.Cosmos().Handle):
		return coin.Cosmos().ID
	case strings.ToLower(coin.Dash().Handle):
		return coin.Dash().ID
	case strings.ToLower(coin.Ethereum().Handle):
		return coin.Ethereum().ID
	case strings.ToLower(coin.Ontology().Handle):
		return coin.Ontology().ID
	case strings.ToLower(coin.Qtum().Handle):
		return coin.Qtum().ID
	case strings.ToLower(coin.Stellar().Handle):
		return coin.Stellar().ID
	case strings.ToLower(coin.Vechain().Handle):
		return coin.Vechain().ID
	case strings.ToLower(coin.Waves().Handle):
		return coin.Waves().ID
	case strings.ToLower(coin.Tron().Handle):
		return coin.Tron().ID
	case strings.ToLower(coin.Classic().Handle):
		return coin.Tron().ID
	case strings.ToLower(coin.Gochain().Handle):
		return coin.Gochain().ID
	case strings.ToLower(coin.Icon().Handle):
		return coin.Icon().ID
	}

	return unknownCoinID
}

func isBasicCoin(symbol string) bool {
	for _, c := range coin.Coins {
		if strings.ToLower(c.Symbol) == strings.ToLower(symbol) {
			return true
		}
	}
	return false
}

func getCoinBySymbol(symbol string) coin.Coin {
	for _, c := range coin.Coins {
		if strings.ToLower(c.Symbol) == strings.ToLower(symbol) {
			return c
		}
	}
	return coin.Coin{}
}
