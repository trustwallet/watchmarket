package compound

import (
	"github.com/trustwallet/blockatlas/coin"
	c "github.com/trustwallet/watchmarket/market/clients/compound"
	"github.com/trustwallet/watchmarket/market/ticker"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"time"
)

const (
	id = "compound"
)

type Market struct {
	ticker.Market
	client *c.Client
}

func InitMarket(api, updateTime string) ticker.TickerProvider {
	m := &Market{
		Market: ticker.Market{
			Id:         id,
			UpdateTime: updateTime,
		},
		client: c.NewClient(api),
	}
	return m
}

func (m *Market) GetData() (result watchmarket.Tickers, err error) {
	coinPrices, err := m.client.GetData()
	if err != nil {
		return
	}
	result = normalizeTickers(coinPrices, m.GetId())
	return result, nil
}

func normalizeTicker(ctoken c.CToken, provider string) (*watchmarket.Ticker, error) {
	// TODO: add value24 calculation
	return &watchmarket.Ticker{
		CoinName: coin.Ethereum().Symbol,
		CoinType: watchmarket.TypeToken,
		TokenId:  ctoken.TokenAddress,
		Price: watchmarket.TickerPrice{
			Value:    ctoken.UnderlyingPrice.Value,
			Currency: coin.Coins[coin.ETH].Symbol,
			Provider: provider,
		},
		LastUpdate: time.Now(),
	}, nil
}

func normalizeTickers(prices c.CoinPrices, provider string) (tickers watchmarket.Tickers) {
	for _, price := range prices.Data {
		t, err := normalizeTicker(price, provider)
		if err != nil {
			continue
		}
		tickers = append(tickers, t)
	}
	return
}
