package cache

import (
	"crypto/sha1"
	"encoding/base64"
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

func GenerateKey(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (i Instance) GetRates() {

}
