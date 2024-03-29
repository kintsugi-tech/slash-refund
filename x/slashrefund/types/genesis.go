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
		RefundList:           []Refund{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {

	// Check for duplicated index and invalid shares amount in deposit
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
			return fmt.Errorf("duplicated index in genesis state for deposit (acc: %s / val: %s)", elem.DepositorAddress, elem.ValidatorAddress)
		}
		depositIndexMap[index] = struct{}{}

		if elem.Shares.IsNil() || !elem.Shares.IsPositive() {
			return fmt.Errorf("non-positive shares in genesis state for deposit (acc: %s / val: %s)", elem.DepositorAddress, elem.ValidatorAddress)
		}
	}

	// Check for duplicated index or invalid tokens and shares in depositPool
	depositPoolIndexMap := make(map[string]struct{})
	for _, elem := range gs.DepositPoolList {
		valOperAddr, err := sdk.ValAddressFromBech32(elem.OperatorAddress)
		if err != nil {
			panic(err)
		}
		index := string(DepositPoolKey(valOperAddr))
		if _, ok := depositPoolIndexMap[index]; ok {
			return fmt.Errorf("duplicated index in genesis state for depositPool (val: %s)", elem.OperatorAddress)
		}
		depositPoolIndexMap[index] = struct{}{}

		if elem.Tokens.IsNil() || !elem.Tokens.IsValid() || elem.Tokens.Amount.IsZero() {
			return fmt.Errorf("invalid denom or non-positive tokens amount in genesis state for depositPool (val: %s)", elem.OperatorAddress)
		}
		if elem.Shares.IsNil() || !elem.Shares.IsPositive() {
			return fmt.Errorf("non-positive shares in genesis state for depositPool (val: %s)", elem.OperatorAddress)
		}
	}

	// Check for duplicated index in unbondingDeposit
	unbondingDepositIndexMap := make(map[string]struct{})
	for _, elem := range gs.UnbondingDepositList {
		depAddr := sdk.MustAccAddressFromBech32(elem.DepositorAddress)
		valAddr, err := sdk.ValAddressFromBech32(elem.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		index := string(GetUBDKey(depAddr, valAddr))
		if _, ok := unbondingDepositIndexMap[index]; ok {
			return fmt.Errorf("duplicated index in genesis state for unbondingDeposit")
		}
		unbondingDepositIndexMap[index] = struct{}{}

		//TODO: check for max entries if implemented in withdraw logic
		for i, entry := range elem.Entries {
			if entry.InitialBalance.IsNil() || !entry.InitialBalance.IsPositive() {
				return fmt.Errorf("non-positive initial balance in genesis state for unbonding deposit entry (acc: %s / val: %s / entry: %d)", elem.DepositorAddress, elem.ValidatorAddress, i)
			}
			if entry.Balance.IsNil() || entry.Balance.IsNegative() {
				return fmt.Errorf("unset or negative balance in genesis state for unbonding deposit entry (acc: %s / val: %s / entry: %d)", elem.DepositorAddress, elem.ValidatorAddress, i)
			}
		}
	}

	// Check for duplicated index and invalid tokens and shares in refundPool
	refundPoolIndexMap := make(map[string]struct{})
	for _, elem := range gs.RefundPoolList {
		valOperAddr, err := sdk.ValAddressFromBech32(elem.OperatorAddress)
		if err != nil {
			panic(err)
		}
		index := string(RefundPoolKey(valOperAddr))
		if _, ok := refundPoolIndexMap[index]; ok {
			return fmt.Errorf("duplicated index in genesis state for refundPool (val: %s)", elem.OperatorAddress)
		}
		refundPoolIndexMap[index] = struct{}{}

		if elem.Tokens.IsNil() || !elem.Tokens.IsValid() || elem.Tokens.Amount.IsZero() {
			return fmt.Errorf("invalid denom or non-positive tokens amount in genesis state for refundPool (val: %s)", elem.OperatorAddress)
		}
		if elem.Shares.IsNil() || !elem.Shares.IsPositive() {
			return fmt.Errorf("non-positive shares in genesis state for refundPool (val: %s)", elem.OperatorAddress)
		}
	}

	// Check for duplicated index and invalid shares in refund
	refundIndexMap := make(map[string]struct{})
	for _, elem := range gs.RefundList {
		delegator, err := sdk.AccAddressFromBech32(elem.DelegatorAddress)
		if err != nil {
			return err
		}
		validator, err := sdk.ValAddressFromBech32(elem.ValidatorAddress)
		if err != nil {
			return err
		}
		index := string(RefundKey(delegator, validator))
		if _, ok := refundIndexMap[index]; ok {
			return fmt.Errorf("duplicated index in genesis state for refund (acc: %s / val: %s)", elem.DelegatorAddress, elem.ValidatorAddress)
		}
		refundIndexMap[index] = struct{}{}

		if elem.Shares.IsNil() || !elem.Shares.IsPositive() {
			return fmt.Errorf("non-positive shares in genesis state for refund (acc: %s / val: %s)", elem.DelegatorAddress, elem.ValidatorAddress)
		}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
