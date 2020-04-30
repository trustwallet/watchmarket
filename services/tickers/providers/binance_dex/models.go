package binance_dex

type CoinPrice struct {
	BaseAssetName      string `json:"baseAssetName"`
	QuoteAssetName     string `json:"quoteAssetName"`
	PriceChangePercent string `json:"priceChangePercent"`
	LastPrice          string `json:"lastPrice"`
}
