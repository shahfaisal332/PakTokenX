package keeper

import (
	"context"
	"errors"
	"fmt"

	"paktoken/x/tokenization/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// assertHoldingAuthority ensures only the module authority (x/gov) can manage
// holding records directly. Regular users must go through BuyTokens so that
// payments and released-supply limits are always enforced.
func (k msgServer) holdingAuthority(ctx context.Context, signer string) error {
	authority, err := k.addressCodec.BytesToString(k.authority)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrLogic, "failed to encode authority")
	}
	if signer != authority {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "only the module authority may manage holdings; use BuyTokens instead")
	}
	return nil
}

func (k msgServer) CreateHolding(ctx context.Context, msg *types.MsgCreateHolding) (*types.MsgCreateHoldingResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}
	if err := k.holdingAuthority(ctx, msg.Creator); err != nil {
		return nil, err
	}

	// Check if the value already exists
	ok, err := k.Holding.Has(ctx, msg.Index)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var holding = types.Holding{
		Creator:   msg.Creator,
		Index:     msg.Index,
		Amount:    msg.Amount,
		ProjectId: msg.ProjectId,
		Investor:  msg.Investor,
	}

	if err := k.Holding.Set(ctx, holding.Index, holding); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateHoldingResponse{}, nil
}

func (k msgServer) UpdateHolding(ctx context.Context, msg *types.MsgUpdateHolding) (*types.MsgUpdateHoldingResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	if err := k.holdingAuthority(ctx, msg.Creator); err != nil {
		return nil, err
	}

	// Check if the value exists
	val, err := k.Holding.Get(ctx, msg.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	var holding = types.Holding{
		Creator:   val.Creator,
		Index:     msg.Index,
		Amount:    msg.Amount,
		ProjectId: msg.ProjectId,
		Investor:  msg.Investor,
	}

	if err := k.Holding.Set(ctx, holding.Index, holding); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update holding")
	}

	return &types.MsgUpdateHoldingResponse{}, nil
}

func (k msgServer) DeleteHolding(ctx context.Context, msg *types.MsgDeleteHolding) (*types.MsgDeleteHoldingResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	if err := k.holdingAuthority(ctx, msg.Creator); err != nil {
		return nil, err
	}

	// Check if the value exists
	_, err := k.Holding.Get(ctx, msg.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	if err := k.Holding.Remove(ctx, msg.Index); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove holding")
	}

	return &types.MsgDeleteHoldingResponse{}, nil
}
