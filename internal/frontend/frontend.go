// Hi
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

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

var (
	ErrInvalidRegion = errors.New("invalid region")
)

// Frontend serves [templ.Component].
type Frontend struct {
	router *http.ServeMux

	datasource *internal.Datasource

	store internal.Store
}

func New(store internal.Store, datasource *internal.Datasource) *Frontend {
	router := http.NewServeMux()

	frontend := Frontend{
		router:     router,
		store:      store,
		datasource: datasource,
	}

	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	router.HandleFunc("GET /", frontend.getHomePage)

	router.HandleFunc("GET /summoner/{platform}/{name}/{tagline}", frontend.getSumonerPage)

	router.HandleFunc("POST /search", frontend.serveSearchResults)

	router.HandleFunc("POST /summoner/fetch", frontend.updateSummoner)

	router.HandleFunc("POST /summoner/matchlist", frontend.serveMatchlist)

	router.HandleFunc("POST /summoner/live", frontend.serveLiveMatch)

	router.HandleFunc("POST /summoner/champions", frontend.serveChampions)

	return &frontend
}

func (f *Frontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.router.ServeHTTP(w, r)
}

// UpdateSummoner
func (f *Frontend) UpdateSummoner(ctx context.Context, region riot.Region, name, tag string) error {
	puuid, err := f.datasource.GetPUUID(ctx, name, tag)
	if err != nil {
		return err
	}

	if err := f.datasource.UpdateSummoner(ctx, region, puuid); err != nil {
		return err
	}

	return nil
}

// GetLiveMatch the live match view if summoner is in a game in the region. If
// no such game is found, return a view indicating such.
func (f *Frontend) GetLiveMatch(ctx context.Context, region riot.Region, puuid string) (view templ.Component, err error) {
	match, err := f.datasource.GetLiveMatch(ctx, region, puuid)
	if err != nil {
		if errors.Is(err, internal.ErrNoLiveMatch) {
			return NoLiveMatchModalWindow{}, err
		}

		return LiveMatchModalWindow{}, err
	}

	redSide := []LiveMatchSummonerCard{}
	blueSide := []LiveMatchSummonerCard{}

	for _, p := range match.Participants {
		card := LiveMatchSummonerCard{
			Champion:  p.ChampionID,
			Summoners: p.SummonersIDs,
			RunePage:  p.Runes,
			Name:      "Doublelift",
			Tag:       "",
			Rank:      &internal.RankDetail{},
			TeamID:    p.TeamID,
		}
		if p.TeamID == 100 {
			redSide = append(redSide, card)
		} else {
			blueSide = append(blueSide, card)
		}
	}

	window := LiveMatchModalWindow{
		AverageRank: &internal.RankDetail{},
		StartTime:   match.Date,
		RedSide:     redSide,
		BlueSide:    blueSide,
	}

	return window, nil
}

// GetHomePage returns [HomePage].
func (f *Frontend) GetHomePage(ctx context.Context) (templ.Component, error) {
	page := HomePage{}
	return page, nil
}

// GetSummonerMatchHistoryPage returns a [SummonerPage].
func (f *Frontend) GetSummonerPage(ctx context.Context, region riot.Region, name, tag string) (templ.Component, error) {
	puuid, err := f.datasource.GetPUUID(ctx, name, tag)
	if err != nil {
		return nil, err
	}

	summoner, err := f.store.GetSummoner(ctx, puuid)
	if err != nil {
		return nil, err
	}

	rank, err := f.store.GetRank(ctx, puuid, time.Now(), true)
	if err != nil {
		return nil, err
	}

	page := SummonerPage{
		PUUID:       puuid,
		Name:        summoner.Name,
		Tag:         summoner.Tagline,
		Rank:        rank.Detail,
		LastUpdated: rank.EffectiveDate,
	}

	return page, nil
}

