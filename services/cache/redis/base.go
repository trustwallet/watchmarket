package rediscache

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/trustwallet/watchmarket/redis"
	"time"
)

type Instance struct {
	id             string
	redis          redis.Redis
	chartsCaching  time.Duration
	tickersCaching time.Duration
	ratesCaching   time.Duration
	detailsCaching time.Duration
}

const id = "redis"

func Init(redis redis.Redis, chartsCaching, tickersCaching, ratesCaching, detailsCaching time.Duration) Instance {
	return Instance{id, redis, chartsCaching, tickersCaching, ratesCaching, detailsCaching}
}

func (i Instance) GetID() string {
	return i.id
}

func (i Instance) GenerateKey(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])
}
