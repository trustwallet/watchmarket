package controllers

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
	"strconv"
	"strings"
)

type Controller struct {
	dataCache        cache.Provider
	database         db.Instance
	chartsPriority   []string
	coinInfoPriority []string
	ratesPriority    []string
	tickersPriority  []string
	api              markets.ChartsAPIs
	configuration    config.Configuration
}

func NewController(
	cache cache.Provider,
	database db.Instance,
	chartsPriority, coinInfoPriority, ratesPriority, tickersPriority []string,
	api markets.ChartsAPIs,
	configuration config.Configuration,
) Controller {
	return Controller{
		cache,
		database,
		chartsPriority,
		coinInfoPriority,
		ratesPriority,
		tickersPriority,
		api,
		configuration,
	}
}

func ParseID(id string) (uint, string, watchmarket.CoinType, error) {
	rawResult := strings.Split(id, "_")
	resLen := len(rawResult)
	if !(resLen > 0 && resLen <= 2) {
		return 0, "", watchmarket.Coin, errors.E("Bad ID")
	}

	coin, err := strconv.Atoi(rawResult[0])
	if err != nil {
		return 0, "", watchmarket.Coin, errors.E("Bad coin")
	}

	if resLen == 1 || rawResult[1] == "" {
		return uint(coin), "", watchmarket.Coin, nil
	}

	return uint(coin), rawResult[1], watchmarket.Token, nil
}

func BuildID(coin uint, token string) string {
	c := strconv.Itoa(int(coin))
	if token != "" {
		return c + "_" + token
	}
	return c
}
