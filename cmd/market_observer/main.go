package main

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/market"
	marketprovider "github.com/trustwallet/watchmarket/market/market"
	marketcmc "github.com/trustwallet/watchmarket/market/market/cmc"
	marketcoingecko "github.com/trustwallet/watchmarket/market/market/coingecko"
	marketcompound "github.com/trustwallet/watchmarket/market/market/compound"
	marketdex "github.com/trustwallet/watchmarket/market/market/dex"
	rateprovider "github.com/trustwallet/watchmarket/market/rate"
	ratecmc "github.com/trustwallet/watchmarket/market/rate/cmc"
	ratecoingecko "github.com/trustwallet/watchmarket/market/rate/coingecko"
	ratecompound "github.com/trustwallet/watchmarket/market/rate/compound"
	ratefixer "github.com/trustwallet/watchmarket/market/rate/fixer"
	"github.com/trustwallet/watchmarket/storage"
	"time"
)

const (
	defaultConfigPath              = "../../config.yml"
	gracefulShutdownTimeoutSeconds = 1
)

var (
	cache           *storage.Storage
	rateProviders   *rateprovider.Providers
	marketProviders *marketprovider.Providers
)

func init() {
	_, confPath := internal.ParseArgs("", defaultConfigPath)
	internal.InitConfig(confPath)
	logger.InitLogger()

	redisHost := viper.GetString("storage.redis")
	cache = internal.InitRedis(redisHost)

	rateProviders = &rateprovider.Providers{
		// Add Market Quote Providers:
		0: ratecmc.InitRate(
			viper.GetString("market.cmc.api"),
			viper.GetString("market.cmc.api_key"),
			viper.GetString("market.cmc.map_url"),
			viper.GetString("market.rate_update_time"),
		),
		1: ratefixer.InitRate(
			viper.GetString("market.fixer.api"),
			viper.GetString("market.fixer.api_key"),
			viper.GetString("market.fixer.rate_update_time"),
		),
		2: ratecompound.InitRate(
			viper.GetString("market.compound.api"),
			viper.GetString("market.rate_update_time"),
		),
		3: ratecoingecko.InitRate(
			viper.GetString("market.coingecko.api"),
			viper.GetString("market.rate_update_time"),
		),
	}

	marketProviders = &marketprovider.Providers{
		// Add Market Quote Providers:
		0: marketcmc.InitMarket(
			viper.GetString("market.cmc.api"),
			viper.GetString("market.cmc.api_key"),
			viper.GetString("market.cmc.map_url"),
			viper.GetString("market.quote_update_time"),
		),
		1: marketcompound.InitMarket(
			viper.GetString("market.compound.api"),
			viper.GetString("market.quote_update_time"),
		),
		2: marketcoingecko.InitMarket(
			viper.GetString("market.coingecko.api"),
			viper.GetString("market.quote_update_time"),
		),
		3: marketdex.InitMarket(
			viper.GetString("market.dex.api"),
			viper.GetString("market.dex.quote_update_time"),
		),
	}
}

func main() {
	rateCron := market.InitRates(cache, rateProviders)
	defer gracefullyShutDown(rateCron)
	rateCron.Start()
	marketCron := market.InitMarkets(cache, marketProviders)
	defer gracefullyShutDown(marketCron)
	marketCron.Start()
	internal.WaitingForExitSignal()
	logger.Info("Waiting for all observer jobs to stop")
}

func gracefullyShutDown(job *cron.Cron) {
	c := job.Stop()
	select {
	case <-time.After(gracefulShutdownTimeoutSeconds * time.Second):
	case <-c.Done():
	}
}
