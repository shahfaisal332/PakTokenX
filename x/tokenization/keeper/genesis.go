package keeper

import (
	"context"

	"paktoken/x/tokenization/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.ProjectList {
		if err := k.Project.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.ProjectSeq.Set(ctx, genState.ProjectCount); err != nil {
		return err
	}
	for _, elem := range genState.HoldingMap {
		if err := k.Holding.Set(ctx, elem.Index, elem); err != nil {
			return err
		}
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = k.Project.Walk(ctx, nil, func(key uint64, elem types.Project) (bool, error) {
		genesis.ProjectList = append(genesis.ProjectList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.ProjectCount, err = k.ProjectSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}
	if err := k.Holding.Walk(ctx, nil, func(_ string, val types.Holding) (stop bool, err error) {
		genesis.HoldingMap = append(genesis.HoldingMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return genesis, nil
}
