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

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
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
var staticAssets embed.FS

// New creates a [Server]. If handler is nil, the default [Handler] is used.
func New(handler *Handler, opts ...FrontendOption) *Server {
	frontend := Server{
		handler: handler,
	}

	for _, opt := range opts {
		opt(&frontend)
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /", frontend.getHomePage)

	router.HandleFunc("GET /{riotID}", frontend.getSumonerPage)

	router.HandleFunc("POST /search", frontend.serveSearchResults)

	router.HandleFunc("POST /summoner/fetch", frontend.updateSummoner)

	router.HandleFunc("POST /summoner/matchlist", frontend.serveMatchlist)

	router.HandleFunc("POST /summoner/live", frontend.serveLiveMatch)

	router.HandleFunc("POST /summoner/champions", frontend.serveChampions)

	router.HandleFunc("POST /match", frontend.serveMatchDetail)

	loggedRouter := frontend.addLoggingMiddleware(router)

	main := http.NewServeMux()
	main.Handle("/", loggedRouter)
	main.Handle("GET /static/", http.FileServer(http.FS(staticAssets)))

	frontend.router = main

	return &frontend
}

type FrontendOption func(*Server)

func WithLogger(logger *slog.Logger) FrontendOption {
	return func(f *Server) {
		f.Logger = logger
	}
}

func (f *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.router.ServeHTTP(w, r)
}

func (f *Server) getHomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := fromCtx(ctx)

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	region := r.FormValue("region")

	payload := slog.Group("payload", "region", region)

	riotRegion := convertStringToRiotRegion(region)

	component, err := f.handler.GetHomePage(ctx, riotRegion)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.ToTempl(ctx).Render(ctx, w)
}

func (f *Server) getSumonerPage(w http.ResponseWriter, r *http.Request) {
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

	component, err := f.handler.GetSummonerPage(ctx, riotRegion, name, tag)
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

	w.WriteHeader(http.StatusOK)
	component.ToTempl(ctx).Render(ctx, w)
}

func (f *Server) serveSearchResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	var (
		region = r.FormValue("region")
		q      = r.FormValue("q")
	)

	payload := slog.Group("payload", "region", region, "q", q)

	riotRegion := convertStringToRiotRegion(region)

	component, err := f.handler.GetSearchResults(ctx, riotRegion, q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", slog.Any("err", err), payload)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.ToTempl(ctx).Render(ctx, w)
}

func (f *Server) serveMatchlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	req, err := decode[MatchHistoryRequest](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("bad request", "err", err)
		return
	}

	payload := slog.Any("request", req)

	component, err := f.handler.GetMatchHistory(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.ToTempl(ctx).Render(ctx, w)
}

func (f *Server) updateSummoner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	decoded, _ := decode[UpdateSummonerRequest](r)

	payload := slog.Group("payload", "region", decoded.Region, "name", decoded.Name, "tag", decoded.Tag)

	if err := f.handler.UpdateSummoner(ctx, decoded); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("service failed", "err", err, payload)
		return
	}

	// Redirect to summoner page
	w.Header().Set("HX-Location", fmt.Sprintf("/%s-%s", decoded.Name, decoded.Tag))
	w.WriteHeader(http.StatusOK)
}

func (f *Server) serveChampions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	decoded, err := decode[GetSummonerChampionsRequest](r)

	payload := slog.Group("payload", "region", decoded.Region, "puuid", decoded.PUUID, "week", decoded.Week)

	component, err := f.handler.GetSummonerChampions(ctx, decoded)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.ToTempl(ctx).Render(ctx, w)
}

func (f *Server) serveLiveMatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	req, err := decode[LiveMatchRequest](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("bad request", "err", err)
		return
	}

	payload := slog.Group("payload", "region", req.Region, "puuid", req.PUUID)

	component, err := f.handler.GetLiveMatch(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.ToTempl(ctx).Render(ctx, w)
}

func (f *Server) serveMatchDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decode[MatchDetailRequest](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("bad request", "err", err)
		return
	}

	payload := slog.Group("payload", "region", req.Region, "match_id", req.MatchID)

	component, err := f.handler.GetMatchDetail(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Debug("failed service", "err", err, payload)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.ToTempl(ctx).Render(ctx, w)
}

func (f *Server) addLoggingMiddleware(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now()

		requestLogger := f.Logger.With(slog.Group("request", "method", r.Method, "endpoint", r.URL))

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
