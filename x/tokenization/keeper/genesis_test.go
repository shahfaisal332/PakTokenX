package keeper_test

import (
	"testing"

	"paktoken/x/tokenization/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:       types.DefaultParams(),
		ProjectList:  []types.Project{{Id: 0}, {Id: 1}},
		ProjectCount: 2,
		HoldingMap:   []types.Holding{{Index: "0"}, {Index: "1"}}}
	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.ProjectList, got.ProjectList)
	require.Equal(t, genesisState.ProjectCount, got.ProjectCount)
	require.EqualExportedValues(t, genesisState.HoldingMap, got.HoldingMap)

}
