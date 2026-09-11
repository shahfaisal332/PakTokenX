package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"paktoken/x/tokenization/keeper"
	module "paktoken/x/tokenization/module"
	"paktoken/x/tokenization/types"
)

type mockBankKeeper struct {
	balances map[string]sdk.Coins
}

func (m *mockBankKeeper) SpendableCoins(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	return m.balances[addr.String()]
}

func (m *mockBankKeeper) SendCoins(_ context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	from := m.balances[fromAddr.String()]
	for _, coin := range amt {
		if from.AmountOf(coin.Denom).LT(coin.Amount) {
			return errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, "insufficient funds")
		}
	}
	m.balances[fromAddr.String()] = from.Sub(amt...)
	m.balances[toAddr.String()] = m.balances[toAddr.String()].Add(amt...)
	return nil
}

func (m *mockBankKeeper) Fund(addr sdk.AccAddress, coins sdk.Coins) {
	m.balances[addr.String()] = m.balances[addr.String()].Add(coins...)
}

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	bank         *mockBankKeeper
	addressCodec address.Codec
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)
	bank := &mockBankKeeper{balances: make(map[string]sdk.Coins)}

	k := keeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		authority,
		bank,
	)

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		bank:         bank,
		addressCodec: addressCodec,
	}
}

// authorityAddr returns the bech32 address of the module authority (x/gov).
func (f *fixture) authorityAddr(t *testing.T) string {
	t.Helper()
	auth, err := f.addressCodec.BytesToString(authtypes.NewModuleAddress(types.GovModuleName))
	if err != nil {
		t.Fatalf("failed to encode authority: %v", err)
	}
	return auth
}
