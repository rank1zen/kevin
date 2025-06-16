package riot

import (
	"net/http"
)

// Client abstracts the logic pertaining to riot types and errors.
type Client struct {
	httpClient *http.Client
	apiKey     string
}

type ClientOption func(*Client) error

func NewClient(opts ...ClientOption) *Client {
	client := Client{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(&client)
	}

	return &client
}

func WithHttpClient(c *http.Client) ClientOption {
	return func(client *Client) error {
		client.httpClient = c
		return nil
	}
}

func WithApiKey(apiKey string) ClientOption {
	return func(client *Client) error {
		client.apiKey = apiKey
		return nil
	}
}
