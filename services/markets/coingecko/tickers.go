package coingecko

import (
	"errors"
	"sort"
	"strings"

	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (p Provider) GetTickers() (watchmarket.Tickers, error) {
	coins, err := p.client.fetchCoins()
	if err != nil {
		return watchmarket.Tickers{}, err
	}

	rates := p.client.fetchRates(coins, p.currency)
	tickersList := p.normalizeTickers(rates, coins, p.id, p.currency)
	return tickersList, nil
}

func (p Provider) normalizeTickers(prices CoinPrices, coins Coins, provider, currency string) watchmarket.Tickers {
	var (
		tickersList = make(watchmarket.Tickers, 0)
		cgCoinsMap  = createCgCoinsMap(coins)
	)

	for _, price := range prices {
		_, ok := cgCoinsMap[price.Id]
		if !ok {
			continue
		}
		t := p.normalizeTicker(price, cgCoinsMap, provider, currency)
		tickersList = append(tickersList, t...)
	}
	return tickersList
}

func (p Provider) normalizeTicker(price CoinPrice, coinsMap map[string][]CoinResult, provider, currency string) watchmarket.Tickers {
	var (
		tickersList = make(watchmarket.Tickers, 0)
		tokenId     = ""
		coinName    = strings.ToUpper(price.Symbol)
		coinType    = watchmarket.Coin
	)

	coins, err := getCgCoinsById(coinsMap, price.Id)
	if err != nil {
		t := createTicker(price, coinType, watchmarket.UnknownCoinID, coinName, tokenId, provider, currency)
		tickersList = append(tickersList, t)
		return tickersList
	}

	for _, cg := range coins {
		coinName = strings.ToUpper(cg.Symbol)
		if cg.CoinType == watchmarket.Coin {
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
		return nil, errors.New("no coin found by id")
	}
	return coins, nil
}

func createCgCoinsMap(coins Coins) map[string][]CoinResult {
	var (
		coinsMap   = getCoinsMap(coins)
		cgCoinsMap = make(map[string][]CoinResult)
	)

	for _, c := range coins {
		if isBasicCoin(c.Symbol) {
			cr := CoinResult{
				Symbol:          c.Symbol,
				TokenId:         "",
				CoinType:        watchmarket.Coin,
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
				CoinType:        watchmarket.Token,
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

func createTicker(price CoinPrice, coinType watchmarket.CoinType, coinID uint, coinName, tokenId, provider, currency string) watchmarket.Ticker {
	return watchmarket.Ticker{
		Coin:     coinID,
		CoinName: coinName,
		CoinType: coinType,
		TokenId:  strings.ToLower(tokenId),
		Price: watchmarket.Price{
			Value:     price.CurrentPrice,
			Change24h: price.PriceChangePercentage24h,
			Currency:  currency,
			Provider:  provider,
		},
		LastUpdate:        price.LastUpdated,
		MarketCap:         price.MarketCap,
		Volume:            price.TotalVolume,
		CirculatingSupply: price.CirculatingSupply,
		TotalSupply:       price.TotalSupply,
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
	case strings.ToLower(coin.Polkadot().Handle):
		return coin.Polkadot().ID
	case strings.ToLower(coin.Elrond().Handle):
		return coin.Elrond().ID
	case strings.ToLower(coin.Filecoin().Handle):
		return coin.Filecoin().ID
	}

	return watchmarket.UnknownCoinID
}

func isBasicCoin(symbol string) bool {
	for _, c := range coin.Coins {
		if strings.EqualFold(c.Symbol, symbol) {
			return true
		}
	}
	return false
}

func getCoinBySymbol(symbol string) coin.Coin {
	ids := []int{}
	for _, c := range coin.Coins {
		ids = append(ids, int(c.ID))
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] > ids[j]
	})
	for _, id := range ids {
		c := coin.Coins[uint(id)]
		if strings.EqualFold(c.Symbol, symbol) {
			return c
		}
	}
	return coin.Coin{}
}
