package keeper

import (
	"context"
	"errors"

	"paktoken/x/tokenization/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListHolding(ctx context.Context, req *types.QueryAllHoldingRequest) (*types.QueryAllHoldingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	holdings, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Holding,
		req.Pagination,
		func(_ string, value types.Holding) (types.Holding, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllHoldingResponse{Holding: holdings, Pagination: pageRes}, nil
}

func (q queryServer) GetHolding(ctx context.Context, req *types.QueryGetHoldingRequest) (*types.QueryGetHoldingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Holding.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetHoldingResponse{Holding: val}, nil
}
