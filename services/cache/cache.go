package cache

import "context"

type (
	Provider interface {
		GetID() string
		GenerateKey(data string) string

		Get(key string, ctx context.Context) ([]byte, error)
		Set(key string, data []byte, ctx context.Context) error
		GetWithTime(key string, time int64, ctx context.Context) ([]byte, error)
		SetWithTime(key string, data []byte, time int64, ctx context.Context) error
	}
)
