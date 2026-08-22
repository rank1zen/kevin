package profile

import (
	"context"

	"buf.build/gen/go/kevin-labs/lol-service/connectrpc/go/kevin/lolservice/v1/lolservicev1connect"
	lolservicev1 "buf.build/gen/go/kevin-labs/lol-service/protocolbuffers/go/kevin/lolservice/v1"
	"connectrpc.com/connect"
)

type Handler struct {
	lolservicev1connect.UnimplementedProfileServiceHandler

	service *Service
}

func (s Handler) GetProfile(ctx context.Context, c *connect.Request[lolservicev1.GetProfileRequest]) (*connect.Response[lolservicev1.GetProfileResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s Handler) RefreshProfile(ctx context.Context, c *connect.Request[lolservicev1.RefreshProfileRequest]) (*connect.Response[lolservicev1.RefreshProfileResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s Handler) SearchProfile(ctx context.Context, c *connect.Request[lolservicev1.SearchProfileRequest]) (*connect.Response[lolservicev1.SearchProfileResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s Handler) GetMatch(ctx context.Context, c *connect.Request[lolservicev1.GetMatchRequest]) (*connect.Response[lolservicev1.GetMatchResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s Handler) GetMatchHistory(ctx context.Context, c *connect.Request[lolservicev1.GetMatchHistoryRequest]) (*connect.Response[lolservicev1.GetMatchHistoryResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s Handler) GetRankHistory(ctx context.Context, c *connect.Request[lolservicev1.GetRankHistoryRequest]) (*connect.Response[lolservicev1.GetRankHistoryResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s Handler) GetChampionAggregate(ctx context.Context, c *connect.Request[lolservicev1.GetChampionAggregateRequest]) (*connect.Response[lolservicev1.GetChampionAggregateResponse], error) {
	//TODO implement me
	panic("implement me")
}
