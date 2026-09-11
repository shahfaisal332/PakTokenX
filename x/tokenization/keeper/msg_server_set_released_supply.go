package keeper

import (
"context"
"fmt"
"strconv"

"paktoken/x/tokenization/types"

errorsmod "cosmossdk.io/errors"
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

newReleased, err := strconv.ParseUint(msg.ReleasedSupply, 10, 64)
if err != nil {
return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid released supply value")
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

project.ReleasedSupply = msg.ReleasedSupply

if err := k.Project.Set(ctx, msg.ProjectId, project); err != nil {
return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update released supply")
}

return &types.MsgSetReleasedSupplyResponse{}, nil
}
