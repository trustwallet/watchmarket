package cache

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/trustwallet/watchmarket/redis"
	"time"
)

type Instance struct {
	redis          redis.Redis
	chartsCaching  time.Duration
	tickersCaching time.Duration
	ratesCaching   time.Duration
	detailsCaching time.Duration
}

func Init(redis redis.Redis, chartsCaching, tickersCaching, ratesCaching, detailsCaching time.Duration) Instance {
	return Instance{redis, chartsCaching, tickersCaching, ratesCaching, detailsCaching}
}

func GenerateKey(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (i Instance) GetRates() {

}
