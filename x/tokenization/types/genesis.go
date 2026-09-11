package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		ProjectList: []Project{}, HoldingMap: []Holding{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	projectIdMap := make(map[uint64]bool)
	projectCount := gs.GetProjectCount()
	for _, elem := range gs.ProjectList {
		if _, ok := projectIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for project")
		}
		if elem.Id >= projectCount {
			return fmt.Errorf("project id should be lower or equal than the last id")
		}
		projectIdMap[elem.Id] = true
	}
	holdingIndexMap := make(map[string]struct{})

	for _, elem := range gs.HoldingMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := holdingIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for holding")
		}
		holdingIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
