package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Update the tokens of an existing validator, update the validators power index key

func (k Keeper) AddValidatorTokensAndShares(
	ctx sdk.Context, 
	validator stakingtypes.Validator,
	tokensToAdd sdk.Coin,
) (valOut stakingtypes.Validator, addedShares sdk.Dec) {
	k.DeleteValidatorByPowerIndex(ctx, validator)
	validator, addedShares = validator.AddTokensFromDel(tokensToAdd)
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)

	return validator, addedShares
}