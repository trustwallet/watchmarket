package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/golibs/asset"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/services/controllers"
	"net/http"
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
		coinId, err := controllers.GetCoinId(c.Query("coin"))
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}
		request := controllers.ChartRequest{
			Asset: controllers.Asset{
				CoinId:  coinId,
				TokenId: c.Query("token"),
			},
			Currency:  controllers.GetCurrency(c.Query("currency")),
			TimeStart: controllers.GetTimeStart(c.Query("time_start")),
			MaxItems:  controllers.GetMaxItems(c.Query("max_items")),
		}

		response, err := controller.HandleChartsRequest(request)
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
		coinId, tokenId, err := asset.ParseID(c.Param("id"))
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}
		if _, ok := coin.Coins[coinId]; !ok {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}
		request := controllers.ChartRequest{
			Asset: controllers.Asset{
				CoinId:  coinId,
				TokenId: tokenId,
			},
			Currency:  controllers.GetCurrency(c.Query("currency")),
			TimeStart: controllers.GetTimeStart(c.Query("time_start")),
			MaxItems:  controllers.GetMaxItems(c.Query("max_items")),
		}

		response, err := controller.HandleChartsRequest(request)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
