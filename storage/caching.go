package storage

const EntityCache = "MARKET_CACHE"

func (s *Storage) Set(key string, data []byte) error {
	err := s.AddHM(EntityCache, key, &data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Get(key string) ([]byte, error) {
	var cd []byte
	err := s.GetHMValue(EntityCache, key, &cd)
	if err != nil {
		return nil, err
	}
	return cd, nil
}

func (s *Storage) Delete(key string) error {
	err := s.DeleteHM(EntityCache, key)
	if err != nil {
		return err
	}
	return nil
}