// GetSummonerChampions returns a [ChampionsModal].
func (f *Frontend) GetSummonerChampions(ctx context.Context, region riot.Region, puuid string) (templ.Component, error) {
	champions, err := f.store.GetChampions(ctx, puuid)
	if err != nil {
		return nil, err
	}

	modal := ChampionsModal{
		Cards: []SummonerChampionCard{},
	}

	for _, champion := range champions {
		modal.Cards = append(modal.Cards, SummonerChampionCard{
			Champion:             int(champion.Champion),
			Kills:                champion.Kills,
			Deaths:               champion.Deaths,
			Assists:              champion.Assists,
			KillParticipation:    champion.KillParticipation,
			CS:                   champion.CreepScore,
			CSPerMinute:          champion.CreepScorePerMinute,
			DamageDealt:          champion.DamageDealt,
			DamageTaken:          champion.DamageTaken,
			DamageDeltaEnemy:     champion.DamageDeltaEnemy,
			DamagePercentageTeam: champion.DamagePercentageTeam,
			GoldEarned:           champion.GoldEarned,
			GoldDeltaEnemy:       champion.GoldDeltaEnemy,
			GoldPercentageTeam:   champion.GoldPercentageTeam,
			VisionScore:          champion.VisionScore,
			PinkWardsBought:      champion.PinkWardsBought,
		})
	}

	return modal, nil
}

// GetSummonerMatchHistoryList returns a list of 10 [MatchHistoryCard] with an
// additional component that will fetch the next page. It will fetch from riot
// to ensure the store is synced.
func (f *Frontend) GetSummonerMatchHistoryList(ctx context.Context, region riot.Region, puuid string, page int) (templ.Component, error) {
	if err := f.datasource.UpdateMatchHistory(ctx, region, puuid, page*10, 10); err != nil {
		return nil, err
	}

	storeMatches, err := f.store.GetMatches(ctx, puuid, page)
	if err != nil {
		return nil, fmt.Errorf("storage failure: %w", err)
	}

	cards := []MatchHistoryCard{}

	for _, m := range storeMatches {
		cards = append(cards, MatchHistoryCard{
			Champion:    m.ChampionID,
			Summoners:   m.SummonerIDs,
			Kills:       m.Kills,
			Deaths:      m.Deaths,
			Assists:     m.Assists,
			CS:          m.CreepScore,
			CSPerMinute: m.CreepScorePerMinute,
			RunePage:    m.Runes,
			Items:       m.Items,
		})
	}

	// blocks := groupMatchesByDay(storeMatches)

	// days := []List{}

	// for _, block := range blocks {
	// 	if len(block) == 0 {
	// 		continue
	// 	}

	// 	day := List{
	// 		Title: block[0].Date.Truncate(24 * time.Hour).Format("Monday, Jan 2"),
	// 		Items: []struct{ ListItemChildren []templ.Component }{},
	// 	}

	// 	for _, m := range block {
	// 		scoreboardRow := struct {
	// 			ListItemChildren []templ.Component
	// 		}{
	// 			ListItemChildren: []templ.Component{
	// 				ChampionWidget{
	// 					Champion:  m.ChampionID,
	// 					Summoners: &m.SummonerIDs,
	// 				},
	// 				TextKDA{
	// 					Kills:   m.Kills,
	// 					Deaths:  m.Deaths,
	// 					Assists: m.Assists,
	// 				},
	// 				Text{
	// 					S:     fmt.Sprintf("%d (%.1f)", m.CreepScore, m.CreepScorePerMinute),
	// 					Width: "w-24",
	// 				},
	// 				RuneWidget{
	// 					RunePage: m.Runes,
	// 				},
	// 				ItemWidget{
	// 					Items: m.Items,
	// 				},
	// 				RankDeltaWidget{},
	// 				LinkButton{
	// 					HREF:     templ.URL(fmt.Sprintf("/match/%s", m.MatchID)),
	// 					IconPath: StaticExternalLinkSymbol,
	// 				},
	// 			},
	// 		}

	// 		day.Items = append(day.Items, scoreboardRow)
	// 	}

	// 	days = append(days, day)
	// }

	var components []templ.Component

	for _, card := range cards {
		components = append(components, card)
	}

	return templ.Join(components...), nil
}

// GetMatchScoreboard returns the scoreboard of a match.
func (f *Frontend) GetMatchScoreboard(ctx context.Context, id string) (scoreboard templ.Component, err error) {
	_, participants, err := f.store.GetMatch(ctx, id)
	if err != nil {
		return nil, err
	}

	redSide := List{
		Title: "Red Side",
		Items: []struct{ ListItemChildren []templ.Component }{},
	}

	blueSide := List{
		Title: "Blue Side",
		Items: []struct{ ListItemChildren []templ.Component }{},
	}

	for _, p := range participants {

		particpantRow := struct {
			ListItemChildren []templ.Component
		}{
			ListItemChildren: []templ.Component{
				ChampionWidget{
					Champion:  p.ChampionID,
					Summoners: &p.SummonerIDs,
				},
				TextKDA{
					Kills:   p.Kills,
					Deaths:  p.Deaths,
					Assists: p.Assists,
				},
				Text{
					S:     fmt.Sprintf("%d (%.1f)", p.CreepScore, p.CreepScorePerMinute),
					Width: "w-24",
				},
				RuneWidget{
					RunePage: p.Runes,
				},
				ItemWidget{
					Items: p.Items,
				},
			},
		}

		if p.TeamID == 100 {
			blueSide.Items = append(blueSide.Items, particpantRow)
		} else {
			redSide.Items = append(redSide.Items, particpantRow)
		}
	}

	scoreboard = templ.Join(blueSide, redSide)
	return scoreboard, nil
}

