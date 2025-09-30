package frontend

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/page"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/view/profile"
)

var (
	ErrInvalidRegion = errors.New("invalid region")
	ErrInvalidRiotID = errors.New("invalid riot id")
	ErrInvalidPUUID  = errors.New("invalid puuid")
)

// Server serves [templ.Component].
type Server struct {
	router *http.ServeMux

	// Logger is used to log events emitted from Server. A nil value
	// indicates Server will use [slog.Default]
	Logger *slog.Logger

	handler *Handler
}

//go:embed static
var StaticAssets embed.FS

// New creates a [Server]. If handler is nil, the default [Handler] is used.
func New(handler *Handler, opts ...FrontendOption) *Server {
	frontend := Server{
		handler: handler,
	}

	for _, opt := range opts {
		opt(&frontend)
	}

	router := http.NewServeMux()

	var (
		profile = ProfileService{}
		search  = SearchService{}
	)

	profile.RegisterRoutes(router)

	router.HandleFunc("GET /", frontend.getHomePage)

	router.HandleFunc("GET /{riotID}", frontend.getSumonerPage)

	router.HandleFunc("POST /search", frontend.serveSearchResults)

	loggedRouter := frontend.addLoggingMiddleware(router)

	main := http.NewServeMux()
	main.Handle("/", loggedRouter)
	main.Handle("GET /static/", http.FileServer(http.FS(StaticAssets)))

	frontend.router = main

	return &frontend
}

func (s *Server) Open() error {
	err := http.ListenAndServe(s.Address, s)
	return err
}

type FrontendOption func(*Server)

func WithLogger(logger *slog.Logger) FrontendOption {
	return func(f *Server) {
		f.Logger = logger
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) getHomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := fromCtx(ctx)

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	region := r.FormValue("region")

	payload := slog.Group("payload", "region", region)

	riotRegion := convertStringToRiotRegion(region)

	v, err := s.handler.GetHomePage(ctx, riotRegion)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	if err := page.HomePage(ctx, *v).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed rendering", "err", err)
		return
	}
}

func (s *Server) getSumonerPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := fromCtx(ctx)

	region := r.FormValue("region")

	payload := slog.Group("payload", "region", region)

	riotID := r.PathValue("riotID")
	name, tag, err := ParseRiotID(riotID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("failed to resolve riot id", "err", err, payload)
		return
	}

	riotRegion := convertStringToRiotRegion(region)

	data, err := s.handler.GetSummonerPage(ctx, riotRegion, name, tag, time.UTC)
	if err != nil {
		if errors.Is(err, internal.ErrSummonerNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Info("summoner is not found", "name", name, "tag", tag, payload)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	championListCh := make(chan profile.ChampionListData)
	data.ChampionListCh = championListCh

	rankCardCh := make(chan profile.RankCardData)
	data.RankCardCh = rankCardCh

	historyEntryCh := make(chan profile.HistoryEntryData)
	data.HistoryEntryCh = historyEntryCh
	go func() {
		defer close(historyEntryCh)
		defer close(championListCh)
		defer close(rankCardCh)

		profileHandler := ProfileHandler{Datasource: s.handler.Datasource}
		days := GetDays(time.Now())
		for i := range len(days) - 1 {
			historyEntryData, err := profileHandler.GetMatchHistory(ctx, MatchHistoryRequest{
				Region:  riotRegion,
				PUUID:   data.PUUID,
				StartTS: days[i+1],
				EndTS:   days[i],
			})
			if err == nil {
				historyEntryCh <- *historyEntryData
			}
		}
	}()

	component := page.ProfilePage(ctx, *data)

	templ.Handler(component, templ.WithStreaming()).ServeHTTP(w, r)
	// if err := component.Render(ctx, w); err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	logger.Debug("failed rendering", "err", err)
	// 	return
	// }
}

func (s *Server) serveSearchResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	var (
		region = r.FormValue("region")
		q      = r.FormValue("q")
	)

	payload := slog.Group("payload", "region", region, "q", q)

	riotRegion := convertStringToRiotRegion(region)

	c, err := s.handler.GetSearchResults(ctx, riotRegion, q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", slog.Any("err", err), payload)
		return
	}

	if err := c.Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed rendering", "err", err)
		return
	}
}

func (s *Server) addLoggingMiddleware(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now()

		requestLogger := s.Logger.With(slog.Group("request", "method", r.Method, "endpoint", r.URL))

		r = r.WithContext(newCtx(r.Context(), requestLogger))

		handler.ServeHTTP(w, r)

		requestLogger.Info(fmt.Sprintf("request completed in %v", time.Since(ts)))
	}

	return http.HandlerFunc(fn)
}

type ctxKey struct{}

func newCtx(parent context.Context, logger *slog.Logger) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if lp, ok := parent.Value(ctxKey{}).(*slog.Logger); ok {
		// if parent already has the same loggger
		if lp == logger {
			return parent
		}
	}

	return context.WithValue(parent, ctxKey{}, logger)
}

func fromCtx(parent context.Context) *slog.Logger {
	if logger, ok := parent.Value(ctxKey{}).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}

// ParseRiotID parses a name-tag serperated by a '-' character. It returns
// [ErrInvalidRiotID] if riotID is not exactly "name-tag".
func ParseRiotID(riotID string) (name, tag string, err error) {
	index := strings.Index(riotID, "-")
	if index == -1 {
		return "", "", ErrInvalidRiotID
	}

	if index == len(riotID)-1 {
		return "", "", ErrInvalidRiotID
	}

	name = riotID[:index]
	tag = riotID[index+1:]

	if index := strings.Index(tag, "-"); index != -1 {
		return "", "", ErrInvalidRiotID
	}

	return name, tag, nil
}

// ParsePUUID parses a string puuid to [riot.PUUID], and returns
// [ErrInvalidPUUID] if puuid is not valid.
func ParsePUUID(puuid string) (riot.PUUID, error) {
	if len(puuid) != 78 {
		return "", ErrInvalidPUUID
	}

	// lol
	return riot.PUUID(puuid), nil
}

type Validator interface {
	Validate() (problems map[string]string)
}

func decode[T Validator](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}

	if problems := v.Validate(); len(problems) != 0 {
		return v, errors.New("validation error")
	}

	return v, nil
}
