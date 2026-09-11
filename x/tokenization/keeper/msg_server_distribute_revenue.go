package keeper

import (
	"context"
	"fmt"

	"paktoken/x/tokenization/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) DistributeRevenue(ctx context.Context, msg *types.MsgDistributeRevenue) (*types.MsgDistributeRevenueResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	project, err := k.Project.Get(ctx, msg.ProjectId)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("project %d not found", msg.ProjectId))
	}

	if msg.Creator != project.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the project owner can distribute revenue")
	}

	var totalHeld uint64
	holdings := []types.Holding{}

	err = k.Holding.Walk(ctx, nil, func(idx string, h types.Holding) (bool, error) {
		if h.ProjectId == msg.ProjectId {
			totalHeld += h.Amount
			holdings = append(holdings, h)
		}
		return false, nil
	})
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to read holdings")
	}

	if totalHeld == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no token holders for this project yet")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	for _, h := range holdings {
		share := (msg.Amount * h.Amount) / totalHeld

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				"revenue_distributed",
				sdk.NewAttribute("project_id", fmt.Sprintf("%d", msg.ProjectId)),
				sdk.NewAttribute("investor", h.Investor),
				sdk.NewAttribute("amount", fmt.Sprintf("%d", share)),
			),
		)
	}

	return &types.MsgDistributeRevenueResponse{}, nil
}
