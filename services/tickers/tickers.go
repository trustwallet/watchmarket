package ticker

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
)

type Market struct {
	Id string
}

func (m *Market) GetId() string {
	return m.Id
}

func (m *Market) Init() error {
	logger.Info("Init Market Quote Provider", logger.Params{"market": m.GetId()})
	if len(m.Id) == 0 {
		return errors.E("Market Quote: Id cannot be empty")
	}

	return nil
}
