package main

import (
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/storage"
)

const (
	defaultConfigPath = "../../config.yml"
)

var (
	cache    *storage.Storage
)

func init() {
	_, confPath := internal.ParseArgs("", defaultConfigPath)
	internal.InitConfig(confPath)
	logger.InitLogger()

	redisHost := viper.GetString("storage.redis")
	cache = internal.InitRedis(redisHost)
}

func main() {
	market.InitRates(cache)
	market.InitMarkets(cache)
	<-make(chan struct{})
}
