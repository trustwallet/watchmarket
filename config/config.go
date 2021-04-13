package config

import (
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	Markets struct {
		Priority struct {
			Charts     []string `mapstructure:"charts"`
			CoinInfo   []string `mapstructure:"coin_info"`
			Tickers    []string `mapstructure:"tickers"`
			Rates      []string `mapstructure:"rates"`
			RatesAllow []string `mapstructure:"rates_allow"`
		} `mapstructure:"priority"`
		Coinmarketcap struct {
			API       string `mapstructure:"api"`
			Key       string `mapstructure:"key"`
			Currency  string `mapstructure:"currency"`
			WebAPI    string `mapstructure:"web_api"`
			WidgetAPI string `mapstructure:"widget_api"`
		} `mapstructure:"coinmarketcap"`
		Coingecko struct {
			API      string `mapstructure:"api"`
			Key      string `mapstructure:"key"`
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
		Tickers string `mapstructure:"tickers"`
		Rates   string `mapstructure:"rates"`
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
		UseMemoryCache bool          `mapstructure:"use_memory_cache"`
		UpdateTime     struct {
			Tickers string `mapstructure:"memory_cache_tickers"`
			Rates   string `mapstructure:"memory_cache_rates"`
		} `mapstructure:"update_time"`
	} `mapstructure:"rest_api"`

	Sentry struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"sentry"`
}

func Init(confPath string) (Configuration, error) {
	c := Configuration{}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if confPath != "" {
		confPath, err := filepath.Abs(confPath)
		if err != nil {
			return c, err
		}
		viper.SetConfigFile(confPath)
	}
	if err := viper.ReadInConfig(); err != nil {
		return c, errors.Wrapf(err, "Fatal error reading config")
	}
	log.Info("Viper using config")

	bindEnvs(c)
	if err := viper.Unmarshal(&c); err != nil {
		return c, errors.Wrap(err, "Error Unmarshal Viper Config File")
	}
	return c, nil
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
