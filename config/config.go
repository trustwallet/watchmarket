package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"reflect"
	"strings"
	"time"
)

type Configuration struct {
	Markets struct {
		Priority struct {
			Charts   []string `mapstructure:"charts"`
			CoinInfo []string `mapstructure:"coin_info"`
			Tickers  []string `mapstructure:"tickers"`
			Rates    []string `mapstructure:"rates"`
		} `mapstructure:"priority"`
		BinanceDex struct {
			API string `mapstructure:"api"`
		} `mapstructure:"binancedex"`
		Coinmarketcap struct {
			API       string `mapstructure:"api"`
			Key       string `mapstructure:"key"`
			Currency  string `mapstructure:"currency"`
			WebAPI    string `mapstructure:"web_api"`
			WidgetAPI string `mapstructure:"widget_api"`
		} `mapstructure:"coinmarketcap"`
		Coingecko struct {
			API      string `mapstructure:"api"`
			Currency string `mapstructure:"currency"`
		} `mapstructure:"coingecko"`
		Fixer struct {
			API      string `mapstructure:"api"`
			Currency string `mapstructure:"currency"`
			Key      string `mapstructure:"key"`
		} `mapstructure:"fixer"`
		Assets string `mapstructure:"assets"`
	} `mapstructure:"markets"`

	Storage struct {
		Redis struct {
			Url string `mapstructure:"url"`
		} `mapstructure:"redis"`
		Postgres struct {
			Url  string `mapstructure:"url"`
			Logs bool   `mapstructure:"logs"`
		} `mapstructure:"postgres"`
	} `mapstructure:"storage"`

	Worker struct {
		Tickers    string `mapstructure:"tickers"`
		Rates      string `mapstructure:"rates"`
		BatchLimit uint   `mapstructure:"batch_limit"`
	} `mapstructure:"worker"`

	RestAPI struct {
		Mode    string `mapstructure:"mode"`
		Port    string `mapstructure:"port"`
		Tickers struct {
			RespsectableMarketCap float64       `mapstructure:"respectable_market_cap"`
			RespsectableVolume    float64       `mapstructure:"respectable_volume"`
			RespectableUpdateTime time.Duration `mapstructure:"respectable_update_time"`
			CacheControl          time.Duration `mapstructure:"cache_control"`
		}
		Charts struct {
			CacheControl time.Duration `mapstructure:"cache_control"`
		} `mapstructure:"charts"`
		Info struct {
			CacheControl time.Duration `mapstructure:"cache_control"`
		} `mapstructure:"info"`
		Cache          time.Duration `mapstructure:"cache"`
		RequestLimit   int           `mapstructure:"request_limit"`
		UseMemoryCache bool          `mapstructure:"use_memory_cache"`
		UpdateTime     struct {
			Tickers string `mapstructure:"memory_cache_tickers"`
			Rates   string `mapstructure:"memory_cache_rates"`
		} `mapstructure:"update_time"`
	} `mapstructure:"rest_api"`
}

func Init(confPath string) Configuration {
	c := Configuration{}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if confPath == "" {
		err := viper.ReadInConfig()
		if err != nil {
			log.Panic(err, "Fatal error reading default config")
		} else {
			log.WithFields(log.Fields{"config": viper.ConfigFileUsed()}).Info("Viper using default config")
		}
	} else {
		viper.SetConfigFile(confPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Panic(err, "Fatal error reading supplied config")
		} else {
			log.WithFields(log.Fields{"config": viper.ConfigFileUsed()}).Info("Viper using supplied config")
		}
	}

	bindEnvs(c)
	if err := viper.Unmarshal(&c); err != nil {
		log.Panic(err, "Error Unmarshal Viper Config File")
	}
	return c
}

func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			if err := viper.BindEnv(strings.Join(append(parts, tv), ".")); err != nil {
				log.Fatal(err)
			}
		}
	}
}
