package endpoint

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"net/http"
)

type (
	ErrorResponse struct {
		Error ErrorDetails `json:"error"`
	}
	ErrorDetails struct {
		Message string    `json:"message"`
		Code    ErrorCode `json:"code,omitempty"`
	}

	ErrorCode int
)

func errorResponse(err error) ErrorResponse {
	var message string
	if err != nil {
		message = err.Error()
	}
	return ErrorResponse{Error: ErrorDetails{
		Message: message,
	}}
}

func createErrorResponseAndStatusCode(err error) (int, ErrorResponse) {
	switch err.Error() {
	case watchmarket.ErrInternal:
		return http.StatusInternalServerError, errorResponse(errors.E("Internal Fail"))
	case watchmarket.ErrBadRequest:
		return http.StatusBadRequest, errorResponse(errors.E("Invalid request payload"))
	case watchmarket.ErrNotFound:
		return http.StatusNotFound, errorResponse(errors.E("Not found"))
	default:
		return http.StatusBadRequest, errorResponse(errors.E("Invalid request payload"))
	}
}
