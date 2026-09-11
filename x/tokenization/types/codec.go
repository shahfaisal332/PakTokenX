package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateHolding{},
		&MsgUpdateHolding{},
		&MsgDeleteHolding{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDistributeRevenue{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBuyTokens{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateProject{},
		&MsgUpdateProject{},
		&MsgDeleteProject{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
