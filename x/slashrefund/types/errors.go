package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/slashrefund module sentinel errors
var (
	ErrZeroDeposit                                = sdkerrors.Register(ModuleName, 11, "cannot deposit zero amount")
	ErrZeroWithdraw                               = sdkerrors.Register(ModuleName, 12, "cannot withdraw zero amount")
	ErrZeroClaim                                  = sdkerrors.Register(ModuleName, 13, "cannot claim zero amount")
	ErrNoDepositForAddress                        = sdkerrors.Register(ModuleName, 14, "no deposit for (address, validator) tuple")
	ErrNoDepositPoolForValidator                  = sdkerrors.Register(ModuleName, 15, "no deposit pool for validator")
	ErrInsufficientShares                         = sdkerrors.Register(ModuleName, 16, "insufficient shares")
	ErrInsufficientTokens                         = sdkerrors.Register(ModuleName, 17, "insufficient tokens")
	ErrZeroTokensQuotient                         = sdkerrors.Register(ModuleName, 18, "cannot divide shares by zero tokens")
	ErrDepositorShareExRateInvalid                = sdkerrors.Register(ModuleName, 19, "cannot deposit for validators with invalid (zero) ex-rate")
	ErrNotEnoughDepositShares                     = sdkerrors.Register(ModuleName, 20, "not enough deposit shares")
	ErrNoUnbondingDeposit                         = sdkerrors.Register(ModuleName, 21, "no unbonding deposit found")
	ErrNoRefundForAddress                         = sdkerrors.Register(ModuleName, 22, "no refund for (address, validator) tuple")
	ErrNoRefundPoolForValidator                   = sdkerrors.Register(ModuleName, 23, "no refund pool for validator")
	ErrNotEnoughRefundShares                      = sdkerrors.Register(ModuleName, 24, "not enough refund shares")
	ErrCantGetValidatorFromSlashEvent             = sdkerrors.Register(ModuleName, 25, "cannot get validator address from slash event")
	ErrUnknownSlashingReasonFromSlashEvent        = sdkerrors.Register(ModuleName, 26, "cannot get slashing reason from slash event")
	ErrCantGetValidatorBurnedTokensFromSlashEvent = sdkerrors.Register(ModuleName, 27, "cannot get validator burned tokens from slash event")
	ErrCantGetInfractionHeightFromSlashEvent      = sdkerrors.Register(ModuleName, 28, "cannot get infraction height from slash event")
	ErrZeroDepositAvailable                       = sdkerrors.Register(ModuleName, 29, "cannot refund: zero total deposit available for validator")
)
