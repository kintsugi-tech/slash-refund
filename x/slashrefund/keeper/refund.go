package keeper

import (

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleRefundsFromSlash(ctx sdk.Context, slashEvent sdk.Event) {
	// Iterate attributes to find which validators has been slashed
	logger := k.Logger(ctx)

	var valAddr sdk.ValAddress
	var refundTokens sdk.Int
	for _, attr := range slashEvent.Attributes {

		if string(attr.GetKey()) == "address" {
			valAddr, _ = sdk.ValAddressFromBech32(string(attr.GetValue()))
		}

		if string(attr.GetKey()) == "burned_coins" {
			burnedTokens, ok := math.NewIntFromString(string(attr.GetValue()))
			if !ok {
				logger.Error("Error in converting burnedTokens into int.")
			}
			refundTokens = sdk.Int(burnedTokens)
		}

	}

	depPool, isFound := k.GetDepositPool(ctx, valAddr)
	if !isFound {
		logger.Error("No pool for the slashed validator: ", valAddr.String())
	}

	// TODO: depPool Tokens has even a denom, should be managed
	if refundTokens.GT(depPool.Tokens.Amount) {
		depPool.Tokens.Amount = sdk.ZeroInt()
		// TODO: remove pool
	} else {
		depPool.Tokens.Amount = depPool.Tokens.Amount.Sub(refundTokens)
	}

	k.SetDepositPool(ctx,depPool)




	

			//deposit slashing modifico deposit
			//
			//
			//if depositCreationTime > slashingEventTime
			//	non slashare deposito
			//else
			//	slasha deposito
			////end
}