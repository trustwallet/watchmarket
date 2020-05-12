package cache

import (
	"github.com/trustwallet/watchmarket/redis"
	"time"
)

type Instance struct {
	chartsCaching time.Duration
	redis         redis.Redis
}

func Init(redis redis.Redis, chartsCaching time.Duration) Instance {
	return Instance{redis: redis, chartsCaching: chartsCaching}
}

func (i Instance) GetTickers() {

}

func (i Instance) GetRates() {

}
