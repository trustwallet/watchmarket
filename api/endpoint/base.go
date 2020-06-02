package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/api/model"
	"github.com/trustwallet/blockatlas/pkg/errors"
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
