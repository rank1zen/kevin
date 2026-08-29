package runtime

import (
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal/profile"
)

func initializeServer(cfg Config) *http.Server {
	profileHandler := profile.Handler{}

	mux := http.NewServeMux()

	routeProfileService(mux, profileHandler)

	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)

	return &http.Server{
		Addr:      fmt.Sprintf(":%d", cfg.Port),
		Handler:   mux,
		Protocols: p,
	}
}
