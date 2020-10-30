package endpoint

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/golibs/asset"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"go.elastic.co/apm"
	"net/http"
	"strconv"
)

// @Summary Get charts data for a specific coin
// @Id get_charts_data
// @Description Get the charts data from an market and coin/token
// @Accept json
// @Produce json
// @Tags Charts
// @Param coin query int true "Coin id" default(60)
// @Param token query string false "Token id"
// @Param time_start query int false "Start timestamp" default(1574483028)
// @Param max_items query int false "Max number of items in result prices array" default(64)
// @Param currency query string false "The currency to show charts" default(USD)
// @Success 200 {object} watchmarket.Chart
// @Router /v1/market/charts [get]
func GetChartsHandler(controller controllers.ChartsController) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v1/market/charts", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		request := controllers.ChartRequest{
			CoinQuery:    c.Query("coin"),
			Token:        c.Query("token"),
			Currency:     c.DefaultQuery("currency", watchmarket.DefaultCurrency),
			TimeStartRaw: c.Query("time_start"),
			MaxItems:     c.Query("max_items"),
		}

		response, err := controller.HandleChartsRequest(request, ctx)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Get charts data for a specific id
// @Id get_charts_data_v2
// @Description Get the charts data from an market and coin/token
// @Accept json
// @Produce json
// @Tags Charts
// @Param id path string true "id" default(c60)
// @Param time_start query int false "Start timestamp" default(1574483028)
// @Param max_items query int false "Max number of items in result prices array" default(64)
// @Param currency query string false "The currency to show charts" default(USD)
// @Success 200 {object} watchmarket.Chart
// @Router /v2/market/charts/{id} [get]
func GetChartsHandlerV2(controller controllers.ChartsController) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v2/market/charts/:id", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		coin, token, err := asset.ParseID(c.Param("id"))
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		request := controllers.ChartRequest{
			CoinQuery:    strconv.Itoa(int(coin)),
			Token:        token,
			Currency:     c.DefaultQuery("currency", watchmarket.DefaultCurrency),
			TimeStartRaw: c.Query("time_start"),
			MaxItems:     c.Query("max_items"),
		}

		response, err := controller.HandleChartsRequest(request, ctx)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
