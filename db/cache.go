package db

func (i *Instance) AddCache(key string, data []byte) error {
	return nil
}

func (i *Instance) GetCache(key string) ([]byte, error) {
	return nil, nil
}

func (i *Instance) GetCachingKeyForInterval(intervalKey string, time int64) (string, error) {
	return "", nil
}

func (i *Instance) AddInterval(intervalKey string, time int64) error {
	return nil
}
