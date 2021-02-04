package assets

import (
	"errors"
	"fmt"
	"github.com/trustwallet/watchmarket/services/controllers"
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

func (c Client) GetCoinInfo(asset controllers.Asset) (info watchmarket.Info, err error) {
	coinObject, ok := coin.Coins[asset.CoinId]
	if !ok {
		err = errors.New(fmt.Sprint("coin not found ", asset.CoinId, "; token ", asset.TokenId))
		return
	}

	path := fmt.Sprintf("/%s/info.json", getPathForCoin(coinObject, asset.TokenId))
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
