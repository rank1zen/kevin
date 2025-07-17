package frontend

import (
	"context"
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
	ErrInvalidPUUID = errors.New("invalid puuid")
)

// Frontend serves [templ.Component].
type Frontend struct {
	router *http.ServeMux

	// logs emitted from [Frontend] will use Logger.
	Logger *slog.Logger

	handler *Handler
}

// New creates a [Frontend]. If handler is nil, the default [Handler] is used.
func New(handler *Handler, opts ...FrontendOption) *Frontend {
	frontend := Frontend{
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

	loggedRouter := frontend.addLoggingMiddleware(router)

	main := http.NewServeMux()
	main.Handle("/", loggedRouter)
	main.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	frontend.router = main

	return &frontend
}

type FrontendOption func(*Frontend)

func WithLogger(logger *slog.Logger) FrontendOption {
	return func(f *Frontend) {
		f.Logger = logger
	}
}

func (f *Frontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.router.ServeHTTP(w, r)
}

func (f *Frontend) getHomePage(w http.ResponseWriter, r *http.Request) {
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
	component.Render(ctx, w)
}

func (f *Frontend) getSumonerPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := fromCtx(ctx)

	var (
		region = r.FormValue("region")
	)

	payload := slog.Group("payload", "region", region)

	riotID := r.PathValue("riotID")
	name, tag, err := ParseRiotID(riotID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("failed to resolve riot id", "err", err , payload)
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
		logger.Debug("failed service", "err", err , payload)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.Render(ctx, w)
}

func (f *Frontend) serveSearchResults(w http.ResponseWriter, r *http.Request) {
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
	component.Render(ctx, w)
}

func (f *Frontend) serveMatchlist(w http.ResponseWriter, r *http.Request) {
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
	component.Render(ctx, w)
}

func (f *Frontend) updateSummoner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	var (
		region = r.FormValue("region")
		name   = r.FormValue("name")
		tag    = r.FormValue("tag")
	)

	payload := slog.Group("payload", "region", region, "name", name, "tag", tag)

	logger.Debug("updating summoner", payload)

	riotRegion:= convertStringToRiotRegion(region)

	if err := f.handler.UpdateSummoner(ctx, riotRegion, name, tag); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("service failed", "err", err, payload)
		return
	}

	// Redirect to summoner page
	w.Header().Set("HX-Location", fmt.Sprintf("/%s-%s", name, tag))
	w.WriteHeader(http.StatusOK)
}

func (f *Frontend) serveChampions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	decoded, err := decode[ZGetSummonerChampionsRequest](r)

	payload := slog.Group("payload", "region", decoded.Region, "puuid", decoded.PUUID, "week", decoded.Week)

	component, err := f.handler.ZGetSummonerChampions(ctx, decoded)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.Render(ctx, w)
}

func (f *Frontend) serveLiveMatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decode[GetLiveMatchRequest](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("bad request", "err", err)
		return
	}

	payload := slog.Group("payload", "region", req.Region, "puuid", req.PUUID)

	component, err := f.handler.GetLiveMatch(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Debug("failed service", "err", err, payload)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.Render(ctx, w)
}

func (f *Frontend) addLoggingMiddleware(handler http.Handler) http.Handler {
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

	if index == len(riotID) - 1 {
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
