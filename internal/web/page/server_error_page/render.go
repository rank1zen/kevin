package server_error_page

import (
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
)

func Render(w http.ResponseWriter, r *http.Request, err error, msg string) {
	frontend.LogError(r, fmt.Errorf("%s: %w", msg, err))
	w.WriteHeader(http.StatusInternalServerError)
	_ = Index(r.Context(), &IndexData{}).Render(r.Context(), w)
}
