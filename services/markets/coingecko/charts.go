package coingecko

import (
	"errors"
	"github.com/trustwallet/watchmarket/services/controllers"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/golibs/coin"

	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (p Provider) GetChartData(asset controllers.Asset, currency string, timeStart int64) (watchmarket.Chart, error) {
	chartsData := watchmarket.Chart{}

	coins, err := p.client.fetchCoins()
	if err != nil {
		return chartsData, err
	}

	symbolsMap := createSymbolsMap(coins)

	coinResult, err := getCoinByID(symbolsMap, asset)
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

func (p Provider) GetCoinData(asset controllers.Asset, currency string) (watchmarket.CoinDetails, error) {
	coins, err := p.client.fetchCoins()
	if err != nil {
		return watchmarket.CoinDetails{}, err
	}

	symbolsMap := createSymbolsMap(coins)

	coinResult, err := getCoinByID(symbolsMap, asset)
	if err != nil {
		return watchmarket.CoinDetails{}, err
	}

	ratesData := p.client.fetchRates(Coins{coinResult}, currency)
	if len(ratesData) == 0 {
		return watchmarket.CoinDetails{}, errors.New("no rates found")
	}

	infoData, err := p.info.GetCoinInfo(asset)
	if err != nil {
		log.WithFields(log.Fields{"coin": asset.CoinId, "token": asset.TokenId}).Warn("No assets assets about that coin")
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

func getCoinByID(coinMap map[string]Coin, asset controllers.Asset) (Coin, error) {
	c := Coin{}
	coinObj, ok := coin.Coins[asset.CoinId]
	if !ok {
		return c, errors.New("coin not found")
	}

	c, ok = coinMap[createID(coinObj.Symbol, asset.TokenId)]
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
