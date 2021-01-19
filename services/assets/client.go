package assets

import (
	"errors"
	"fmt"
	"strings"

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

func (c Client) GetCoinInfo(coinId uint, token string) (watchmarket.Info, error) {
	coinObject, ok := coin.Coins[coinId]
	if !ok {
		return watchmarket.Info{}, errors.New("coin not found " + "token " + token)
	}

	var (
		path   = c.api + fmt.Sprintf("/%s/info.json", getPathForCoin(coinObject, token))
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
	result = normalize(result)

	return result, nil
}

func getPathForCoin(c coin.Coin, token string) string {
	if len(token) == 0 {
		return c.Handle + "/info"
	}
	return c.Handle + "/assets/" + token
}

func normalize(info watchmarket.Info) watchmarket.Info {
	info.Website = normalizeUrl(info.Website)
	info.Research = normalizeUrl(info.Research)
	info.WhitePaper = normalizeUrl(info.WhitePaper)
	info.Explorer = normalizeUrl(info.Explorer)
	return info
}

func normalizeUrl(url string) string {
	return strings.Replace(url, "www.", "", -1)
}
