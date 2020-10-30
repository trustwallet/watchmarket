package endpoint

import (
	"errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"net/http"
)

type (
	ErrorResponse struct {
		Error ErrorDetails `json:"error"`
	}
	ErrorDetails struct {
		Message string `json:"message"`
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
		return http.StatusInternalServerError, errorResponse(errors.New("Internal Fail"))
	case watchmarket.ErrBadRequest:
		return http.StatusBadRequest, errorResponse(errors.New("Invalid request payload"))
	case watchmarket.ErrNotFound:
		return http.StatusNotFound, errorResponse(errors.New("Not found"))
	default:
		return http.StatusBadRequest, errorResponse(errors.New("Invalid request payload"))
	}
}
