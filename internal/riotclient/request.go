package riotclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	userAgent      = "server:yujin:v1.0"
	defaultTimeout = 80
)

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

type errorResponse struct {
	Message string `json:"message"`
}

func execute(ctx context.Context, req *http.Request) (io.ReadCloser, error) {
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Riot-Token", os.Getenv("RIOT_API_KEY"))

	// NOTE: currently using default http client
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, fmt.Errorf("http: %v", err)
		}
	}

	switch res.StatusCode {
	case http.StatusBadRequest:
		defer res.Body.Close()
		var e errorResponse
		err := json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return nil, fmt.Errorf("%w (%v)", ErrBadRequest, err)
		}

		return nil, fmt.Errorf("%w (%s)", ErrBadRequest, e.Message)
	case http.StatusUnauthorized:
		res.Body.Close()
		return nil, ErrUnauthorized
	case http.StatusForbidden:
		res.Body.Close()
		return nil, ErrForbidden
	case http.StatusNotFound:
		res.Body.Close()
		return nil, ErrNotFound
	case http.StatusMethodNotAllowed:
	case http.StatusUnsupportedMediaType:
	case http.StatusTooManyRequests:
	case http.StatusInternalServerError:
	case http.StatusBadGateway:
	case http.StatusServiceUnavailable:
	case http.StatusGatewayTimeout:
	default:
	}

	return res.Body, nil
}
