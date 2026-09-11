package keeper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"paktoken/x/tokenization/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) SetReleasedSupply(ctx context.Context, msg *types.MsgSetReleasedSupply) (*types.MsgSetReleasedSupplyResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	project, err := k.Project.Get(ctx, msg.ProjectId)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("project %d not found", msg.ProjectId))
	}

	if msg.Creator != project.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the project owner can change released supply")
	}

	totalSupply, err := strconv.ParseUint(project.TotalSupply, 10, 64)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid project total supply")
	}

	// Accepts either an absolute amount ("400") or a percentage ("2%").
	// A percentage is computed against totalSupply at issue time, so the
	// sponsor can issue e.g. "2%" during tokenization or after creating the project.
	var newReleased uint64
	if strings.HasSuffix(msg.ReleasedSupply, "%") {
		pctStr := strings.TrimSuffix(msg.ReleasedSupply, "%")
		pct, err := strconv.ParseUint(pctStr, 10, 64)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid percentage value, expected e.g. \"2%\"")
		}
		if pct > 100 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "percentage cannot exceed 100")
		}

		p := math.NewIntFromUint64(pct)
		t := math.NewIntFromUint64(totalSupply)
		newReleased = p.Mul(t).Quo(math.NewInt(100)).Uint64()
	} else {
		newReleased, err = strconv.ParseUint(msg.ReleasedSupply, 10, 64)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid released supply value")
		}
	}

	if newReleased > totalSupply {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "released supply cannot exceed total supply")
	}

	// never allow releasing below what's already been sold to investors
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

	if newReleased < alreadySold {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("cannot set released supply below %d tokens already sold", alreadySold))
	}

	project.ReleasedSupply = strconv.FormatUint(newReleased, 10)

	if err := k.Project.Set(ctx, msg.ProjectId, project); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update released supply")
	}

	return &types.MsgSetReleasedSupplyResponse{}, nil
}
