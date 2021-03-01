package endpoint

import (
	"errors"
	"github.com/trustwallet/golibs/asset"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
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
		request := controllers.TickerRequest{Currency: watchmarket.DefaultCurrency}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid request payload")))
			return
		}
		tickers, err := controller.HandleTickersRequest(request)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}
		if len(tickers) == 0 {
			handleTickersError(c, request)
			return
		}

		c.JSON(http.StatusOK, controllers.TickerResponse{
			Currency: request.Currency,
			Tickers:  tickers,
		})
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
		request := controllers.TickerRequest{
			Currency: c.DefaultQuery("currency", watchmarket.DefaultCurrency),
			Assets:   parseAssetIds([]string{c.Param("id")}),
		}
		tickers, err := controller.HandleTickersRequest(request)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}
		c.JSON(http.StatusOK, mapToResponse(request.Currency, tickers))
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
		var request controllers.TickerRequestV2
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid request payload")))
			return
		}
		request.Ids = removeDuplicates(request.Ids)
		currency := request.Currency
		if len(currency) == 0 {
			currency = watchmarket.DefaultCurrency
		}
		tickers, err := controller.HandleTickersRequest(controllers.TickerRequest{
			Currency: currency,
			Assets:   parseAssetIds(request.Ids),
		})
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		c.JSON(http.StatusOK, mapToResponse(currency, tickers))
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
		assets := c.Param("assets")
		if len(assets) == 0 {
			c.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid request payload")))
			return
		}
		assetsIds := removeDuplicates(strings.Split(assets, ","))
		request := controllers.TickerRequest{
			Currency: c.DefaultQuery("currency", watchmarket.DefaultCurrency),
			Assets:   parseAssetIds(assetsIds),
		}
		tickers, err := controller.HandleTickersRequest(request)
		if err != nil {
			code, response := createErrorResponseAndStatusCode(err)
			c.AbortWithStatusJSON(code, response)
			return
		}

		c.JSON(http.StatusOK, mapToResponse(request.Currency, tickers))
	}
}

func mapToResponse(currency string, tickers watchmarket.Tickers) controllers.TickerResponseV2 {
	response := controllers.TickerResponseV2{
		Currency: currency,
	}
	response.Tickers = make([]controllers.TickerPrice, 0, len(tickers))
	for _, ticker := range tickers {
		response.Tickers = append(response.Tickers, controllers.TickerPrice{
			Change24h: ticker.Price.Change24h,
			Provider:  ticker.Price.Provider,
			Price:     ticker.Price.Value,
			ID:        asset.BuildID(ticker.Coin, ticker.TokenId),
		})
	}
	return response
}

func handleTickersError(c *gin.Context, req controllers.TickerRequest) {
	if len(req.Assets) == 0 || req.Assets == nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid request payload")))
		return
	}
	emptyResponse := controllers.TickerResponse{
		Currency: req.Currency,
	}
	tickers := make(watchmarket.Tickers, 0, len(req.Assets))
	for _, t := range req.Assets {
		tickers = append(tickers, watchmarket.Ticker{
			Coin:     t.CoinId,
			TokenId:  t.TokenId,
			CoinType: t.CoinType,
		})
	}
	emptyResponse.Tickers = tickers
	c.JSON(http.StatusOK, emptyResponse)
}

func removeDuplicates(values []string) (result []string) {
	keys := make(map[string]bool)
	for _, entry := range values {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}

func parseAssetIds(ids []string) (assets []controllers.Asset) {
	for _, id := range ids {
		if coinId, tokenId, err := asset.ParseID(id); err == nil {
			assets = append(assets, controllers.Asset{
				CoinId:  coinId,
				TokenId: tokenId,
			})
		}
	}
	return assets
}
