package coinmarketcap

import (
	"github.com/trustwallet/watchmarket/services/rates"
	tickersClient "github.com/trustwallet/watchmarket/services/tickers/providers/coinnmarketcap"
	"math/big"
)

const (
	id = "coinnmarketcap"
)

type Parser struct {
	ID       string
	client   tickersClient.Client
	currency string
}

func InitParser(api, key, currency string) Parser {
	return Parser{
		ID:       id,
		client:   tickersClient.NewClient(api, key),
		currency: currency,
	}
}

func (p Parser) GetData() (rates rates.Rates, err error) {
	prices, err := p.client.GetData(p.currency)
	if err != nil {
		return
	}
	rates = normalizeRates(prices, p.ID)
	return
}

func normalizeRates(prices tickersClient.CoinPrices, provider string) rates.Rates {
	var res rates.Rates

	for _, price := range prices.Data {
		if price.Platform != nil {
			continue
		}
		res = append(res, rates.Rate{
			Currency:         price.Symbol,
			Rate:             1.0 / price.Quote.USD.Price,
			Timestamp:        price.LastUpdated.Unix(),
			PercentChange24h: *big.NewFloat(price.Quote.USD.PercentChange24h),
			Provider:         provider,
		})
	}
	return res
}
