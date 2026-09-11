package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"paktoken/x/tokenization/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) DistributeRevenue(ctx context.Context, msg *types.MsgDistributeRevenue) (*types.MsgDistributeRevenueResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "distribution amount must be greater than zero")
	}

	project, err := k.Project.Get(ctx, msg.ProjectId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("project %d not found", msg.ProjectId))
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get project")
	}

	if msg.Creator != project.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the project owner can distribute revenue")
	}

	// accumulate total held tokens using big-precision math (no uint64 overflow)
	totalHeld := math.ZeroInt()
	var holdings []types.Holding

	err = k.Holding.Walk(ctx, nil, func(idx string, h types.Holding) (bool, error) {
		if h.ProjectId == msg.ProjectId {
			totalHeld = totalHeld.Add(math.NewIntFromUint64(h.Amount))
			holdings = append(holdings, h)
		}
		return false, nil
	})
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to read holdings")
	}

	if totalHeld.IsZero() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no token holders for this project yet")
	}

	sponsor, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid sponsor address: %s", err))
	}

	revenue := math.NewIntFromUint64(msg.Amount)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	totalDistributed := math.ZeroInt()

	for _, h := range holdings {
		// share = revenue * holdingAmount / totalHeld (floor)
		share := revenue.Mul(math.NewIntFromUint64(h.Amount)).Quo(totalHeld)
		if share.IsZero() {
			continue
		}

		investor, err := sdk.AccAddressFromBech32(h.Investor)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid investor address: %s", err))
		}

		coins := sdk.NewCoins(sdk.NewCoin(types.PaymentDenom, share))
		if err := k.bankKeeper.SendCoins(ctx, sponsor, investor, coins); err != nil {
			return nil, err
		}

		totalDistributed = totalDistributed.Add(share)

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				"revenue_distributed",
				sdk.NewAttribute("project_id", fmt.Sprintf("%d", msg.ProjectId)),
				sdk.NewAttribute("investor", h.Investor),
				sdk.NewAttribute("amount", share.String()),
			),
		)
	}

	if totalDistributed.IsZero() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no distribution occurred")
	}

	return &types.MsgDistributeRevenueResponse{}, nil
}