// GetSearchResults returns a list of [SearchResultCard] for accounts that
// match q. If no results were found, return [SearchNotFoundCard] instead.
//
// q should be of the form name#tag, if q has no tag, region is used as the
// tag.
//
// POST /search
func (f *Frontend) GetSearchResults(ctx context.Context, region riot.Region, q string) (results templ.Component, err error) {
	storeSearchResults, err := f.store.SearchSummoner(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(storeSearchResults) == 0 {
		var name, tag string
		if i := strings.Index(q, "#"); i != -1 {
			name = q[:i]
			if i+1 == len(q) {
				tag = string(region)
			} else {
				tag = q[i+1:]
			}
		} else {
			name = q
			tag = string(region)
		}

		return SearchNotFoundCard{
			Name: name,
			Tag:  tag,
		}, nil
	}

	searchResults := []SearchResultCard{}

	for _, r := range storeSearchResults {
		rank, err := f.store.GetRank(ctx, r.Puuid, time.Now(), true)
		if err != nil {
			return nil, fmt.Errorf("getting rank for %s#%s: %w", r.Name, r.Tagline, err)
		}

		row := SearchResultCard{
			PUUID: r.Puuid,
			Name:  r.Name,
			Tag:   r.Tagline,
			Rank:  rank.Detail,
		}

		searchResults = append(searchResults, row)
	}

	v := []templ.Component{}
	for _, r := range searchResults {
		v = append(v, r)
	}

	return templ.Join(v...), nil
}

func (f *Frontend) getHomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	component, err := f.GetHomePage(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	component.Render(ctx, w)
}

func (f *Frontend) getSumonerPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		region = r.PathValue("region")
		name = r.PathValue("name")
		tag  = r.PathValue("tag")
	)

	riotRegion, err := convertStringToRiotRegion(region)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("failed to parse region", slog.Group("payload", "region", region, "name", name, "tag", tag))
		return
	}

	component, err := f.GetSummonerPage(ctx, riotRegion, name, tag)
	if err != nil {
		if errors.Is(err, internal.ErrSummonerNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
		}

		slog.Debug("getting summoner", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
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

	component, err := f.GetSearchResults(ctx, riotRegion, q)
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
		region = r.FormValue("region")
		puuid     = r.FormValue("puuid")
		page  int = 0
	)

	riotRegion, err := convertStringToRiotRegion(region)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("failed to parse region", slog.Group("payload", "region", region, "puuid", puuid, "page", page))
		return
	}

	if pageQuery := r.FormValue("page"); pageQuery != "" {
		if pageVal, err := strconv.Atoi(pageQuery); err == nil {
			page = pageVal
		}
	}

	component, err := f.GetSummonerMatchHistoryList(ctx, riotRegion, puuid, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	riotRegion, err := convertStringToRiotRegion(region)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("failed to parse region", "err", err, slog.Group("payload", "region", region, "name", name, "tag", tag))
		return
	}

	if err := f.UpdateSummoner(ctx, riotRegion, name, tag); err != nil {
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
		puuid = r.PathValue("puuid")
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


	component, err := f.GetSummonerChampions(ctx, riotRegion, puuid)
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

	component, err := f.GetLiveMatch(ctx, riotRegion, puuid)
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

func groupMatchesByDay(matches []internal.SummonerMatch) [][]internal.SummonerMatch {
	var blocks [][]internal.SummonerMatch

	if len(matches) == 0 {
		return blocks
	}

	last := 0
	for i := range matches {
		if matches[last].Date.Truncate(24*time.Hour) == matches[i].Date.Truncate(24*time.Hour) {
			continue
		} else {
			blocks = append(blocks, matches[last:i-1])
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
