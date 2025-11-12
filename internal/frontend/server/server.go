// server contains all http server things.
package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/frontend/page"
	"github.com/rank1zen/kevin/internal/frontend/view/profile"
)

type Server struct {
	router *http.ServeMux

	// Logger is used to log events emitted from Server. A nil value
	// indicates Server will use [slog.Default]
	Logger *slog.Logger
}

type ServerOption func(*Server)

func WithLogger(logger *slog.Logger) ServerOption {
	return func(f *Server) {
		f.Logger = logger
	}
}

func New(handler *internal.Datasource, opts ...ServerOption) *Server {
	srvr := Server{}

	for _, opt := range opts {
		opt(&srvr)
	}

	router := http.NewServeMux()

	router.Handle("GET /{$}", (*page.HomePageHandler)(handler))
	router.Handle("GET /profile/{riotID}/{$}", (*page.ProfilePageHandler)(handler))
	router.Handle("GET /profile/{riotID}/live/{$}", (*page.ProfileLiveMatchPageHandler)(handler))

	router.Handle("POST /partial/profile.HistoryEntry", (*profile.HistoryEntryHandler)(handler))
	router.Handle("POST /partial/profile.ChampionList", (*profile.ChampionListHandler)(handler))
	router.Handle("POST /partial/profile.MatchDetailBox", (*profile.MatchDetailBoxHandler)(handler))
	router.Handle("POST /partial/profile.UpdateProfile", (*profile.UpdateProfileHandler)(handler))

	loggedRouter := srvr.addLoggingMiddleware(router)

	main := http.NewServeMux()
	main.Handle("/", loggedRouter)
	main.Handle("GET /static/", http.FileServer(http.FS(frontend.StaticAssets)))

	srvr.router = main

	return &srvr
}

func (f *Server) addLoggingMiddleware(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now()

		requestLogger := f.Logger.With(slog.Group("request", "method", r.Method, "endpoint", r.URL))

		r = r.WithContext(frontend.LoggerNewContext(r.Context(), requestLogger))

		handler.ServeHTTP(w, r)

		requestLogger.Info(fmt.Sprintf("request completed in %v", time.Since(ts)))
	}

	return http.HandlerFunc(fn)
}

func (s *Server) Open() error {
	return http.ListenAndServe("0.0.0.0:4001", s.router)
}
