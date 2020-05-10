package coingecko

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/services/charts"
	"github.com/trustwallet/watchmarket/services/charts/info"
	"sort"
	"strings"
	"time"
)

const (
	id            = "coingecko"
	chartDataSize = 2
)

type Provider struct {
	ID     string
	client Client
	info   info.Client
}

func InitProvider(webApi, infoApi string) Provider {
	return Provider{ID: id, client: NewClient(webApi), info: info.NewClient(infoApi)}
}

func (p Provider) GetChartData(coinID uint, token, currency string, timeStart int64) (charts.Data, error) {
	chartsData := charts.Data{}

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

func (p Provider) GetCoinData(coinId uint, token, currency string) (charts.CoinDetails, error) {
	coins, err := p.client.fetchCoins()
	if err != nil {
		return charts.CoinDetails{}, err
	}

	symbolsMap := createSymbolsMap(coins)

	coinResult, err := getCoinByID(symbolsMap, coinId, token)
	if err != nil {
		return charts.CoinDetails{}, err
	}

	data := p.client.fetchRates(coins, currency, 500)
	if len(data) == 0 {
		return charts.CoinDetails{}, errors.E("No rates found", errors.Params{"id": coinResult.Id})
	}

	return normalizeInfo(data[0]), nil
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
			if strings.EqualFold(platformCoin.Symbol, addr) {
				symbolsMap[createID(platformCoin.Symbol, "")] = c
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
	if len(token) > 0 {
		strings.ToLower(symbol + token)
	}
	return strings.ToLower(symbol)
}

func getCoinByID(coinMap map[string]Coin, coinId uint, token string) (Coin, error) {
	c := Coin{}
	coinObj, ok := coin.Coins[coinId]
	if !ok {
		return c, errors.E("Coin not found", errors.Params{"coindId": coinId})
	}

	c, err := getCoinBySymbol(coinMap, coinObj.Symbol, token)
	if err != nil {
		return c, err
	}

	return c, nil
}


func getCoinBySymbol(coinMap map[string]Coin, symbol, token string) (Coin, error) {
	c, ok := coinMap[createID(symbol, token)]
	if !ok {
		return c, errors.E("No coin found by symbol", errors.Params{"symbol": symbol, "token": token})
	}
	return c, nil
}

func normalizeCharts(c Charts) charts.Data {
	chartsData := charts.Data{}
	prices := make([]charts.Price, 0)
	for _, quote := range c.Prices {
		if len(quote) != chartDataSize {
			continue
		}

		date := time.Unix(int64(quote[0])/1000, 0)
		prices = append(prices, charts.Price{
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

func normalizeInfo(data CoinPrice) charts.CoinDetails {
	return charts.CoinDetails{
		Vol24:             data.TotalVolume,
		MarketCap:         data.MarketCap,
		CirculatingSupply: data.CirculatingSupply,
		TotalSupply:       data.TotalSupply,
		Provider:          id,
	}
}
