package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const userAgent = "kevin"

// Client manages communication with the Riot API.
type Client struct {
	// HTTP handles http client details. A nil value indicates a zero
	// [http.Client] is used.
	HTTP *http.Client

}

// DispatchRequest executes the request and decodes the json response in dst.
// It will return non-nil error in the case that an error, specified in the
// API, occurs.
func (c *Client) DispatchRequest(ctx context.Context, req *Request, dst any) error {
	if c.HTTP == nil {
		c.HTTP = &http.Client{}
	}

	httpReq, err := req.MakeHTTPRequest(ctx)

	res, err := c.HTTP.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}

	defer res.Body.Close()

	if err := GetError(res.StatusCode); err != nil {
		return fmt.Errorf("http response for host %s: %w", req.BaseURL, err)
	}

	if err := json.NewDecoder(res.Body).Decode(&dst); err != nil {
		return fmt.Errorf("json decoder: %w", err)
	}

	return nil
}

// Request represents a request to the Riot API.
type Request struct {
	// Method is the http method. A zero value defaults to GET.
	Method string

	// BaseURL is the host I guess... A zero value defaults to
	// https://na1.api.riotgames.com
	BaseURL string

	// Endpoint is an api endpoint.
	Endpoint string

	APIKey string

	Query url.Values

	Retry bool
}

// MakeHTTPRequest creates an http request.
func (r *Request) MakeHTTPRequest(ctx context.Context) (*http.Request, error) {
	// NOTE: currently never adding a body to http requests.
	u := r.BaseURL + r.Endpoint

	req, err := http.NewRequestWithContext(ctx, r.Method, u, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = r.Query.Encode()

	req.Header.Add("Accept", "application/json")

	req.Header.Add("User-Agent", userAgent)

	if r.APIKey != "" {
		req.Header.Set("X-Riot-Token", r.APIKey)
	}

	return req, err
}
