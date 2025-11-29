package server

import (
	"net/http"

	"github.com/rank1zen/kevin/internal/service"
)

type ReadyHandler service.Service

func (h *ReadyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := (*service.Service)(h).CheckHealth(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}
