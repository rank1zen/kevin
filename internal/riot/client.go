package riot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Client manages communication with the Riot API.
type Client struct {
	apiKey string

	httpClient *http.Client

	common service

	Account   *AccountService
	League    *LeagueService
	Match     *MatchService
	Spectator *SpectatorService
	Summoner  *SummonerService
}

type service struct {
	client *Client
}

func NewClient(apiKey string, opts ...ClientOption) *Client {

	c := Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}

	c.common.client = &c
	c.Account = (*AccountService)(&c.common)
	c.League = (*LeagueService)(&c.common)
	c.Match = (*MatchService)(&c.common)
	c.Spectator = (*SpectatorService)(&c.common)
	c.Summoner = (*SummonerService)(&c.common)

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			panic("setting up riot client")
		}
	}

	return &c
}

type ClientOption func(*Client) error

func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) error {
		client.httpClient = c
		return nil
	}
}

// makeAndDispatchRequest basically does everything
func (c *Client) makeAndDispatchRequest(ctx context.Context, region Region, endpoint string, dst any) error {
	u, err := url.JoinPath(regionToHost[region], endpoint)
	if err != nil {
		panic(err)
	}

	return c.makeAndDispatch(ctx, u, dst)
}

func (c *Client) makeAndDispatchRequestOnContinent(ctx context.Context, continent Continent, endpoint string, dst any) error {
	u, err := url.JoinPath(continentToHost[continent], endpoint)
	if err != nil {
		panic(err)
	}

	return c.makeAndDispatch(ctx, u, dst)
}

func (c *Client) makeAndDispatch(ctx context.Context, url string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Riot-Token", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http client: %w", err)
	}

	defer res.Body.Close()

	if err := getError(res.StatusCode); err != nil {
		return err
	}

	if err := json.NewDecoder(res.Body).Decode(&dst); err != nil {
		return fmt.Errorf("json decoder: %w", err)
	}

	return nil
}

// should have an Error type

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
	// ErrUnknown              = errors.New("riot: unknown")
)

func getError(status int) error {
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
		// assume status is OK!
		return nil
	}
}
