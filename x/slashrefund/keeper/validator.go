package keeper

// Update the tokens of an existing validator, update the validators power index key
/*
func (k Keeper) AddValidatorTokensAndShares(ctx sdk.Context, validator types.Validator,
	tokensToAdd sdk.Int,
) (valOut types.Validator, addedShares sdk.Dec) {
	k.DeleteValidatorByPowerIndex(ctx, validator)
	validator, addedShares = validator.AddTokensFromDel(tokensToAdd)
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)

	return validator, addedShares
}
*/