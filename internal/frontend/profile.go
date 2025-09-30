package frontend

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/page"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/view/profile"
)

type ProfileService struct {
	Handler *ProfileHandler
}

func (s *ProfileService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /summoner/fetch", s.updateSummoner)

	router.HandleFunc("POST /summoner/matchlist", s.serveMatchlist)

	router.HandleFunc("POST /summoner/live", s.serveLiveMatch)

	router.HandleFunc("POST /summoner/champions", s.serveChampions)

	router.HandleFunc("POST /match", s.serveMatchDetail)
}

func (s *ProfileService) serveSumonerPage(w http.ResponseWriter, r *http.Request) {
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

	data, async, err := s.Handler.GetSummonerPage(ctx, riotRegion, name, tag, time.UTC)
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

	component := page.ProfilePage(ctx, *data)

	go async()
	templ.Handler(component, templ.WithStreaming()).ServeHTTP(w, r)
}

func (s *ProfileService) serveMatchlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	req, err := decode[MatchHistoryRequest](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("bad request", "err", err)
		return
	}

	payload := slog.Any("request", req)

	component, err := s.Handler.GetMatchHistory(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	if err := profile.HistoryEntry(ctx, *component).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed rendering", "err", err)
		return
	}
}

func (s *ProfileService) updateSummoner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	decoded, _ := decode[UpdateSummonerRequest](r)

	payload := slog.Group("payload", "region", decoded.Region, "name", decoded.Name, "tag", decoded.Tag)

	if err := s.Handler.UpdateSummoner(ctx, decoded); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("service failed", "err", err, payload)
		return
	}

	// Redirect to summoner page
	w.Header().Set("HX-Location", fmt.Sprintf("/%s-%s", decoded.Name, decoded.Tag))
	w.WriteHeader(http.StatusOK)
}

func (s *ProfileService) serveChampions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	decoded, _ := decode[GetSummonerChampionsRequest](r)

	payload := slog.Group("payload", "region", decoded.Region, "puuid", decoded.PUUID, "week", decoded.Week)

	v, err := s.Handler.GetSummonerChampions(ctx, decoded)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	if err := profile.ChampionList(ctx, *v).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed rendering", "err", err)
		return
	}
}

func (s *ProfileService) serveLiveMatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	req, err := decode[LiveMatchRequest](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("bad request", "err", err)
		return
	}

	payload := slog.Group("payload", "region", req.Region, "puuid", req.PUUID)

	v, err := s.Handler.GetLiveMatch(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	if err := profile.LiveMatchModal(ctx, *v).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed rendering", "err", err)
		return
	}
}

func (s *ProfileService) serveMatchDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	req, err := decode[MatchDetailRequest](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("bad request", "err", err)
		return
	}

	payload := slog.Group("payload", "region", req.Region, "match_id", req.MatchID)

	v, err := s.Handler.GetMatchDetail(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	// TODO: adjust timezone here

	if err := profile.MatchDetailBox(ctx, *v).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed rendering", "err", err)
		return
	}
}

type ProfileHandler struct {
	Datasource *internal.Datasource
}

func (h *ProfileHandler) GetSummonerPage(ctx context.Context, region riot.Region, name, tag string, tz *time.Location) (*page.ProfilePageData, func(), error) {
	storeProfile, err := h.Datasource.GetProfileDetailByRiotID(ctx, region, name, tag)
	if err != nil {
		return nil, nil, err
	}

	data := page.ProfilePageData{
		PUUID:          storeProfile.PUUID,
		Region:         region,
		Name:           name,
		Tag:            tag,
		HistoryEntryCh: make(chan profile.HistoryEntryData),
		RankCardCh:     make(chan profile.RankCardData),
		ChampionListCh: make(chan profile.ChampionListData),
	}

	grr := func() {
		defer close(data.HistoryEntryCh)
		defer close(data.ChampionListCh)
		defer close(data.RankCardCh)

		days := GetDays(time.Now())

		for i := range len(days) - 1 {
			historyEntryData, err := h.GetMatchHistory(ctx, MatchHistoryRequest{
				Region:  region,
				PUUID:   data.PUUID,
				StartTS: days[i+1],
				EndTS:   days[i],
			})

			if err == nil {
				data.HistoryEntryCh <- *historyEntryData
			}
		}

		championListData, err := h.GetSummonerChampions(ctx, GetSummonerChampionsRequest{
			Region: region,
			PUUID:  data.PUUID,
			Week:   GetCurrentWeek(),
		})

		if err == nil {
			data.ChampionListCh <- *championListData
		}
	}

	return &data, grr, nil
}

