package assets

import (
	"errors"
	"fmt"
	"github.com/trustwallet/watchmarket/services/controllers"

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

func (c Client) GetCoinInfo(asset controllers.Asset) (watchmarket.Info, error) {
	coinObject, ok := coin.Coins[asset.CoinId]
	if !ok {
		return watchmarket.Info{}, errors.New("coin not found " + "token " + asset.TokenId)
	}

	var (
		path   = c.api + fmt.Sprintf("/%s/info.json", getPathForCoin(coinObject, asset.TokenId))
		result watchmarket.Info
	)

	resp, err := c.r.Get(path)
	if err != nil {
		return watchmarket.Info{}, err
	}
	err = resp.ToJSON(&result)
	if err != nil {
		log.WithFields(log.Fields{
			"url":      resp.Request().URL.String(),
			"status":   resp.Response().Status,
			"response": resp,
		}).Error("Assets Get Coin Info: ", resp.Response().Status)

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
