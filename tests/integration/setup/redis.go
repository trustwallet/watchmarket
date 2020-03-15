package setup

import (
	"errors"
	"fmt"
	"github.com/ory/dockertest"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/storage"
	"log"
)

var (
	Cache         *storage.Storage
	redisResource *dockertest.Resource
)

func runRedisContainer() error {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	redisResource, err = pool.Run("redis", "latest", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		Cache = internal.InitRedis(fmt.Sprintf("redis://localhost:%s", redisResource.GetPort("6379/tcp")))
		if Cache == nil {
			return errors.New("failed to init cache")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func stopRedisContainer() error {
	return redisResource.Close()
}

func RunRedisContainer() {
	if err := runRedisContainer(); err != nil {
		log.Fatal(err)
	}
	if Cache == nil {
		log.Fatal(errors.New("failed to init cache"))
	}
}

func StopRedisContainer() {
	if err := stopRedisContainer(); err != nil {
		log.Fatal(err)
	}
}
