package compound

import (
	c "github.com/trustwallet/watchmarket/market/clients/compound"
	"github.com/trustwallet/watchmarket/market/rate"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
	"time"
)

const (
	compound = "compound"
)

type Compound struct {
	rate.Rate
	client *c.Client
}

func InitRate(api string, updateTime string) rate.Provider {
	return &Compound{
		Rate: rate.Rate{
			Id:         compound,
			UpdateTime: updateTime,
		},
		client: c.NewClient(api),
	}
}

func (c *Compound) FetchLatestRates() (rates watchmarket.Rates, err error) {
	coinPrices, err := c.client.GetData()
	if err != nil {
		return
	}
	rates = normalizeRates(coinPrices, c.GetId())
	return
}

func normalizeRates(coinPrices c.CoinPrices, provider string) (rates watchmarket.Rates) {
	for _, cToken := range coinPrices.Data {
		rates = append(rates, watchmarket.Rate{
			Currency:  strings.ToUpper(cToken.Symbol),
			Rate:      1.0 / cToken.UnderlyingPrice.Value,
			Timestamp: time.Now().Unix(),
			Provider:  provider,
		})
	}
	return
}
