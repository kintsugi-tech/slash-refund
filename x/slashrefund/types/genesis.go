package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DepositList:          []Deposit{},
		UnbondingDepositList: []UnbondingDeposit{},
		DepositPoolList:      []DepositPool{},
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

		depositor, err := sdk.AccAddressFromBech32(elem.DepositorAddress)
		if err != nil {
			return err
		}

		validator, err := sdk.ValAddressFromBech32(elem.ValidatorAddress)
		if err != nil {
			return err
		}
		index := string(DepositKey(depositor, validator))
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
	// Check for duplicated index in depositPool
	depositPoolIndexMap := make(map[string]struct{})

	for _, elem := range gs.DepositPoolList {
		valOperAddr, err := sdk.ValAddressFromBech32(elem.OperatorAddress)
		if err != nil {
			panic(err)
		}
		index := string(DepositPoolKey(valOperAddr))
		if _, ok := depositPoolIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for depositPool")
		}
		depositPoolIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
