package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/slashrefund module sentinel errors
var (
	ErrNoDepositForAddress         = sdkerrors.Register(ModuleName, 1100, "no deposit for (address, validator) tuple")
	ErrNoDepositPoolForValidator   = sdkerrors.Register(ModuleName, 1200, "no deposit pool for validator")
	ErrZeroWithdraw                = sdkerrors.Register(ModuleName, 1300, "cannot withdraw zero amount")
	ErrInsufficientShares          = sdkerrors.Register(ModuleName, 22, "insufficient deposit shares")
	ErrDepositorShareExRateInvalid = sdkerrors.Register(ModuleName, 2, "cannot deposit for validators with invalid (zero) ex-rate")
	ErrNotEnoughDepositShares      = sdkerrors.Register(ModuleName, 24, "not enough deposit shares")
	ErrNoUnbondingDeposit          = sdkerrors.Register(ModuleName, 26, "no unbonding deposit found")
)
