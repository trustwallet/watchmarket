package memory

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type Instance struct {
	id string
	*gocache.Cache
}

const id = "memory"

func Init() Instance {
	return Instance{id: id, Cache: gocache.New(time.Hour, time.Hour*10)}
}

func (i Instance) GetID() string {
	return i.id
}

func (i Instance) GenerateKey(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (i Instance) Set(key string, data []byte) error {
	i.Cache.Set(key, data, gocache.NoExpiration)
	return nil
}

func (i Instance) Get(key string) ([]byte, error) {
	res, ok := i.Cache.Get(key)
	if !ok {
		return nil, errors.New(watchmarket.ErrNotFound)
	}
	return res.([]byte), nil
}

func (i Instance) SetWithTime(key string, data []byte, time int64) error {
	return nil
}

func (i Instance) GetWithTime(key string, time int64) ([]byte, error) {
	return nil, nil
}

func (i Instance) GetLenOfSavedItems() int {
	items := i.Cache.Items()
	return len(items)
}
