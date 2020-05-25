package redis

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"go.elastic.co/apm/module/apmgoredis"
	"time"
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

func (db Redis) Get(key string, ctx context.Context) ([]byte, error) {
	client := apmgoredis.Wrap(&db.client).WithContext(ctx)
	cmd := client.Get(key)
	if cmd.Err() == redis.Nil {
		return nil, errors.E("Not found", errors.Params{"key": key})
	} else if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return []byte(cmd.Val()), nil
}

func (db Redis) Set(key string, value []byte, expiration time.Duration, ctx context.Context) error {
	client := apmgoredis.Wrap(&db.client).WithContext(ctx)
	cmd := client.Set(key, value, expiration)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (db Redis) Delete(key string, ctx context.Context) error {
	client := apmgoredis.Wrap(&db.client).WithContext(ctx)
	cmd := client.Del(key)
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
