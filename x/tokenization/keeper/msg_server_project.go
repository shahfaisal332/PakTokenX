package keeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"paktoken/x/tokenization/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateProject(ctx context.Context, msg *types.MsgCreateProject) (*types.MsgCreateProjectResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	totalSupply, err := strconv.ParseUint(msg.TotalSupply, 10, 64)
	if err != nil || totalSupply == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total supply must be a positive integer")
	}

	if msg.TokenPrice == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token price must be greater than zero")
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

	newTotalSupply, err := strconv.ParseUint(msg.TotalSupply, 10, 64)
	if err != nil || newTotalSupply == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total supply must be a positive integer")
	}

	if msg.TokenPrice == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token price must be greater than zero")
	}

	// prevent shrinking total supply below what investors already hold
	var alreadySold uint64
	err = k.Holding.Walk(ctx, nil, func(idx string, h types.Holding) (bool, error) {
		if h.ProjectId == msg.Id {
			alreadySold += h.Amount
		}
		return false, nil
	})
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to read existing holdings")
	}

	if newTotalSupply < alreadySold {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total supply cannot be lower than the number of tokens already sold")
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
