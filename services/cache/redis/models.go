package rediscache

type CachedInterval struct {
	Timestamp int64
	Duration  int64
	Key       string
}
