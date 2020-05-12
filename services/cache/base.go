package cache

import (
	"github.com/trustwallet/watchmarket/redis"
)

type Instance struct {
	cachingDuration uint
	redis           redis.Redis
}

func Init(redis redis.Redis) Instance {
	return Instance{redis: redis}
}

func (i Instance) GetTickers() {

}

func (i Instance) GetRates() {

}
