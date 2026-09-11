package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"paktoken/x/tokenization/keeper"
	"paktoken/x/tokenization/types"
)

func mustAddr(t *testing.T, f *fixture, label string) string {
	t.Helper()
	addr, err := f.addressCodec.BytesToString([]byte(label))
	require.NoError(t, err)
	return addr
}

func mustCreateProject(t *testing.T, f *fixture, srv types.MsgServer, creator, totalSupply string, tokenPrice uint64) uint64 {
	t.Helper()
	resp, err := srv.CreateProject(f.ctx, &types.MsgCreateProject{
		Creator:     creator,
		Name:        "test-project",
		TotalSupply: totalSupply,
		Owner:       creator,
		TokenPrice:  tokenPrice,
		Category:    "test",
	})
	require.NoError(t, err)
	return resp.Id
}

func TestBuyTokensCollectsPaymentAndEnforcesReleasedSupply(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	sponsor := mustAddr(t, f, "sponsorAddr____________________")
	buyer := mustAddr(t, f, "buyerAddr_______________________")

	projectId := mustCreateProject(t, f, srv, sponsor, "1000", 2)

	// sponsor issues 2% of the project tokens at any time (after creation)
	_, err := srv.SetReleasedSupply(f.ctx, &types.MsgSetReleasedSupply{Creator: sponsor, ProjectId: projectId, ReleasedSupply: "2%"})
	require.NoError(t, err)

	proj, err := f.keeper.Project.Get(f.ctx, projectId)
	require.NoError(t, err)
	require.Equal(t, "20", proj.ReleasedSupply)

	f.bank.Fund(sdk.MustAccAddressFromBech32(buyer), sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))

	// buy 20 tokens at price 2 -> sponsor receives 40 stake
	_, err = srv.BuyTokens(f.ctx, &types.MsgBuyTokens{Creator: buyer, ProjectId: projectId, Amount: 20})
	require.NoError(t, err)

	require.Equal(t, int64(40), f.bank.SpendableCoins(f.ctx, sdk.MustAccAddressFromBech32(sponsor)).AmountOf("stake").Int64())
	require.Equal(t, int64(960), f.bank.SpendableCoins(f.ctx, sdk.MustAccAddressFromBech32(buyer)).AmountOf("stake").Int64())

	holding, err := f.keeper.Holding.Get(f.ctx, "0-"+buyer)
	require.NoError(t, err)
	require.Equal(t, uint64(20), holding.Amount)

	// released supply cap reached -> cannot buy more
	_, err = srv.BuyTokens(f.ctx, &types.MsgBuyTokens{Creator: buyer, ProjectId: projectId, Amount: 1})
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
}

func TestBuyTokensRejectsZeroAmountAndBeforeRelease(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	sponsor := mustAddr(t, f, "sponsorAddr____________________")
	buyer := mustAddr(t, f, "buyerAddr_______________________")

	projectId := mustCreateProject(t, f, srv, sponsor, "1000", 2)
	f.bank.Fund(sdk.MustAccAddressFromBech32(buyer), sdk.NewCoins(sdk.NewInt64Coin("stake", 5000)))

	// nothing released yet -> rejected
	_, err := srv.BuyTokens(f.ctx, &types.MsgBuyTokens{Creator: buyer, ProjectId: projectId, Amount: 10})
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)

	_, err = srv.SetReleasedSupply(f.ctx, &types.MsgSetReleasedSupply{Creator: sponsor, ProjectId: projectId, ReleasedSupply: "1000"})
	require.NoError(t, err)

	// zero amount -> rejected
	_, err = srv.BuyTokens(f.ctx, &types.MsgBuyTokens{Creator: buyer, ProjectId: projectId, Amount: 0})
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)

	// insufficient funds -> rejected (100 tokens at price 2 = 200, buyer only has 5000 -> ok), so shrink funds
	f.bank.balances[sdk.MustAccAddressFromBech32(buyer).String()] = nil
	f.bank.Fund(sdk.MustAccAddressFromBech32(buyer), sdk.NewCoins(sdk.NewInt64Coin("stake", 150)))

	_, err = srv.BuyTokens(f.ctx, &types.MsgBuyTokens{Creator: buyer, ProjectId: projectId, Amount: 100})
	require.ErrorIs(t, err, sdkerrors.ErrInsufficientFunds)
}

func TestDistributeRevenuePaysOutProportionally(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	sponsor := mustAddr(t, f, "sponsorAddr____________________")
	buyer1 := mustAddr(t, f, "buyer1Addr______________________")
	buyer2 := mustAddr(t, f, "buyer2Addr______________________")

	projectId := mustCreateProject(t, f, srv, sponsor, "1000", 1)

	_, err := srv.SetReleasedSupply(f.ctx, &types.MsgSetReleasedSupply{Creator: sponsor, ProjectId: projectId, ReleasedSupply: "100%"})
	require.NoError(t, err)

	f.bank.Fund(sdk.MustAccAddressFromBech32(buyer1), sdk.NewCoins(sdk.NewInt64Coin("stake", 5000)))
	f.bank.Fund(sdk.MustAccAddressFromBech32(buyer2), sdk.NewCoins(sdk.NewInt64Coin("stake", 5000)))

	_, err = srv.BuyTokens(f.ctx, &types.MsgBuyTokens{Creator: buyer1, ProjectId: projectId, Amount: 10})
	require.NoError(t, err)
	_, err = srv.BuyTokens(f.ctx, &types.MsgBuyTokens{Creator: buyer2, ProjectId: projectId, Amount: 30})
	require.NoError(t, err)

	// sponsor distributes 100 stake; holders split by holding ratio 10 : 30
	f.bank.Fund(sdk.MustAccAddressFromBech32(sponsor), sdk.NewCoins(sdk.NewInt64Coin("stake", 5000)))

	_, err = srv.DistributeRevenue(f.ctx, &types.MsgDistributeRevenue{Creator: sponsor, ProjectId: projectId, Amount: 100})
	require.NoError(t, err)

	// buyer1 paid 10, received 25 -> final balance 5000 - 10 + 25
	require.Equal(t, int64(5015), f.bank.SpendableCoins(f.ctx, sdk.MustAccAddressFromBech32(buyer1)).AmountOf("stake").Int64())
	// buyer2 paid 30, received 75 -> final balance 5000 - 30 + 75
	require.Equal(t, int64(5045), f.bank.SpendableCoins(f.ctx, sdk.MustAccAddressFromBech32(buyer2)).AmountOf("stake").Int64())

	// unauthorized caller cannot distribute
	buyer3 := mustAddr(t, f, "buyer3Addr______________________")
	_, err = srv.DistributeRevenue(f.ctx, &types.MsgDistributeRevenue{Creator: buyer3, ProjectId: projectId, Amount: 10})
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
}
