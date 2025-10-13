// server contains all http server things.
package server

import (
	"log/slog"
	"net/http"

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

func New(handler *frontend.Handler, opts ...ServerOption) *Server {
	frontend := Server{}

	for _, opt := range opts {
		opt(&frontend)
	}

	router := http.NewServeMux()

	router.Handle("GET /{$}", (*page.HomePageHandler)(handler))
	router.Handle("GET /{riotID}/{$}", (*page.ProfilePageHandler)(handler))
	router.Handle("GET /summoner/fetch/{$}", nil)

	router.Handle("POST /summoner/matchlist", (*profile.HistoryEntryHandler)(handler))
	router.Handle("POST /summoner/live", nil)
	router.Handle("POST /summoner/champions", (*profile.ChampionListHandler)(handler))
	router.Handle("POST /match", (*profile.MatchDetailBoxHandler)(handler))

	frontend.router = router

	return &frontend
}

func (s *Server) Open() error {
	// err := http.ListenAndServe(s.Address, s)
	return nil
}
