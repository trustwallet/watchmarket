package assets

import (
	"context"
	"errors"
	"fmt"

	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type Client struct {
	api string
	r   *req.Req
}

func Init(api string) Client {
	return Client{r: req.New(), api: api}
}

func (c Client) GetCoinInfo(coinId uint, token string, ctx context.Context) (watchmarket.Info, error) {
	coinObject, ok := coin.Coins[coinId]
	if !ok {
		return watchmarket.Info{}, errors.New("coin not found " + "token " + token)
	}

	var (
		path   = c.api + fmt.Sprintf("/%s/info.json", getPathForCoin(coinObject, token))
		result watchmarket.Info
	)

	resp, err := c.r.Get(path, ctx)
	if err != nil {
		return watchmarket.Info{}, err
	}
	err = resp.ToJSON(&result)
	if err != nil {
		log.Error("URL: " + resp.Request().URL.String())
		log.Error("Status code: " + resp.Response().Status)
		return watchmarket.Info{}, err
	}
	return result, nil
}

func getPathForCoin(c coin.Coin, token string) string {
	if len(token) == 0 {
		return c.Handle + "/info"
	}
	return c.Handle + "/assets/" + token
}
