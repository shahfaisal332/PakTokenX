package keeper

import (
	"context"
	"fmt"
	"strconv"

	"paktoken/x/tokenization/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) BuyTokens(ctx context.Context, msg *types.MsgBuyTokens) (*types.MsgBuyTokensResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than zero")
	}

	project, err := k.Project.Get(ctx, msg.ProjectId)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("project %d not found", msg.ProjectId))
	}

	totalSupply, err := strconv.ParseUint(project.TotalSupply, 10, 64)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid project total supply")
	}

	releasedSupply, err := strconv.ParseUint(project.ReleasedSupply, 10, 64)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid project released supply")
	}
	if releasedSupply == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no tokens have been released for this project yet")
	}

	var alreadySold uint64
	err = k.Holding.Walk(ctx, nil, func(idx string, h types.Holding) (bool, error) {
		if h.ProjectId == msg.ProjectId {
			alreadySold += h.Amount
		}
		return false, nil
	})
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to read existing holdings")
	}

	// overflow-safe checks: remaining = sold cap - alreadySold
	if msg.Amount > releasedSupply-alreadySold {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("cannot buy %d tokens: only %d have been released", msg.Amount, releasedSupply-alreadySold))
	}
	if msg.Amount > totalSupply-alreadySold {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "not enough tokens remaining in this project")
	}

	// collect payment: cost = token price * amount (uint64 multiplication overflow check)
	cost := project.TokenPrice * msg.Amount
	if project.TokenPrice != 0 && cost/project.TokenPrice != msg.Amount {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token price overflow")
	}
	if cost == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token price must be greater than zero")
	}

	buyer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid buyer address: %s", err))
	}
	sponsor, err := sdk.AccAddressFromBech32(project.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid sponsor address: %s", err))
	}

	payment := sdk.NewCoins(sdk.NewCoin(types.PaymentDenom, math.NewIntFromUint64(cost)))

	// transfer payment from the investor to the project sponsor
	if err := k.bankKeeper.SendCoins(ctx, buyer, sponsor, payment); err != nil {
		return nil, err
	}

	index := fmt.Sprintf("%d-%s", msg.ProjectId, msg.Creator)

	holding, err := k.Holding.Get(ctx, index)
	if err != nil {
		holding = types.Holding{
			Index:     index,
			ProjectId: msg.ProjectId,
			Investor:  msg.Creator,
			Amount:    0,
			Creator:   msg.Creator,
		}
	}

	holding.Amount += msg.Amount

	if err := k.Holding.Set(ctx, index, holding); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to record holding")
	}

	return &types.MsgBuyTokensResponse{}, nil
}
