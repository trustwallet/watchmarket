package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	client redis.Client
}

func Init(host string) (Redis, error) {
	options, err := redis.ParseURL(host)
	if err != nil {
		return Redis{}, err
	}
	client := redis.NewClient(options)
	if err := client.Ping().Err(); err != nil {
		return Redis{}, err
	}

	return Redis{client: *client}, nil
}

func (db Redis) Get(key string) ([]byte, error) {
	cmd := db.client.Get(key)
	if cmd.Err() == redis.Nil {
		return nil, errors.New("not found")
	} else if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return []byte(cmd.Val()), nil
}

func (db Redis) Set(key string, value []byte, expiration time.Duration) error {
	cmd := db.client.Set(key, value, expiration)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (db Redis) Delete(key string) error {
	cmd := db.client.Del(key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (db Redis) IsAvailable() bool {
	return db.client.Ping().Err() == nil
}

func (db Redis) Reconnect(host string) bool {
	options, err := redis.ParseURL(host)
	if err != nil {
		return false
	}
	client := redis.NewClient(options)
	if err := client.Ping().Err(); err != nil {
		return false
	}
	db.client = *client
	if err := db.client.Ping().Err(); err != nil {
		return false
	}
	return true
}
