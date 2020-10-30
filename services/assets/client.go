package assets

import (
	"context"
	"errors"
	"fmt"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type Client struct {
	blockatlas.Request
}

func Init(api string) Client {
	return Client{blockatlas.InitClient(api)}
}

func (c Client) GetCoinInfo(coinId uint, token string, ctx context.Context) (watchmarket.Info, error) {
	coinObject, ok := coin.Coins[coinId]
	if !ok {
		return watchmarket.Info{}, errors.New("coin not found " + "token " + token)
	}

	var (
		path   = fmt.Sprintf("%s/info.json", getPathForCoin(coinObject, token))
		result watchmarket.Info
	)

	err := c.GetWithContext(&result, path, nil, ctx)
	if err != nil {
		return result, err
	}

	return result, nil
}

func getPathForCoin(c coin.Coin, token string) string {
	if len(token) == 0 {
		return c.Handle + "/info"
	}
	return c.Handle + "/assets/" + token
}
