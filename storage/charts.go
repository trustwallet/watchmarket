package storage

import "time"

const (
	EntityCharts              = "ATLAS_CHARTS"
	defaultChartsCacheTimeout = 600
)

func (s *Storage) SaveCharts(key string, data *ChartData) (SaveResult, error) {
	err := s.AddHM(EntityCharts, key, &data)
	if err != nil {
		return SaveResultAddHMFailure, err
	}
	return SaveResultSuccess, err
}

func (s *Storage) GetCharts(key string) (*ChartData, error) {
	var data ChartData
	err := s.GetHMValue(EntityCharts, key, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (data *ChartData) IsOutdated() bool {
	timeNow := time.Now().Unix()
	return timeNow-data.Timestamp > defaultChartsCacheTimeout
}

func (data *ChartData) IsEmpty() bool {
	return len(data.Prices) == 0
}
