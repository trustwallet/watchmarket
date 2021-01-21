package rediscache

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"time"

	"github.com/trustwallet/watchmarket/redis"
)

type Instance struct {
	id            string
	redis         redis.Redis
	cachingPeriod time.Duration
}

const id = "redis"

func Init(host string, cachingPeriod time.Duration) (Instance, error) {
	c, err := redis.Init(host)
	if err != nil {
		return Instance{}, err
	}
	return Instance{id, c, cachingPeriod}, nil
}

func (i Instance) GetID() string {
	return i.id
}

func (i Instance) GenerateKey(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (i Instance) Get(key string) ([]byte, error) {
	raw, err := i.redis.Get(key)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (i Instance) Set(key string, data []byte) error {
	if data == nil {
		return errors.New("data is empty")
	}
	err := i.redis.Set(key, data, i.cachingPeriod)
	if err != nil {
		return err
	}
	return nil
}

func (i Instance) GetLenOfSavedItems() int {
	return 0
}
