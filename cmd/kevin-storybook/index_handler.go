package main

import (
	"net/http"
)

type IndexHandler struct{}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := Index(r.Context(), IndexData{Routes: routes})
	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(500)
	}
}
