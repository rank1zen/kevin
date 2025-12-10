// server contains all http server things.
package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/frontend/page"
	"github.com/rank1zen/kevin/internal/frontend/view/profile"
	"github.com/rank1zen/kevin/internal/frontend/view/search/searchmenu"
	"github.com/rank1zen/kevin/internal/service"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

func New(handler *service.Service, port int, opts ...ServerOption) *Server {
	srvr := Server{
		Logger:  slog.Default(),
		Address: fmt.Sprintf(":%d", port),
	}

	for _, opt := range opts {
		opt(&srvr)
	}

	router := http.NewServeMux()

	handleFunc := func(pattern string, handler http.Handler) {
		router.Handle(pattern, otelhttp.WithRouteTag(pattern, handler))
	}

	handleFunc("GET /{$}", (*page.HomePageHandler)(handler))

	handleFunc("GET /profile/{riotID}/{$}", (*page.ProfilePageHandler)(handler))
	handleFunc("GET /profile/{riotID}/live/{$}", (*page.ProfileLiveMatchPageHandler)(handler))

	handleFunc("GET /livematch/{matchID}/{$}", (*page.LivematchPageHandler)(handler))

	handleFunc("POST /partial/profile.HistoryEntry", (*profile.HistoryEntryHandler)(handler))
	handleFunc("POST /partial/profile.ChampionList", (*profile.ChampionListHandler)(handler))
	handleFunc("POST /partial/profile.MatchDetailBox", (*profile.MatchDetailBoxHandler)(handler))
	handleFunc("POST /partial/profile.UpdateProfile", (*profile.UpdateProfileHandler)(handler))
	handleFunc("POST /partial/searchmenu.SearchMenu", (*searchmenu.SearchMenuHandler)(handler))

	// Wrap with logging middleware
	loggedRouter := srvr.addLoggingMiddleware(router)

	// Wrap with OpenTelemetry HTTP instrumentation
	// otelhttp.WithRouteTag provides a cleaner span name than the full URL path.
	otelRouter := otelhttp.NewHandler(loggedRouter, "HTTP Server")

	main := http.NewServeMux()
	main.Handle("/", otelRouter) // Use the OpenTelemetry wrapped handler
	main.Handle("GET /static/", http.FileServer(http.FS(frontend.StaticAssets)))
	main.Handle("GET /ready", (*ReadyHandler)(handler))

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
