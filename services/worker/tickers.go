package worker

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/markets"
)

func (w Worker) FetchAndSaveTickers() {
	log.Info("Fetching Tickers ...")
	fetchedTickers := fetchTickers(w.tickersApis)
	normalizedTickers := toTickersModel(fetchedTickers)

	if err := w.db.AddTickers(normalizedTickers); err != nil {
		log.Error(err)
	}
}

func fetchTickers(tickersApis markets.TickersAPIs) watchmarket.Tickers {
	wg := new(sync.WaitGroup)
	s := new(tickers)
	for _, t := range tickersApis {
		wg.Add(1)
		go fetchTickersByProvider(t, wg, s)
	}
	wg.Wait()

	return s.tickers
}

func fetchTickersByProvider(t markets.TickersAPI, wg *sync.WaitGroup, s *tickers) {
	defer wg.Done()

	tickers, err := t.GetTickers()
	if err != nil {
		log.WithFields(log.Fields{"provider": t.GetProvider(), "details": err}).Error("Failed to fetch tickers")
		return
	}

	log.WithFields(log.Fields{"provider": t.GetProvider(), "tickers": len(tickers)}).Info("Tickers fetching done")

	s.Add(tickers)
}
