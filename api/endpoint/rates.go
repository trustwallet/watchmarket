package endpoint

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"go.elastic.co/apm"
	"net/http"
	"strconv"
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
// @Router /v1/market/fiat/rate [get]
func GetRate(controller controllers.RatesController) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v1/market/fiat/rate ", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		from := c.DefaultQuery("from", watchmarket.DefaultCurrency)
		to := c.Query("to")
		amount, err := strconv.ParseFloat(c.Query("amount"), 64)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		request := controllers.RateRequest{From: from, To: to, Amount: amount}

		response, err := controller.HandleRatesRequest(request, ctx)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}
