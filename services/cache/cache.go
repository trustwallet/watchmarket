package cache

type (
	Provider interface {
		GetID() string
		GenerateKey(data string) string

		Get(key string) ([]byte, error)
		Set(key string, data []byte) error
		GetWithTime(key string, time int64) ([]byte, error)
		SetWithTime(key string, data []byte, time int64) error
		GetLenOfSavedItems() int
	}
)
