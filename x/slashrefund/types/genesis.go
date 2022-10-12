package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DepositList:          []Deposit{},
		UnbondingDepositList: []UnbondingDeposit{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in deposit
	depositIndexMap := make(map[string]struct{})

	for _, elem := range gs.DepositList {
		index := string(DepositKey(elem.Address, elem.ValidatorAddress))
		if _, ok := depositIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for deposit")
		}
		depositIndexMap[index] = struct{}{}
	}
	// Check for duplicated ID in unbondingDeposit
	unbondingDepositIdMap := make(map[uint64]bool)
	unbondingDepositCount := gs.GetUnbondingDepositCount()
	for _, elem := range gs.UnbondingDepositList {
		if _, ok := unbondingDepositIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for unbondingDeposit")
		}
		if elem.Id >= unbondingDepositCount {
			return fmt.Errorf("unbondingDeposit id should be lower or equal than the last id")
		}
		unbondingDepositIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
