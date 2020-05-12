package controllers

import "github.com/trustwallet/watchmarket/pkg/watchmarket"

func (c Controller) HandleChartsRequest() (watchmarket.ChartsPrice, error) {
	return watchmarket.ChartsPrice{}, nil
}

func (c Controller) verifyRequestData() {

}
