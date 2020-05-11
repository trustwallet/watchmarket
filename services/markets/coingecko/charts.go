package coingecko

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/services/markets"
	"sort"
	"strings"
	"time"
)

func (p Provider) GetChartData(coinID uint, token, currency string, timeStart int64) (markets.Data, error) {
	chartsData := markets.Data{}

	coins, err := p.client.fetchCoins()
	if err != nil {
		return chartsData, err
	}

	symbolsMap := createSymbolsMap(coins)

	coinResult, err := getCoinByID(symbolsMap, coinID, token)
	if err != nil {
		return chartsData, err
	}

	timeEndDate := time.Now().Unix()

	c, err := p.client.fetchCharts(coinResult.Id, currency, timeStart, timeEndDate)
	if err != nil {
		return chartsData, err
	}

	return normalizeCharts(c), nil
}

func (p Provider) GetCoinData(coinID uint, token, currency string) (markets.CoinDetails, error) {
	coins, err := p.client.fetchCoins()
	if err != nil {
		return markets.CoinDetails{}, err
	}

	symbolsMap := createSymbolsMap(coins)

	coinResult, err := getCoinByID(symbolsMap, coinID, token)
	if err != nil {
		return markets.CoinDetails{}, err
	}

	ratesData := p.client.fetchRates(coins)
	if len(ratesData) == 0 {
		return markets.CoinDetails{}, errors.E("No rates found", errors.Params{"id": coinResult.Id})
	}

	infoData, err := p.info.GetCoinInfo(coinID, token)
	if err != nil {
		logger.Warn("No assets info about that coin", logger.Params{"coin": coinID, "token": token})
	}
	return normalizeInfo(ratesData[0], infoData), nil
}

func createSymbolsMap(coins Coins) map[string]Coin {
	var (
		symbolsMap = make(map[string]Coin, 0)
		coinsMap   = createCoinsMap(coins)
	)

	for _, c := range coins {
		if len(c.Platforms) == 0 {
			symbolsMap[createID(c.Symbol, "")] = c
		}
		for platform, addr := range c.Platforms {
			if len(platform) == 0 || len(addr) == 0 {
				continue
			}
			platformCoin, ok := coinsMap[platform]
			if !ok {
				continue
			}
			symbolsMap[createID(platformCoin.Symbol, addr)] = c
		}
	}

	return symbolsMap
}

func createCoinsMap(coins Coins) map[string]Coin {
	coinsMap := make(map[string]Coin)
	for _, c := range coins {
		coinsMap[c.Id] = c
	}
	return coinsMap
}

func createID(symbol, token string) string {
	if token != "" {
		return strings.ToLower(symbol + token)
	}
	return strings.ToLower(symbol)
}

func getCoinByID(coinMap map[string]Coin, coinId uint, token string) (Coin, error) {
	c := Coin{}
	coinObj, ok := coin.Coins[coinId]
	if !ok {
		return c, errors.E("Coin not found", errors.Params{"coindId": coinId})
	}

	c, err := getCoinByParams(coinMap, coinObj.Symbol, token)
	if err != nil {
		return c, err
	}

	return c, nil
}

func getCoinByParams(coinMap map[string]Coin, symbol, token string) (Coin, error) {
	c, ok := coinMap[createID(symbol, token)]
	if !ok {
		return c, errors.E("No coin found by symbol", errors.Params{"symbol": symbol, "token": token})
	}
	return c, nil
}

func normalizeCharts(c Charts) markets.Data {
	chartsData := markets.Data{}
	prices := make([]markets.ChartsPrice, 0)
	for _, quote := range c.Prices {
		if len(quote) != chartDataSize {
			continue
		}

		date := time.Unix(int64(quote[0])/1000, 0)
		prices = append(prices, markets.ChartsPrice{
			Price: quote[1],
			Date:  date.Unix(),
		})
	}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Date < prices[j].Date
	})

	chartsData.Prices = prices

	return chartsData
}

func normalizeInfo(data CoinPrice, info markets.Info) markets.CoinDetails {
	return markets.CoinDetails{
		Vol24:             data.TotalVolume,
		MarketCap:         data.MarketCap,
		CirculatingSupply: data.CirculatingSupply,
		TotalSupply:       data.TotalSupply,
		Info:              info,
	}
}
