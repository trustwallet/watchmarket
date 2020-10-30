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

// @Summary Get charts coin assets data for a specific coin
// @Id get_charts_coin_info
// @Description Get the charts coin assets data from an market and coin/contract
// @Accept json
// @Produce json
// @Tags Info
// @Param coin query string true "Coin id" default(60)
// @Param token query string false "Token id"
// @Param currency query string false "The currency to show coin assets in" default(USD)
// @Success 200 {object} watchmarket.CoinDetails
// @Router /v1/market/info [get]
func GetCoinInfoHandler(controller controllers.InfoController) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v1/market/info", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		request := controllers.DetailsRequest{
			CoinQuery: c.Query("coin"),
			Token:     c.Query("token"),
			Currency:  c.DefaultQuery("currency", watchmarket.DefaultCurrency),
		}
		response, err := controller.HandleInfoRequest(request, ctx)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Get charts coin assets data for a specific coin
// @Id get_charts_coin_info_v2
// @Description Get the charts coin assets data from an market and coin/contract
// @Accept json
// @Produce json
// @Tags Info
// @Param id path string true "id" default(c714)
// @Param currency query string false "The currency to show coin assets in" default(USD)
// @Success 200 {object} watchmarket.CoinDetails
// @Router /v2/market/info/{id} [get]
func GetCoinInfoHandlerV2(controller controllers.InfoController) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v2/market/info/:id", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		coin, token, err := asset.ParseID(c.Param("id"))
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		request := controllers.DetailsRequest{
			CoinQuery: strconv.Itoa(int(coin)),
			Token:     token,
			Currency:  c.DefaultQuery("currency", watchmarket.DefaultCurrency),
		}
		response, err := controller.HandleInfoRequest(request, ctx)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
