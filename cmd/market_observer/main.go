package main

import (
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/storage"
)

const (
	defaultConfigPath = "../../config.yml"
)

var (
	confPath string
	cache    *storage.Storage
)

func init() {
	_, confPath, _, cache = internal.InitAPIWithRedis("", defaultConfigPath)
}

func main() {
	market.InitRates(cache)
	market.InitMarkets(cache)
	<-make(chan struct{})
}
