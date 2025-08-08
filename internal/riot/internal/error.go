package internal

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message    string
	StatusCode int
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("riot: %s", e.Message)
}

var (
	ErrBadRequest = ErrorResponse{
		Message:    "bad request",
		StatusCode: http.StatusBadRequest,
	}

	ErrUnauthorized = ErrorResponse{
		Message:    "unauthorized",
		StatusCode: http.StatusUnauthorized,
	}

	ErrForbidden = ErrorResponse{
		Message:    "forbidden",
		StatusCode: http.StatusForbidden,
	}

	ErrNotFound = ErrorResponse{
		Message:    "not found",
		StatusCode: http.StatusNotFound,
	}

	ErrMethodNotAllowed = ErrorResponse{
		Message:    "method not allowed",
		StatusCode: http.StatusMethodNotAllowed,
	}

	ErrUnsupportedMediaType = ErrorResponse{
		Message:    "unsupported media type",
		StatusCode: http.StatusUnsupportedMediaType,
	}

	ErrRateLimitExceeded = ErrorResponse{
		Message:    "rate limit exceeded",
		StatusCode: http.StatusTooManyRequests,
	}

	ErrInternalServerError = ErrorResponse{
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
	}

	ErrBadGateway = ErrorResponse{
		Message:    "bad gateway",
		StatusCode: http.StatusBadGateway,
	}

	ErrServiceUnavailable = ErrorResponse{
		Message:    "service unavailable",
		StatusCode: http.StatusServiceUnavailable,
	}

	ErrGatewayTimeout = ErrorResponse{
		Message:    "gateway timeout",
		StatusCode: http.StatusGatewayTimeout,
	}
)

// GetError interprets a response status code. Any code not specified
// by in Riot's docs is interpreted as a success.
func GetError(status int) error {
	switch status {
	case 400:
		return ErrBadRequest
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 405:
		return ErrMethodNotAllowed
	case 415:
		return ErrUnsupportedMediaType
	case 429:
		return ErrRateLimitExceeded
	case 500:
		return ErrInternalServerError
	case 502:
		return ErrBadGateway
	case 503:
		return ErrServiceUnavailable
	case 504:
		return ErrGatewayTimeout
	default:
		// NOTE: assume status is OK!
		return nil
	}
}
