package controllers

import "github.com/trustwallet/watchmarket/pkg/watchmarket"

func (c Controller) HandleChartsRequest() (watchmarket.ChartPrice, error) {
	return watchmarket.ChartPrice{}, nil
}

func (c Controller) verifyRequestData() {

}