// GetMatchHistory returns matches played between the given timestamps. The
// method will fetch riot first to ensure all matches played on date are in
// store.
func (h *ProfileHandler) GetMatchHistory(ctx context.Context, req MatchHistoryRequest) (*profile.HistoryEntryData, error) {
	start := req.StartTS
	end := req.EndTS

	storeMatches, err := h.Datasource.GetMatchHistory(ctx, req.Region, req.PUUID, start, end)
	if err != nil {
		return nil, fmt.Errorf("storage failure: %w", err)
	}

	v := &profile.HistoryEntryData{
		Date:      start,
		Matchlist: []profile.HistoryCardData{},
	}

	for _, match := range storeMatches {
		path, data := makeGetMatchDetailRequest(req.Region, match.MatchID)

		kda := float32(match.Kills+match.Assists) / float32(match.Deaths)

		v.Matchlist = append(v.Matchlist, profile.HistoryCardData{
			ChampionID:     match.ChampionID,
			ChampionLevel:  match.ChampionLevel,
			SummonerIDs:    match.SummonerIDs,
			Kills:          match.Kills,
			Deaths:         match.Deaths,
			Assists:        match.Assists,
			KillDeathRatio: kda,
			CS:             match.CreepScore,
			CSPerMinute:    match.CreepScorePerMinute,
			RunePage:       match.Runes,
			Items:          match.Items,
			VisionScore:    match.VisionScore,
			RankChange:     nil,
			LPChange:       nil,
			Win:            match.Win,
			Path:           path,
			Data:           string(data),
		})
	}

	return v, nil
}

func (h *ProfileHandler) GetMatchDetail(ctx context.Context, req MatchDetailRequest) (*profile.MatchDetailBoxData, error) {
	matchDetail, err := h.Datasource.GetMatchDetail(ctx, req.Region, req.MatchID)
	if err != nil {
		return nil, err
	}

	v := profile.MatchDetailBoxData{
		Date:     matchDetail.Date,
		Duration: matchDetail.Duration,
		RedSide: profile.MatchTeamListData{
			Participants: []profile.MatchParticipantCardData{},
		},
		BlueSide: profile.MatchTeamListData{
			Participants: []profile.MatchParticipantCardData{},
		},
	}

	for i := range 5 {
		v.BlueSide.Participants = append(v.BlueSide.Participants, *profile.NewMatchParticipantCardData(matchDetail.Participants[i]))
	}

	for i := range 5 {
		v.RedSide.Participants = append(v.RedSide.Participants, *profile.NewMatchParticipantCardData(matchDetail.Participants[5+i]))
	}

	return &v, nil
}

func (h *ProfileHandler) GetLiveMatch(ctx context.Context, req LiveMatchRequest) (*profile.LiveMatchModalData, error) {
	match, err := h.Datasource.GetLiveMatch(ctx, req.Region, req.PUUID)
	if err != nil {
		return nil, err
	}

	v := profile.LiveMatchModalData{
		Date:         match.Date,
		Participants: [10]internal.LiveParticipant{},
	}

	for i, p := range match.Participants {
		v.Participants[i] = p
	}

	return &v, nil
}

func (h *ProfileHandler) UpdateSummoner(ctx context.Context, req UpdateSummonerRequest) error {
	if err := h.Datasource.UpdateProfileByRiotID(ctx, req.Region, req.Name, req.Tag); err != nil {
		return err
	}

	return nil
}

func (h *ProfileHandler) GetSummonerChampions(ctx context.Context, req GetSummonerChampionsRequest) (*profile.ChampionListData, error) {
	start := req.Week
	end := start.Add(7 * 24 * time.Hour)

	storeChampions, err := h.Datasource.GetSummonerChampions(ctx, req.Region, req.PUUID, start, end)
	if err != nil {
		return nil, err
	}

	v := profile.ChampionListData{
		Champions: []profile.ChampionItemData{},
	}

	for _, champion := range storeChampions {
		v.Champions = append(v.Champions, profile.ChampionItemData{
			Champion:    int(champion.Champion),
			GamesPlayed: champion.GamesPlayed,
			Wins:        champion.Wins,
			Losses:      champion.Losses,
			LPDelta:     nil,
		})
	}

	return &v, nil
}
