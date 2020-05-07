package charts

import (
	"fmt"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"net/url"
)

const (
	AssetsURL = "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/"
)

type Client struct {
	blockatlas.Request
}

func NewClient(api string) Client {
	return Client{blockatlas.InitClient(api)}
}

func (c Client) GetCoinInfo(coinId int, token string) (Info, error) {
	coinObject, ok := coin.Coins[uint(coinId)]
	if !ok {
		return Info{}, errors.E("coin not found", errors.Params{"coin": coinObject.Handle, "token": token})
	}

	var (
		path   = fmt.Sprintf("%s/info.json", getPathForCoin(coinObject, token))
		result Info
	)

	err := c.Get(&result, path, url.Values{})
	if err != nil {
		return result, err
	}

	return result, nil
}

func getPathForCoin(c coin.Coin, token string) string {
	if len(token) == 0 {
		return AssetsURL + c.Handle + "/info"
	}
	return AssetsURL + c.Handle + "/assets/" + token
}
