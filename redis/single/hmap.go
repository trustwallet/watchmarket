package single

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (db *Redis) GetAllHM(entity string) (map[string]string, error) {
	cmd := db.client.HGetAll(entity)
	if cmd.Err() == redis.Nil {
		return nil, watchmarket.ErrNotFound
	} else if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return cmd.Val(), nil
}

func (db *Redis) GetHMValue(entity, key string, value interface{}) error {
	cmd := db.client.HMGet(entity, key)
	if cmd.Err() == redis.Nil {
		return watchmarket.ErrNotFound
	} else if cmd.Err() != nil {
		return cmd.Err()
	}
	val, ok := cmd.Val()[0].(string)
	if !ok {
		return watchmarket.ErrNotFound
	}
	err := json.Unmarshal([]byte(val), value)
	if err != nil {
		return err
	}
	return nil
}

func (db *Redis) AddHM(entity, key string, value interface{}) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cmd := db.client.HMSet(entity, map[string]interface{}{key: j})
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (db *Redis) DeleteHM(entity, key string) error {
	cmd := db.client.HDel(entity, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
