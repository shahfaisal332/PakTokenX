package keeper_test

import (
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"paktoken/x/tokenization/keeper"
	"paktoken/x/tokenization/types"
)

func TestHoldingMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	// regular users must not be able to mint holdings directly
	_, err = srv.CreateHolding(f.ctx, &types.MsgCreateHolding{Creator: creator, Index: "0"})
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	// the module authority can still manage holdings
	auth := f.authorityAddr(t)
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateHolding{Creator: auth, Index: strconv.Itoa(i)}
		_, err := srv.CreateHolding(f.ctx, expected)
		require.NoError(t, err)
		rst, err := f.keeper.Holding.Get(f.ctx, expected.Index)
		require.NoError(t, err)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestHoldingMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	auth := f.authorityAddr(t)
	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	_, err = srv.CreateHolding(f.ctx, &types.MsgCreateHolding{Creator: auth, Index: strconv.Itoa(0)})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgUpdateHolding
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgUpdateHolding{Creator: "invalid",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgUpdateHolding{Creator: creator,
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "key not found",
			request: &types.MsgUpdateHolding{Creator: auth,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgUpdateHolding{Creator: auth,
				Index: strconv.Itoa(0),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdateHolding(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, err := f.keeper.Holding.Get(f.ctx, tc.request.Index)
				require.NoError(t, err)
				require.Equal(t, auth, rst.Creator)
			}
		})
	}
}

func TestHoldingMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	auth := f.authorityAddr(t)
	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	_, err = srv.CreateHolding(f.ctx, &types.MsgCreateHolding{Creator: auth, Index: strconv.Itoa(0)})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeleteHolding
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgDeleteHolding{Creator: "invalid",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgDeleteHolding{Creator: creator,
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "key not found",
			request: &types.MsgDeleteHolding{Creator: auth,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgDeleteHolding{Creator: auth,
				Index: strconv.Itoa(0),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeleteHolding(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				found, err := f.keeper.Holding.Has(f.ctx, tc.request.Index)
				require.NoError(t, err)
				require.False(t, found)
			}
		})
	}
}
