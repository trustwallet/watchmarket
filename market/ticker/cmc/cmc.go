package cmc

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/address"
	"github.com/trustwallet/watchmarket/market/clients/cmc"
	"github.com/trustwallet/watchmarket/market/ticker"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

const (
	id = "cmc"
)

type Market struct {
	ticker.Market
	mapApi string
	client *cmc.Client
}

func InitMarket(mapApi, updateTime string, client *cmc.Client) ticker.TickerProvider {
	m := &Market{
		Market: ticker.Market{
			Id:         id,
			UpdateTime: updateTime,
		},
		mapApi: mapApi,
		client: client,
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
		tokenId = normalizeTokenId(price.Platform.TokenAddress, price.Platform.Symbol)
		coinType = watchmarket.TypeToken
		coinName = price.Platform.Symbol
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

	for _, result := range cmcCoin {
		coinName = result.Coin.Symbol
		if result.CoinType == watchmarket.TypeCoin {
			tokenId = ""
		} else if len(result.TokenId) > 0 {
			tokenId = normalizeTokenId(result.TokenId, coinName)
		}
		tickers = append(tickers, &watchmarket.Ticker{
			CoinName: coinName,
			CoinType: result.CoinType,
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

func normalizeTokenId(tokenId, symbol string) string {
	switch symbol {
	case coin.Ethereum().Symbol, coin.Classic().Symbol, coin.Poa().Symbol,
		coin.Tomochain().Symbol, coin.Gochain().Symbol, coin.Thundertoken().Symbol,
		coin.Callisto().Symbol, coin.Wanchain().Symbol:
		return address.EIP55Checksum(tokenId)
	default:
		return tokenId
	}
}

func normalizeTickers(prices cmc.CoinPrices, provider string, cmap cmc.CmcMapping) (tickers watchmarket.Tickers) {
	for _, price := range prices.Data {
		t := normalizeTicker(price, provider, cmap)
		tickers = append(tickers, t...)
	}
	return
}
