package endpoint

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
)

// @Summary Get rate
// @Id get_rate
// @Description Get rate
// @Accept json
// @Produce json
// @Tags Rates
// @Param from query string false "From" default(USD)
// @Param to query string false "To" default(RUB)
// @Param amount query string false "Amount" default(100)
// @Success 200 {object} controllers.RateResponse
// @Router /v1/market/rate [get]
func GetRate(controller controllers.RatesController) func(context *gin.Context) {
	return func(c *gin.Context) {
		from := c.DefaultQuery("from", watchmarket.DefaultCurrency)
		to := c.Query("to")
		amount, err := strconv.ParseFloat(c.Query("amount"), 64)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		request := controllers.RateRequest{From: from, To: to, Amount: amount}

		response, err := controller.HandleRatesRequest(request)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

// @Summary Get Fiat Rates
// @Description Get Fiat Rates
// @Accept json
// @Produce json
// @Tags Rates
// @Success 200 {object} controllers.FiatRates
// @Router /v1/fiat_rates [get]
func GetFiatRates(controller controllers.RatesController) func(context *gin.Context) {
	return func(context *gin.Context) {
		rates, err := controller.GetFiatRates()
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			context.AbortWithStatusJSON(code, response)
			return
		}
		context.JSON(http.StatusOK, rates)
	}
}
