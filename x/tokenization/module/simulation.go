package tokenization

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"paktoken/testutil/sample"
	tokenizationsimulation "paktoken/x/tokenization/simulation"
	"paktoken/x/tokenization/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tokenizationGenesis := types.GenesisState{
		Params:      types.DefaultParams(),
		ProjectList: []types.Project{{Id: 0, Creator: sample.AccAddress()}, {Id: 1, Creator: sample.AccAddress()}}, ProjectCount: 2,
		HoldingMap: []types.Holding{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenizationGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateProject          = "op_weight_msg_tokenization"
		defaultWeightMsgCreateProject int = 100
	)

	var weightMsgCreateProject int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateProject, &weightMsgCreateProject, nil,
		func(_ *rand.Rand) {
			weightMsgCreateProject = defaultWeightMsgCreateProject
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateProject,
		tokenizationsimulation.SimulateMsgCreateProject(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateProject          = "op_weight_msg_tokenization"
		defaultWeightMsgUpdateProject int = 100
	)

	var weightMsgUpdateProject int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateProject, &weightMsgUpdateProject, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateProject = defaultWeightMsgUpdateProject
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateProject,
		tokenizationsimulation.SimulateMsgUpdateProject(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteProject          = "op_weight_msg_tokenization"
		defaultWeightMsgDeleteProject int = 100
	)

	var weightMsgDeleteProject int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteProject, &weightMsgDeleteProject, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteProject = defaultWeightMsgDeleteProject
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteProject,
		tokenizationsimulation.SimulateMsgDeleteProject(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgBuyTokens          = "op_weight_msg_tokenization"
		defaultWeightMsgBuyTokens int = 100
	)

	var weightMsgBuyTokens int
	simState.AppParams.GetOrGenerate(opWeightMsgBuyTokens, &weightMsgBuyTokens, nil,
		func(_ *rand.Rand) {
			weightMsgBuyTokens = defaultWeightMsgBuyTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBuyTokens,
		tokenizationsimulation.SimulateMsgBuyTokens(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDistributeRevenue          = "op_weight_msg_tokenization"
		defaultWeightMsgDistributeRevenue int = 100
	)

	var weightMsgDistributeRevenue int
	simState.AppParams.GetOrGenerate(opWeightMsgDistributeRevenue, &weightMsgDistributeRevenue, nil,
		func(_ *rand.Rand) {
			weightMsgDistributeRevenue = defaultWeightMsgDistributeRevenue
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDistributeRevenue,
		tokenizationsimulation.SimulateMsgDistributeRevenue(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCreateHolding          = "op_weight_msg_tokenization"
		defaultWeightMsgCreateHolding int = 100
	)

	var weightMsgCreateHolding int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateHolding, &weightMsgCreateHolding, nil,
		func(_ *rand.Rand) {
			weightMsgCreateHolding = defaultWeightMsgCreateHolding
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateHolding,
		tokenizationsimulation.SimulateMsgCreateHolding(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateHolding          = "op_weight_msg_tokenization"
		defaultWeightMsgUpdateHolding int = 100
	)

	var weightMsgUpdateHolding int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateHolding, &weightMsgUpdateHolding, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateHolding = defaultWeightMsgUpdateHolding
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateHolding,
		tokenizationsimulation.SimulateMsgUpdateHolding(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteHolding          = "op_weight_msg_tokenization"
		defaultWeightMsgDeleteHolding int = 100
	)

	var weightMsgDeleteHolding int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteHolding, &weightMsgDeleteHolding, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteHolding = defaultWeightMsgDeleteHolding
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteHolding,
		tokenizationsimulation.SimulateMsgDeleteHolding(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
