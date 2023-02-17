package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/slashrefund module sentinel errors
var (
	ErrNonPositiveDeposit                         = sdkerrors.Register(ModuleName, 1, "cannot deposit non-positive amount")
	ErrNonPositiveWithdraw                        = sdkerrors.Register(ModuleName, 2, "cannot withdraw non-positive amount")
	ErrNoDepositForAddress                        = sdkerrors.Register(ModuleName, 3, "no deposit for (address, validator) tuple")
	ErrNoDepositPoolForValidator                  = sdkerrors.Register(ModuleName, 4, "no deposit pool for validator")
	ErrInsufficientShares                         = sdkerrors.Register(ModuleName, 5, "insufficient shares")
	ErrInsufficientTokens                         = sdkerrors.Register(ModuleName, 6, "insufficient tokens")
	ErrZeroTokensQuotient                         = sdkerrors.Register(ModuleName, 7, "cannot divide shares by zero tokens")
	ErrDepositorShareExRateInvalid                = sdkerrors.Register(ModuleName, 8, "cannot deposit for validators with invalid (zero) ex-rate")
	ErrNotEnoughDepositShares                     = sdkerrors.Register(ModuleName, 9, "not enough deposit shares")
	ErrNoUnbondingDeposit                         = sdkerrors.Register(ModuleName, 10, "no unbonding deposit found")
	ErrNoRefundForAddress                         = sdkerrors.Register(ModuleName, 11, "no refund for (address, validator) tuple")
	ErrNoRefundPoolForValidator                   = sdkerrors.Register(ModuleName, 12, "no refund pool for validator")
	ErrNotEnoughRefundShares                      = sdkerrors.Register(ModuleName, 13, "not enough refund shares")
	ErrCantGetValidatorFromSlashEvent             = sdkerrors.Register(ModuleName, 14, "cannot get validator address from slash event")
	ErrUnknownSlashingReasonFromSlashEvent        = sdkerrors.Register(ModuleName, 15, "cannot get slashing reason from slash event")
	ErrCantGetValidatorBurnedTokensFromSlashEvent = sdkerrors.Register(ModuleName, 16, "cannot get validator burned tokens from slash event")
	ErrCantGetInfractionHeightFromSlashEvent      = sdkerrors.Register(ModuleName, 17, "cannot get infraction height from slash event")
	ErrZeroDepositAvailable                       = sdkerrors.Register(ModuleName, 18, "cannot refund: zero total deposit available for validator")
	ErrMaxUnbondingDepositEntries                 = sdkerrors.Register(ModuleName, 19, "maximum number of unbonding deposits for the pair (depositr, validator) reached")
)
