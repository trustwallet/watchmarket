package assets

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

const (
	AssetsURL = "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/"
)

type AssetClient interface {
	GetCoinInfo(coinId int, token string) (info *watchmarket.CoinInfo, err error)
}

type HttpAssetClient struct {
	HttpClient *resty.Client
}

func (cl *HttpAssetClient) GetCoinInfo(coinId int, token string) (*watchmarket.CoinInfo, error) {
	c, ok := coin.Coins[uint(coinId)]
	if !ok {
		return nil, watchmarket.ErrNotFound
	}
	url := fmt.Sprintf("%s/info.json", getCoinInfoUrl(c, token))
	resp, err := cl.HttpClient.R().
		EnableTrace().
		SetResult(&watchmarket.CoinInfo{}).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 404 {
		return nil, watchmarket.ErrNotFound
	}

	if !resp.IsSuccess() {
		return nil, errors.New(fmt.Sprintf("Request to %s failed with HTTP %d: %s", url, resp.StatusCode(), resp.String()))
	}

	var info watchmarket.CoinInfo

	// TODO: cover this in tests
	err = json.Unmarshal(resp.Body(), &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func getCoinInfoUrl(c coin.Coin, token string) string {
	if len(token) == 0 {
		return AssetsURL + c.Handle + "/info"
	}
	return AssetsURL + c.Handle + "/assets/" + token
}
