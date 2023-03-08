package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected methods for the account keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected methods for the bank module.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(
		ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins,
	) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// StakingKeeper defines the expected methods for the staking module.
type StakingKeeper interface {
	GetValidatorDelegations(ctx sdk.Context, valAddr sdk.ValAddress) (delegations []stakingtypes.Delegation)
	GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, found bool)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
	UnbondingTime(ctx sdk.Context) (res time.Duration)
	GetUnbondingDelegationsFromValidator(ctx sdk.Context, valAddr sdk.ValAddress) (ubds []stakingtypes.UnbondingDelegation)
	GetRedelegationsFromSrcValidator(ctx sdk.Context, valAddr sdk.ValAddress) (reds []stakingtypes.Redelegation)
	//
	SetValidator(ctx sdk.Context, validator stakingtypes.Validator)
	SetValidatorByPowerIndex(ctx sdk.Context, validator stakingtypes.Validator)
	DeleteValidatorByPowerIndex(ctx sdk.Context, validator stakingtypes.Validator)
	//
	SetUnbondingDelegation(ctx sdk.Context, ubd stakingtypes.UnbondingDelegation)
}

// SlashKeeper defines the expected methods for the slash module.
type SlashingKeeper interface {
	SlashFractionDoubleSign(ctx sdk.Context) (res sdk.Dec)
	SlashFractionDowntime(ctx sdk.Context) (res sdk.Dec)
}
