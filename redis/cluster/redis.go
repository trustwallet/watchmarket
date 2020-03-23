package cluster

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type Redis struct {
	client *redis.ClusterClient
}

func (db *Redis) Init(host string) error {
	return errors.E("Use single redis package for that")
}

func (db *Redis) InitCluster(host []string) error {
	if len(host) == 0 {
		return errors.E("Empty host")
	}

	var (
		addresses []string
		password  string
	)

	for _, h := range host {
		url, err := redis.ParseURL(h)
		if err != nil {
			return err
		}
		addresses = append(addresses, url.Addr)
		password = url.Password
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addresses,
		Password: password,
	})

	if err := client.Ping().Err(); err != nil {
		return err
	}

	db.client = client
	return nil
}

func (db *Redis) GetValue(key string, value interface{}) error {
	cmd := db.client.Get(key)
	if cmd.Err() == redis.Nil {
		return watchmarket.ErrNotFound
	} else if cmd.Err() != nil {
		return cmd.Err()
	}
	err := json.Unmarshal([]byte(cmd.Val()), value)
	if err != nil {
		return err
	}
	return nil
}

func (db *Redis) Add(key string, value interface{}) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cmd := db.client.Set(key, j, 0)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (db *Redis) Delete(key string) error {
	cmd := db.client.Del(key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (db *Redis) IsReady() bool {
	if db.client == nil {
		return false
	}
	return db.client.Ping().Err() == nil
}
