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

// @Summary Get charts coin assets data for a specific coin
// @Id get_charts_coin_info
// @Description Get the charts coin assets data from an market and coin/contract
// @Accept json
// @Produce json
// @Tags Market
// @Param coin query int true "Coin id" default(60)
// @Param token query string false "Token id"
// @Param currency query string false "The currency to show coin assets in" default(USD)
// @Success 200 {object} watchmarket.ChartCoinInfo
// @Router /v1/market/info [get]
func GetCoinInfoHandler(controller controllers.Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v1/market/info", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		request := controllers.DetailsRequest{
			CoinQuery: c.Query("coin"),
			Token:     c.Query("token"),
			Currency:  c.DefaultQuery("currency", watchmarket.DefaultCurrency),
		}
		response, err := controller.HandleDetailsRequest(request, ctx)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCoinInfoHandlerV2(controller controllers.Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v1/market/info/:id", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		coin, token, _, err := controllers.ParseID(c.Param("id"))
		if err != nil {
			handleError(c, err)
		}

		request := controllers.DetailsRequest{
			CoinQuery: strconv.Itoa(int(coin)),
			Token:     token,
			Currency:  c.DefaultQuery("currency", watchmarket.DefaultCurrency),
		}
		response, err := controller.HandleDetailsRequest(request, ctx)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
