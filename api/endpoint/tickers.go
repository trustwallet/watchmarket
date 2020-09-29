package endpoint

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"go.elastic.co/apm"
	"net/http"
	"strings"
)

// @Summary Get ticker values for a specific market
// @Id get_tickers
// @Description Get the ticker values from many market and coin/token
// @Accept json
// @Produce json
// @Tags Tickers
// @Param tickers body controllers.TickerRequest true "Ticker"
// @Success 200 {object} controllers.TickerResponse
// @Router /v1/market/ticker [post]
func GetTickersHandler(controller controllers.TickersController) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("POST /v1/market/ticker", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		request := controllers.TickerRequest{Currency: watchmarket.DefaultCurrency}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(errors.E("Invalid request payload")))
			return
		}
		response, err := controller.HandleTickersRequest(request, ctx)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}
		if len(response.Tickers) == 0 {
			handleTickersError(c, request)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Get ticker for a specific market
// @Id get_ticker
// @Description Get the ticker for specific id
// @Accept json
// @Produce json
// @Tags Tickers
// @Param id path string true "id" default(c714_tXRP-BF2)
// @Param currency query string false "The currency to show coin assets in" default(USD)
// @Success 200 {object} controllers.TickerResponseV2
// @Router /v2/market/ticker/{id} [get]
func GetTickerHandlerV2(controller controllers.TickersController) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v2/market/ticker/:id", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		currency := c.DefaultQuery("currency", watchmarket.DefaultCurrency)
		request := controllers.TickerRequestV2{Currency: currency, Ids: []string{c.Param("id")}}
		response, err := controller.HandleTickersRequestV2(request, ctx)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

// @Summary Get tickers for list of ids
// @Id post_tickers_v2
// @Description Get the tickers for list of ids
// @Accept json
// @Produce json
// @Tags Tickers
// @Param tickers body controllers.TickerRequestV2 true "Ticker"
// @Success 200 {object} controllers.TickerResponseV2
// @Router /v2/market/tickers [post]
func PostTickersHandlerV2(controller controllers.TickersController) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("POST /v2/market/tickers", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		request := controllers.TickerRequestV2{Currency: watchmarket.DefaultCurrency}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(errors.E("Invalid request payload")))
			return
		}
		request.Ids = removeDuplicates(request.Ids)
		response, err := controller.HandleTickersRequestV2(request, ctx)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Get tickers for list of ids
// @Id get_tickers_v2
// @Description Get the tickers for list of ids
// @Accept json
// @Produce json
// @Tags Tickers
// @Param assets path string true "List of asset ids"
// @Param currency query string empty "Currency symbol"
// @Success 200 {object} controllers.TickerResponseV2
// @Router /v2/market/tickers/{assets} [get]
func GetTickersHandlerV2(controller controllers.TickersController) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v2/market/tickers", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		currency := c.DefaultQuery("currency", watchmarket.DefaultCurrency)
		assets := c.Param("assets")
		if len(assets) == 0 {
			c.JSON(http.StatusBadRequest, errorResponse(errors.E("Invalid request payload")))
			return
		}
		assetsIds := removeDuplicates(strings.Split(assets, ","))
		request := controllers.TickerRequestV2{Currency: currency, Ids: assetsIds}
		response, err := controller.HandleTickersRequestV2(request, ctx)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func handleTickersError(c *gin.Context, req controllers.TickerRequest) {
	if len(req.Assets) == 0 || req.Assets == nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.E("Invalid request payload")))
		return
	}
	emptyResponse := controllers.TickerResponse{
		Currency: req.Currency,
	}
	tickers := make(watchmarket.Tickers, 0, len(req.Assets))
	for _, t := range req.Assets {
		tickers = append(tickers, watchmarket.Ticker{
			Coin:     t.Coin,
			TokenId:  t.TokenId,
			CoinType: t.CoinType,
		})
	}
	emptyResponse.Tickers = tickers
	c.JSON(http.StatusOK, emptyResponse)
}

func removeDuplicates(values []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range values {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
