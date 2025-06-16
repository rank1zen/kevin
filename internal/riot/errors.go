package riot

import "errors"

var (
	ErrBadRequest           = errors.New("riot: bad request")
	ErrUnauthorized         = errors.New("riot: unauthorized")
	ErrForbidden            = errors.New("riot: forbidden")
	ErrNotFound             = errors.New("riot: not found")
	ErrMethodNotAllowed     = errors.New("riot: method not allowed")
	ErrUnsupportedMediaType = errors.New("riot: unsupported media type")
	ErrRateLimitExceeded    = errors.New("riot: rate limit exceeded")
	ErrInternalServerError  = errors.New("riot: internal server error")
	ErrBadGateway           = errors.New("riot: bad gateway")
	ErrServiceUnavailable   = errors.New("riot: service unavailable")
	ErrGatewayTimeout       = errors.New("riot: gateway timeout")
)
