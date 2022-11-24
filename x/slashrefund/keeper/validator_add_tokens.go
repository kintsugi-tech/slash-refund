package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) AddValidatorTokens_SR(
	ctx sdk.Context,
	validator stakingtypes.Validator,
	tokensToAdd sdk.Int,
) stakingtypes.Validator {
	k.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Tokens = validator.Tokens.Add(tokensToAdd)
	k.stakingKeeper.SetValidator(ctx, validator)
	k.stakingKeeper.SetValidatorByPowerIndex(ctx, validator)

	return validator
}
