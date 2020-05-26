package rediscache

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/redis"
	"time"
)

type Instance struct {
	id            string
	redis         redis.Redis
	cachingPeriod time.Duration
}

const id = "redis"

func Init(redis redis.Redis, cachingPeriod time.Duration) Instance {
	return Instance{id, redis, cachingPeriod}
}

func (i Instance) GetID() string {
	return i.id
}

func (i Instance) GenerateKey(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (i Instance) Get(key string, ctx context.Context) ([]byte, error) {
	raw, err := i.redis.Get(key, ctx)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (i Instance) Set(key string, data []byte, ctx context.Context) error {
	if data == nil {
		return errors.E("data is empty")
	}
	err := i.redis.Set(key, data, i.cachingPeriod, ctx)
	if err != nil {
		return err
	}
	return nil
}
