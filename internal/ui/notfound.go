package ui

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/rank1zen/yujin/internal/logging"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := logging.FromContext(ctx)

	logger.Warn(r.URL.String())
	templ.Handler(NotFound(), templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, r)
}
