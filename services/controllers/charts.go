package controllers

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
)

func (c Controller) HandleChartsRequest(coinQuery, token, currency, timeStartRaw, maxItems string) (watchmarket.Chart, error) {
	coinId, _ := strconv.Atoi(coinQuery)
	timeStart, _ := strconv.ParseInt(timeStartRaw, 10, 64)
	provider := c.chartsPriority.GetCurrentProvider()
	price, err := c.api.ChartsAPIs[provider].GetChartData(uint(coinId), token, currency, timeStart)
	if err != nil {
		return watchmarket.Chart{}, err
	}
	return price, nil
}

func (c Controller) verifyRequestData() {

}
