package html

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/rank1zen/yujin/internal/logging"
)

func OK(w http.ResponseWriter, r *http.Request, c templ.Component) {
	ctx := r.Context()
	w.Header().Set("ContentType", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	c.Render(ctx, w)
}

func ServerError(w http.ResponseWriter, r *http.Request, c templ.Component, err error) {
	ctx := r.Context()
	w.Header().Set("ContentType", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate, no-store")
	w.WriteHeader(http.StatusInternalServerError)
	logging.FromContext(ctx).Sugar().Debug(err)
	c.Render(ctx, w)
}

func BadRequest(w http.ResponseWriter, r *http.Request, c templ.Component, err error) {
	ctx := r.Context()
	w.Header().Set("ContentType", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate, no-store")
	w.WriteHeader(http.StatusBadRequest)
	c.Render(ctx, w)
}

func NotFound() {

}
