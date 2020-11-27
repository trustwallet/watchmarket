package coingecko

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/golibs/coin"

	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (p Provider) GetChartData(coinID uint, token, currency string, timeStart int64, ctx context.Context) (watchmarket.Chart, error) {
	chartsData := watchmarket.Chart{}

	coins, err := p.client.fetchCoins(ctx)
	if err != nil {
		return chartsData, err
	}

	symbolsMap := createSymbolsMap(coins)

	coinResult, err := getCoinByID(symbolsMap, coinID, token)
	if err != nil {
		return chartsData, err
	}

	timeEndDate := time.Now().Unix()

	c, err := p.client.fetchCharts(coinResult.Id, currency, timeStart, timeEndDate, ctx)
	if err != nil {
		return chartsData, err
	}

	return normalizeCharts(c), nil
}

func (p Provider) GetCoinData(coinID uint, token, currency string, ctx context.Context) (watchmarket.CoinDetails, error) {
	coins, err := p.client.fetchCoins(ctx)
	if err != nil {
		return watchmarket.CoinDetails{}, err
	}

	symbolsMap := createSymbolsMap(coins)

	coinResult, err := getCoinByID(symbolsMap, coinID, token)
	if err != nil {
		return watchmarket.CoinDetails{}, err
	}

	ratesData := p.client.fetchRates(Coins{coinResult}, currency, ctx)
	if len(ratesData) == 0 {
		return watchmarket.CoinDetails{}, errors.New("no rates found")
	}

	infoData, err := p.info.GetCoinInfo(coinID, token, ctx)
	if err != nil {
		log.WithFields(log.Fields{"coin": coinID, "token": token}).Warn("No assets assets about that coin")
	}
	return watchmarket.CoinDetails{
		Info:        &infoData,
		Provider:    id,
		ProviderURL: ratesData[0].getUrl(),
	}, nil
}

func createSymbolsMap(coins Coins) map[string]Coin {
	var (
		symbolsMap = make(map[string]Coin)
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
		return c, errors.New("coin not found")
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
		return c, errors.New("no coin found by symbol")
	}
	return c, nil
}

func normalizeCharts(c Charts) watchmarket.Chart {
	prices := make([]watchmarket.ChartPrice, 0)
	for _, quote := range c.Prices {
		if len(quote) == chartDataSize {
			prices = append(prices, watchmarket.ChartPrice{
				Price: quote[1],
				Date:  time.Unix(int64(quote[0])/1000, 0).Unix(),
			})
		}
	}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Date < prices[j].Date
	})

	return watchmarket.Chart{
		Prices:   prices,
		Provider: id,
	}
}
