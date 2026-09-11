package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"paktoken/x/tokenization/keeper"
	"paktoken/x/tokenization/types"
)

func createNHolding(keeper keeper.Keeper, ctx context.Context, n int) []types.Holding {
	items := make([]types.Holding, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		items[i].Amount = uint64(i)
		items[i].ProjectId = uint64(i)
		items[i].Investor = strconv.Itoa(i)
		_ = keeper.Holding.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestHoldingQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNHolding(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetHoldingRequest
		response *types.QueryGetHoldingResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetHoldingRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetHoldingResponse{Holding: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetHoldingRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetHoldingResponse{Holding: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetHoldingRequest{
				Index: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetHolding(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestHoldingQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNHolding(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllHoldingRequest {
		return &types.QueryAllHoldingRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListHolding(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Holding), step)
			require.Subset(t, msgs, resp.Holding)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListHolding(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Holding), step)
			require.Subset(t, msgs, resp.Holding)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListHolding(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Holding)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListHolding(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
