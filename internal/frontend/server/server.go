// server contains all http server things.
package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/httplog/v3"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/frontend/page"
	"github.com/rank1zen/kevin/internal/frontend/view/profile"
	"github.com/rank1zen/kevin/internal/frontend/view/search/searchmenu"
	"github.com/rank1zen/kevin/internal/service"
)

type Server struct {
	router *http.ServeMux

	// Logger is used to log events emitted from Server. A nil value
	// indicates Server will use [slog.Default]
	Logger *slog.Logger

	// Address is the address the server will listen on.
	Address string
}

type ServerOption func(*Server)

func WithLogger(logger *slog.Logger) ServerOption {
	return func(f *Server) {
		f.Logger = logger
	}
}

func New(s *service.Service, port int, opts ...ServerOption) *Server {
	srvr := Server{
		Logger:  slog.Default(),
		Address: fmt.Sprintf(":%d", port),
	}

	for _, opt := range opts {
		opt(&srvr)
	}

	router := http.NewServeMux()

	router.Handle("GET /{$}", (*page.HomePageHandler)(s))

	router.Handle("GET /profile/{riotID}/{$}", (*page.ProfilePageHandler)(s))
	router.Handle("GET /profile/{riotID}/live/{$}", (*page.ProfileLiveMatchPageHandler)(s))

	router.Handle("GET /livematch/{matchID}/{$}", (*page.LivematchPageHandler)(s))

	router.Handle("POST /partial/profile.HistoryEntry", (*profile.HistoryEntryHandler)(s))
	router.Handle("POST /partial/profile.ChampionList", (*profile.ChampionListHandler)(s))
	router.Handle("POST /partial/profile.MatchDetailBox", (*profile.MatchDetailBoxHandler)(s))
	router.Handle("POST /partial/profile.UpdateProfile", (*profile.UpdateProfileHandler)(s))
	router.Handle("POST /partial/searchmenu.SearchMenu", (*searchmenu.SearchMenuHandler)(s))

	middlewares := []func(http.Handler) http.Handler{
		httplog.RequestLogger(srvr.Logger, &httplog.Options{
			Level:         0,
			Schema:        httplog.SchemaGCP,
			RecoverPanics: true,

			LogRequestHeaders:  []string{},
			LogResponseHeaders: []string{},

			LogRequestBody: func(req *http.Request) bool {
				return req.Header.Get("Debug") == "reveal-body-logs"
			},
			LogResponseBody: func(req *http.Request) bool {
				return req.Header.Get("Debug") == "reveal-body-logs"
			},
		}),
	}

	var h http.Handler = router
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	main := http.NewServeMux()
	main.Handle("/", h) // Use the OpenTelemetry wrapped handler
	main.Handle("GET /static/", http.FileServer(http.FS(frontend.StaticAssets)))
	main.Handle("GET /ready", (*ReadyHandler)(s))

	srvr.router = main

	return &srvr
}

func (f *Server) addLoggingMiddleware(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now()

		requestLogger := f.Logger.With(slog.Group("request", "method", r.Method, "endpoint", r.URL.Path))

		// Set the logger in the context for downstream handlers
		r = r.WithContext(frontend.LoggerNewContext(r.Context(), requestLogger))

		handler.ServeHTTP(w, r)

		requestLogger.Info("request completed", "duration", time.Since(ts))
	}

	return http.HandlerFunc(fn)
}

func (s *Server) Open() error {
	return http.ListenAndServe(s.Address, s.router)
}
