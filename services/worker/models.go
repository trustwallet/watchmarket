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
