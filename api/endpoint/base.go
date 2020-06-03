package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/api/model"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"net/http"
)

func handleError(c *gin.Context, err error) {
	switch err.Error() {
	case watchmarket.ErrInternal:
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(model.InternalFail, errors.E("Internal Fail")))
		return
	case watchmarket.ErrBadRequest:
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
		return
	case watchmarket.ErrNotFound:
		c.JSON(http.StatusNotFound, model.CreateErrorResponse(model.RequestedDataNotFound, errors.E("Not found")))
		return
	default:
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
	}
}
