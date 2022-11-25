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
		DepositPoolList:      []DepositPool{},
		UnbondingDepositList: []UnbondingDeposit{},
		RefundPoolList:       []RefundPool{},
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
	// Check for duplicated index in unbondingDeposit
	unbondingDepositIndexMap := make(map[string]struct{})

	for _, elem := range gs.UnbondingDepositList {
		index := string(UnbondingDepositKey(elem.DepositorAddress, elem.ValidatorAddress))
		if _, ok := unbondingDepositIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for unbondingDeposit")
		}
		unbondingDepositIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in refundPool
	refundPoolIndexMap := make(map[string]struct{})

	for _, elem := range gs.RefundPoolList {
		valOperAddr, err := sdk.ValAddressFromBech32(elem.OperatorAddress)
		if err != nil {
			panic(err)
		}
		index := string(RefundPoolKey(valOperAddr))
		if _, ok := refundPoolIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for refundPool")
		}
		refundPoolIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
