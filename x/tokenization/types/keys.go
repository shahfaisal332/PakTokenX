package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "tokenization"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"

	// PaymentDenom is the denom investors use to pay for tokens and receive
	// revenue distributions. It must match the chain's default_denom (stake).
	PaymentDenom = "stake"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_tokenization")

var (
	ProjectKey      = collections.NewPrefix("project/value/")
	ProjectCountKey = collections.NewPrefix("project/count/")
)
