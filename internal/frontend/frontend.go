package frontend

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rank1zen/kevin/internal"
)

var (
	ErrInvalidRegion = errors.New("invalid region")
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

	router.HandleFunc("GET /summoner/{region}/{name}/{tag}", frontend.getSumonerPage)

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

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	component, err := f.handler.GetHomePage(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.Render(ctx, w)
}

func (f *Frontend) getSumonerPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := fromCtx(ctx)

	var (
		region = r.PathValue("region")
		name   = r.PathValue("name")
		tag    = r.PathValue("tag")
	)

	riotRegion, err := convertStringToRiotRegion(region)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("failed to parse region",
			slog.Group("payload", "region", region, "name", name, "tag", tag),
		)

		return
	}

	component, err := f.handler.GetSummonerPage(ctx, riotRegion, name, tag)
	if err != nil {
		if errors.Is(err, internal.ErrSummonerNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.Render(ctx, w)
}

func (f *Frontend) serveSearchResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		region = r.FormValue("region")
		q      = r.FormValue("q")
	)

	riotRegion, err := convertStringToRiotRegion(region)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("failed to parse region", slog.Group("payload", "region", region, "q", q))
		return
	}

	component, err := f.handler.GetSearchResults(ctx, riotRegion, q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Debug("failed service", slog.Any("err", err), slog.Group("payload", "region", region, "q", q))
		return
	}

	w.WriteHeader(http.StatusOK)
	component.Render(ctx, w)
}

func (f *Frontend) serveMatchlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		region           = r.FormValue("region")
		puuid            = r.FormValue("puuid")
		date   time.Time = time.Now()
	)

	if dateQuery := r.FormValue("date"); dateQuery != "" {
		if dateVal, err := strconv.ParseInt(dateQuery, 10, 64); err == nil {
			date = time.Unix(dateVal, 0)
		}
	}

	riotRegion, err := convertStringToRiotRegion(region)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("failed to parse region",
			slog.Group("payload", "region", region, "puuid", puuid, "date", r.FormValue("date")),
		)
		return
	}

	component, err := f.handler.GetSummonerMatchHistoryList(ctx, riotRegion, puuid, date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Debug("failed service", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.Render(ctx, w)
}

func (f *Frontend) updateSummoner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		region = r.FormValue("region")
		name   = r.FormValue("name")
		tag    = r.FormValue("tag")
	)
	slog.Debug("updating summoner", slog.Group("payload", "region", region, "name", name, "tag", tag))

	riotRegion, err := convertStringToRiotRegion(region)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("failed to parse region",
			"err", err,
			slog.Group("payload", "region", region, "name", name, "tag", tag),
		)
		return
	}

	if err := f.handler.UpdateSummoner(ctx, riotRegion, name, tag); err != nil {
		slog.Debug("service failed",
			slog.Any("err", err),
			slog.Group("payload", "region", region, "name", name, "tag", tag),
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Redirect to summoner page
	w.Header().Set("HX-Location", fmt.Sprintf("/summoner/%s/%s/%s", region, name, tag))
	w.WriteHeader(http.StatusOK)
}

func (f *Frontend) serveChampions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		region = r.PathValue("puuid")
		puuid  = r.PathValue("puuid")
	)

	riotRegion, err := convertStringToRiotRegion(region)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		slog.Debug("failed to parse region",
			slog.Any("err", err),
			slog.Group("payload", "region", region, "puuid", puuid),
		)

		return
	}

	component, err := f.handler.GetSummonerChampions(ctx, riotRegion, puuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.Render(ctx, w)
}

func (f *Frontend) serveLiveMatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		region = r.PathValue("puuid")
		puuid  = r.PathValue("puuid")
	)

	riotRegion, err := convertStringToRiotRegion(region)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		slog.Debug("failed to parse region",
			slog.Any("err", err),
			slog.Group("payload", "region", region, "puuid", puuid),
		)

		return
	}

	component, err := f.handler.GetLiveMatch(ctx, riotRegion, puuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		slog.Debug("failed service",
			slog.Any("err", err),
			slog.Group("payload", "region", region, "puuid", puuid),
		)

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

func groupMatchesByDay(matches []internal.SummonerMatch) [][]internal.SummonerMatch {
	if len(matches) == 0 {
		return [][]internal.SummonerMatch{}
	}

	times := []time.Time{}
	for _, t := range matches {
		times = append(times, t.Date)
	}

	blocks := [][]internal.SummonerMatch{}
	indices := groupTimeByDay(times)

	if len(indices) == 0 {
		blocks = append(blocks, matches[:])
		return blocks
	}

	blocks = append(blocks, matches[:indices[0]])

	for i := 1; i < len(indices); i++ {
		blocks = append(blocks, matches[indices[i-1]:indices[i]])
	}

	blocks = append(blocks, matches[indices[len(indices)-1]:])

	return blocks
}

func groupTimeByDay(times []time.Time) []int {
	blocks := []int{}

	if len(times) == 0 {
		return blocks
	}

	last := 0
	for i := range times {
		if times[i].Truncate(24*time.Hour) == times[last].Truncate(24*time.Hour) {
			continue
		} else {
			blocks = append(blocks, i)
			last = i
		}
	}

	return blocks
}

func joinClasses(classes []string) string {
	return strings.Join(classes, " ")
}

type classBuilder struct {
	class []string
}

func (cb classBuilder) Add(class string) classBuilder {
	cb.class = append(cb.class, class)
	return cb
}

func (cb classBuilder) Build() string {
	return joinClasses(cb.class)
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

	return nil
}
