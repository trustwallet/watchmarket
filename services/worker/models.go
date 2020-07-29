package worker

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"sync"
)

type (
	tickers struct {
		tickers watchmarket.Tickers
		sync.Mutex
	}

	rates struct {
		rates watchmarket.Rates
		sync.Mutex
	}
)

func (t *tickers) Add(tickers watchmarket.Tickers) {
	t.Lock()
	t.tickers = append(t.tickers, tickers...)
	t.Unlock()
}

func (r *rates) Add(rates watchmarket.Rates) {
	r.Lock()
	r.rates = append(r.rates, rates...)
	r.Unlock()
}

func isHigherPriority(priorities []string, current, new string) bool {
	for _, p := range priorities {
		if p == current {
			return false
		} else if p == new {
			return true
		}
	}
	return false
}
