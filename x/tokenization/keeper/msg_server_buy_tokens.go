package keeper

import (
	"context"
	"fmt"
	"strconv"

	"paktoken/x/tokenization/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) BuyTokens(ctx context.Context, msg *types.MsgBuyTokens) (*types.MsgBuyTokensResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	project, err := k.Project.Get(ctx, msg.ProjectId)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("project %d not found", msg.ProjectId))
	}

	totalSupply, err := strconv.ParseUint(project.TotalSupply, 10, 64)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid project total supply")
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

	if alreadySold+msg.Amount > totalSupply {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "not enough tokens remaining in this project")
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
