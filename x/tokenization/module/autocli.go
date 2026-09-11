package tokenization

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"paktoken/x/tokenization/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListProject",
					Use:       "list-project",
					Short:     "List all project",
				},
				{
					RpcMethod:      "GetProject",
					Use:            "get-project [id]",
					Short:          "Gets a project by id",
					Alias:          []string{"show-project"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "ListHolding",
					Use:       "list-holding",
					Short:     "List all holding",
				},
				{
					RpcMethod:      "GetHolding",
					Use:            "get-holding [id]",
					Short:          "Gets a holding",
					Alias:          []string{"show-holding"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateProject",
					Use:            "create-project [name] [total-supply] [owner] [token-price]",
					Short:          "Create project",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "name"}, {ProtoField: "total_supply"}, {ProtoField: "owner"}, {ProtoField: "token_price"}},
				},
				{
					RpcMethod:      "UpdateProject",
					Use:            "update-project [id] [name] [total-supply] [owner] [token-price]",
					Short:          "Update project",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "name"}, {ProtoField: "total_supply"}, {ProtoField: "owner"}, {ProtoField: "token_price"}},
				},
				{
					RpcMethod:      "DeleteProject",
					Use:            "delete-project [id]",
					Short:          "Delete project",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "BuyTokens",
					Use:            "buy-tokens [project-id] [amount]",
					Short:          "Send a buy-tokens tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "DistributeRevenue",
					Use:            "distribute-revenue [project-id] [amount]",
					Short:          "Send a distribute-revenue tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "CreateHolding",
					Use:            "create-holding [index] [amount] [project-id] [investor]",
					Short:          "Create a new holding",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "amount"}, {ProtoField: "project_id"}, {ProtoField: "investor"}},
				},
				{
					RpcMethod:      "UpdateHolding",
					Use:            "update-holding [index] [amount] [project-id] [investor]",
					Short:          "Update holding",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "amount"}, {ProtoField: "project_id"}, {ProtoField: "investor"}},
				},
				{
					RpcMethod:      "DeleteHolding",
					Use:            "delete-holding [index]",
					Short:          "Delete holding",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
			},
		},
	}
}
