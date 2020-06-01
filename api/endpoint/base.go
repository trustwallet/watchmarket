package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/api/model"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"net/http"
)

func handleError(c *gin.Context, err error) {
	switch err.Error() {
	case controllers.ErrInternal:
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(model.InternalFail, errors.E("Internal Fail")))
		return
	case controllers.ErrBadRequest:
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
		return
	case controllers.ErrNotFound:
		c.JSON(http.StatusNotFound, model.CreateErrorResponse(model.RequestedDataNotFound, errors.E("Not found")))
		return
	default:
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
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

func ParseID(id string) (uint, string, watchmarket.CoinType, error) {
	return 0, "", watchmarket.Coin, nil
}
