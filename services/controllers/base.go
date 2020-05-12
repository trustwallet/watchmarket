package controllers

import (
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/redis"
	"github.com/trustwallet/watchmarket/services/markets"
)

type Controller struct {
	database db.Instance
	redis    redis.Redis
	api      markets.APIs
}

func NewController(database db.Instance, redis redis.Redis, api markets.APIs) Controller {
	return Controller{database, redis, api}
}
