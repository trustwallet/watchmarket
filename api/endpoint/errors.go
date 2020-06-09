package endpoint

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"net/http"
)

const (
	Default ErrorCode = iota
	InvalidQuery
	RequestedDataNotFound
	InternalFail
)

type (
	ErrorResponse struct {
		Error ErrorDetails `json:"error"`
	}
	ErrorDetails struct {
		Message string    `json:"message"`
		Code    ErrorCode `json:"code"`
	}

	ErrorCode int
)

func createErrorResponse(code ErrorCode, err error) ErrorResponse {
	var message string
	if err != nil {
		message = err.Error()
	}
	return ErrorResponse{Error: ErrorDetails{
		Message: message,
		Code:    code,
	}}
}

func createError(err error) (int, ErrorResponse) {
	switch err.Error() {
	case watchmarket.ErrInternal:
		return http.StatusInternalServerError, createErrorResponse(InternalFail, errors.E("Internal Fail"))
	case watchmarket.ErrBadRequest:
		return http.StatusBadRequest, createErrorResponse(InvalidQuery, errors.E("Invalid request payload"))
	case watchmarket.ErrNotFound:
		return http.StatusNotFound, createErrorResponse(RequestedDataNotFound, errors.E("Not found"))
	default:
		return http.StatusInternalServerError, createErrorResponse(Default, errors.E("Invalid request payload"))
	}
}
