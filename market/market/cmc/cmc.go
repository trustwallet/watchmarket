package cmc

import (
	"github.com/trustwallet/watchmarket/market/clients/cmc"
	"github.com/trustwallet/watchmarket/market/market"
	watchmarket "github.com/trustwallet/watchmarket/pkg/watchmarket"
)

const (
	id = "cmc"
)

type Market struct {
	market.Market
	mapApi string
	client *cmc.Client
}

func InitMarket(api string, apiKey string, mapApi string, updateTime string) market.Provider {
	m := &Market{
		Market: market.Market{
			Id:         id,
			UpdateTime: updateTime,
		},
		mapApi: mapApi,
		client: cmc.NewClient(api, apiKey),
	}
	return m
}

func (m *Market) GetData() (watchmarket.Tickers, error) {
	cmap, err := cmc.GetCmcMap(m.mapApi)
	if err != nil {
		return nil, err
	}
	prices, err := m.client.GetData()
	if err != nil {
		return nil, err
	}
	return normalizeTickers(prices, m.GetId(), cmap), nil
}

func normalizeTicker(price cmc.Data, provider string, cmap cmc.CmcMapping) (tickers watchmarket.Tickers) {
	tokenId := ""
	coinName := price.Symbol
	coinType := watchmarket.TypeCoin
	if price.Platform != nil {
		coinType = watchmarket.TypeToken
		coinName = price.Platform.Symbol
		tokenId = price.Platform.TokenAddress
		if len(tokenId) == 0 {
			tokenId = price.Symbol
		}
	}

	cmcCoin, err := cmap.GetCoins(price.Id)
	if err != nil {
		tickers = append(tickers, &watchmarket.Ticker{
			CoinName: coinName,
			CoinType: coinType,
			TokenId:  tokenId,
			Price: watchmarket.TickerPrice{
				Value:     price.Quote.USD.Price,
				Change24h: price.Quote.USD.PercentChange24h,
				Currency:  watchmarket.DefaultCurrency,
				Provider:  provider,
			},
			LastUpdate: price.LastUpdated,
		})
		return
	}

	for _, cmc := range cmcCoin {
		coinName = cmc.Coin.Symbol
		if cmc.CoinType == watchmarket.TypeCoin {
			tokenId = ""
		} else if len(cmc.TokenId) > 0 {
			tokenId = cmc.TokenId
		}
		tickers = append(tickers, &watchmarket.Ticker{
			CoinName: coinName,
			CoinType: cmc.CoinType,
			TokenId:  tokenId,
			Price: watchmarket.TickerPrice{
				Value:     price.Quote.USD.Price,
				Change24h: price.Quote.USD.PercentChange24h,
				Currency:  watchmarket.DefaultCurrency,
				Provider:  provider,
			},
			LastUpdate: price.LastUpdated,
		})
	}
	return
}

func normalizeTickers(prices cmc.CoinPrices, provider string, cmap cmc.CmcMapping) (tickers watchmarket.Tickers) {
	for _, price := range prices.Data {
		t := normalizeTicker(price, provider, cmap)
		tickers = append(tickers, t...)
	}
	return
}
