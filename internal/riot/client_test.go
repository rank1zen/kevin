package riot_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
)

// MakeTestClient returns server that will always return the body read from
// specified file.
//
// The mock server will check the X-Riot-Token header is set.
func MakeTestClient(t *testing.T, code int, filename string, validators ...func(*http.Request)) (*riot.Client, *httptest.Server) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	mockServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				for _, v := range validators {
					v(r)
				}

				assert.Equal(t, "test-api-key", r.Header.Get("X-Riot-Token"))

				w.WriteHeader(code)
				_, _ = io.Copy(w, f)

				r.Body.Close()
			},
		),
	)

	client := riot.NewClient("test-api-key", riot.WithHTTPClient(mockServer.Client()), riot.WithBaseURL(mockServer.URL))

	return client, mockServer
}
