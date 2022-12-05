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
	ErrInsufficientShares          = sdkerrors.Register(ModuleName, 21, "insufficient shares")
	ErrInsufficientTokens          = sdkerrors.Register(ModuleName, 22, "insufficient tokens")
	ErrZeroTokensQuotient          = sdkerrors.Register(ModuleName, 23, "cannot divide shares by zero tokens")
	ErrDepositorShareExRateInvalid = sdkerrors.Register(ModuleName, 2, "cannot deposit for validators with invalid (zero) ex-rate")
	ErrNotEnoughDepositShares      = sdkerrors.Register(ModuleName, 24, "not enough deposit shares")
	ErrNoUnbondingDeposit          = sdkerrors.Register(ModuleName, 26, "no unbonding deposit found")
	ErrNoRefundForAddress          = sdkerrors.Register(ModuleName, 1400, "no refund for (address, validator) tuple")
	ErrNoRefundPoolForValidator    = sdkerrors.Register(ModuleName, 1500, "no refund pool for validator")
	ErrNotEnoughRefundShares       = sdkerrors.Register(ModuleName, 1600, "not enough refund shares")
)
