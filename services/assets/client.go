package assets

import (
	"errors"
	"fmt"
	"time"

	"github.com/trustwallet/golibs/client"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/golibs/network/middleware"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type Client struct {
	client.Request
}

func Init(api string) Client {
	return Client{client.InitClient(api, middleware.SentryErrorHandler)}
}

func (c Client) GetCoinInfo(coinId uint, token string) (info watchmarket.Info, err error) {
	coinObject, ok := coin.Coins[coinId]
	if !ok {
		err = errors.New("coin not found " + "token " + token)
		return
	}

	path := fmt.Sprintf("/%s/info.json", getPathForCoin(coinObject, token))
	err = c.GetWithCache(&info, path, nil, time.Hour*12)
	//asset info file now only contains description field.
	info.ShortDescription = info.Description
	return
}

func getPathForCoin(c coin.Coin, token string) string {
	if len(token) == 0 {
		return c.Handle + "/info"
	}
	return c.Handle + "/assets/" + token
}
