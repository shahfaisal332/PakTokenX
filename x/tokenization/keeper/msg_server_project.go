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

func (k msgServer) CreateProject(ctx context.Context, msg *types.MsgCreateProject) (*types.MsgCreateProjectResponse, error) {
if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
}

nextId, err := k.ProjectSeq.Next(ctx)
if err != nil {
return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get next id")
}

var project = types.Project{
Id:             nextId,
Creator:        msg.Creator,
Name:           msg.Name,
TotalSupply:    msg.TotalSupply,
Owner:          msg.Owner,
TokenPrice:     msg.TokenPrice,
Category:       msg.Category,
ReleasedSupply: "0", // nothing is available to investors until the owner explicitly releases some
}

if err = k.Project.Set(ctx, nextId, project); err != nil {
return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to set project")
}

return &types.MsgCreateProjectResponse{Id: nextId}, nil
}

func (k msgServer) UpdateProject(ctx context.Context, msg *types.MsgUpdateProject) (*types.MsgUpdateProjectResponse, error) {
if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
}

val, err := k.Project.Get(ctx, msg.Id)
if err != nil {
if errors.Is(err, collections.ErrNotFound) {
return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
}
return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get project")
}

if msg.Creator != val.Creator {
return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
}

var project = types.Project{
Creator:        msg.Creator,
Id:             msg.Id,
Name:           msg.Name,
TotalSupply:    msg.TotalSupply,
Owner:          msg.Owner,
TokenPrice:     msg.TokenPrice,
Category:       msg.Category,
ReleasedSupply: val.ReleasedSupply, // preserve — UpdateProject must never silently reset this
}

if err := k.Project.Set(ctx, msg.Id, project); err != nil {
return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update project")
}

return &types.MsgUpdateProjectResponse{}, nil
}

func (k msgServer) DeleteProject(ctx context.Context, msg *types.MsgDeleteProject) (*types.MsgDeleteProjectResponse, error) {
if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
}

val, err := k.Project.Get(ctx, msg.Id)
if err != nil {
if errors.Is(err, collections.ErrNotFound) {
return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
}
return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get project")
}

if msg.Creator != val.Creator {
return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
}

if err := k.Project.Remove(ctx, msg.Id); err != nil {
return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete project")
}

return &types.MsgDeleteProjectResponse{}, nil
}
