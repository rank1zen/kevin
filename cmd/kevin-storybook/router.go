package main

import (
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", &IndexHandler{})
	mux.Handle("GET /static/", http.FileServer(http.FS(frontend.StaticAssets)))

	for _, route := range routes {
		mux.HandleFunc(
			fmt.Sprintf("GET /%s/{$}", route.Name),
			func(w http.ResponseWriter, r *http.Request) {
				c := Fragment(route.Renderer)
				if err := c.Render(r.Context(), w); err != nil {
					w.WriteHeader(500)
				}
			},
		)
	}

	return mux
}
