package endpoint

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/api/model"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"go.elastic.co/apm"
	"net/http"
)

// @Summary Get ticker values for a specific market
// @Id get_tickers
// @Description Get the ticker values from many market and coin/token
// @Accept json
// @Produce json
// @Tags Market
// @Param tickers body api.TickerRequest true "Ticker"
// @Success 200 {object} watchmarket.Tickers
// @Router /v1/market/ticker [post]
func GetTickersHandler(controller controllers.Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("POST /v1/market/ticker", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		request := controllers.TickerRequest{Currency: watchmarket.DefaultCurrency}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
			return
		}
		response, err := controller.HandleTickersRequest(request, ctx)
		if err != nil {
			handleError(c, err)
			return
		}
		if len(response.Tickers) == 0 {
			handleTickersError(c, request)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetTickerHandlerV2(controller controllers.Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v2/market/ticker/:id", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		currency := c.DefaultQuery("currency", watchmarket.DefaultCurrency)

		request := controllers.TickerRequestV2{Currency: currency, Ids: []string{c.Param("id")}}
		response, err := controller.HandleTickersRequestV2(request, ctx)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetTickersHandlerV2(controller controllers.Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("POST /v2/market/tickers", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		request := controllers.TickerRequestV2{Currency: watchmarket.DefaultCurrency}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
			return
		}
		response, err := controller.HandleTickersRequestV2(request, ctx)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func handleTickersError(c *gin.Context, req controllers.TickerRequest) {
	if len(req.Assets) == 0 || req.Assets == nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
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
