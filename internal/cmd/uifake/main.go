package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/logging"
	"github.com/rank1zen/yujin/internal/ui"
	"go.uber.org/zap"
)

type mock struct{}

func (s *mock) CheckProfileExists(_ context.Context, _ internal.PUUID) (bool, error) {
	return true, nil
}

func (s *mock) UpdateProfile(_ context.Context, _ internal.Profile) error {
	return nil
}

func (s *mock) GetProfile(_ context.Context, _ internal.PUUID) (internal.Profile, error) {
	var m internal.Profile
	gofakeit.Struct(&m)
	return m, nil
}

func (s *mock) GetRankList(_ context.Context, _ internal.PUUID) ([]internal.RankRecord, error) {
	m := make([]internal.RankRecord, 10)
	for _, n := range m {
		gofakeit.Struct(&n)
	}
	return m, nil
}

func (s *mock) GetMatchList(_ context.Context, _ internal.PUUID, _ int, _ bool) ([]internal.MatchParticipant, error) {
	m := make([]internal.MatchParticipant, 10)
	for i := range len(m) {
		err := gofakeit.Struct(&(m[i]))
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (s *mock) GetChampionList(_ context.Context, _ internal.PUUID) ([]internal.ChampionStats, error) {
	m := make([]internal.ChampionStats, 10)
	for _, n := range m {
		gofakeit.Struct(&n)
	}
	return m, nil
}

func (s *mock) CreateMatch(_ context.Context, _ internal.Match) error {
	return nil
}

type mock2 struct{}

func (s *mock2) GetProfile(_ context.Context, _ internal.PUUID) (_ internal.Profile, _ error) {
	var m internal.Profile
	gofakeit.Struct(&m)
	return m, nil
}

func (s *mock2) GetLiveMatch(_ context.Context, _ internal.PUUID) (_ internal.LiveMatch, _ error) {
	var m internal.LiveMatch
	gofakeit.Struct(&m)
	return m, nil
}

func (s *mock2) GetMatchList(_ context.Context, _ internal.PUUID) (_ []internal.MatchID, _ error) {
	var m []internal.MatchID
	gofakeit.Struct(m)
	return m, nil
}

func (s *mock2) GetMatch(_ context.Context, _ internal.MatchID) (_ internal.Match, _ error) {
	var m internal.Match
	gofakeit.Struct(&m)
	return m, nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	logger := logging.Get()
	defer logger.Sync()

	var port int
	port = 4001
	if os.Getenv("YUJIN_PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("YUJIN_PORT"))
	}

	ui := ui.Routes(&mock{}, &mock2{})

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), ui)
	}()

	logger.Sugar().Infof("started server on %v", zap.Int("port", port))

	<-ctx.Done()

	logger.Sugar().Info("shutting down")
}
