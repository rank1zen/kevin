// riot provides access to endpoints.
//
// NOTE: everything is routed through a region. So for matches, the internals
// access the continent but only returns per region.
package riot

import (
	"net/http"

	"github.com/rank1zen/kevin/internal/riot/internal"
)

// Client manages communication with the Riot API.
type Client struct {
	internals *internal.Client

	apiKey string

	// baseURL will override the platform specific hosts, if provided.
	baseURL string

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

// NewClient creates a new client. Will panic if there is an error in creation.
//
// NOTE: don't know about the concurrency stuff.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := Client{
		apiKey:    apiKey,
		internals: &internal.Client{},
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
		client.internals.HTTP = c
		return nil
	}
}

func WithBaseURL(u string) ClientOption {
	return func(client *Client) error {
		client.baseURL = u
		return nil
	}
}

var (
	ErrBadRequest           = internal.ErrBadRequest
	ErrUnauthorized         = internal.ErrUnauthorized
	ErrForbidden            = internal.ErrForbidden
	ErrNotFound             = internal.ErrNotFound
	ErrMethodNotAllowed     = internal.ErrMethodNotAllowed
	ErrUnsupportedMediaType = internal.ErrUnsupportedMediaType
	ErrRateLimitExceeded    = internal.ErrRateLimitExceeded
	ErrInternalServerError  = internal.ErrInternalServerError
	ErrBadGateway           = internal.ErrBadGateway
	ErrServiceUnavailable   = internal.ErrServiceUnavailable
	ErrGatewayTimeout       = internal.ErrGatewayTimeout
)
