package rate

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
)

const (
	defaultUpdateTime = "5m"
)

type Rate struct {
	Id         string
	UpdateTime string
}

func (r *Rate) GetUpdateTime() string {
	return r.UpdateTime
}

func (r *Rate) GetId() string {
	return r.Id
}

func (r *Rate) GetLogType() string {
	return "market-rate"
}

func (r *Rate) Init(updateTime string) error {
	logger.Info("Init Market Rate Provider", logger.Params{"rate": r.GetId()})
	if len(r.Id) == 0 {
		return errors.E("Market Rate: Id cannot be empty")
	}

	r.UpdateTime = updateTime

	if len(r.UpdateTime) == 0 {
		r.UpdateTime = defaultUpdateTime
	}
	return nil
}
