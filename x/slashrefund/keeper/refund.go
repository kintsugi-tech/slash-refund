package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) HandleRefundsFromSlash(ctx sdk.Context, slashEvent sdk.Event) sdk.Int {
	// Iterate attributes to find which validators has been slashed
	logger := k.Logger(ctx)

	var validator stakingtypes.Validator
	var refundTokens sdk.Int

	for _, attr := range slashEvent.Attributes {

		if string(attr.GetKey()) == "address" {
			var isFound bool
			validator, isFound = k.GetValidatorByConsAddrBytes(ctx, attr.GetValue())
			if !isFound {
				logger.Error("ERROR cant find validator.")
				return sdk.NewInt(0)
			}
		}

		if string(attr.GetKey()) == "burned_coins" {
			burnedTokens, ok := math.NewIntFromString(string(attr.GetValue()))
			if !ok {
				logger.Error("Error in converting burnedTokens into int.")
			}
			refundTokens = sdk.Int(burnedTokens)
		}
	}

	valAddr := validator.GetOperator()
	depPool, isFound := k.GetDepositPool(ctx, valAddr)
	if !isFound {
		logger.Error("No pool for the slashed validator: ", valAddr.String())
		return sdk.NewInt(0)
	}

	// TODO: depPool Tokens has even a denom, should be managed
	if refundTokens.GT(depPool.Tokens.Amount) {
		depPool.Tokens.Amount = sdk.ZeroInt()
		// TODO: remove pool
	} else {
		depPool.Tokens.Amount = depPool.Tokens.Amount.Sub(refundTokens)
	}
	k.SetDepositPool(ctx, depPool)
	k.stakingKeeper.AddValidatorTokens(ctx, validator, refundTokens)

	//deposit slashing modifico deposit
	//
	//
	//if depositCreationTime > slashingEventTime
	//	non slashare deposito
	//else
	//	slasha deposito
	////end

	return refundTokens
}

func (keeper Keeper) GetValidatorByConsAddrBytes(ctx sdk.Context, consAddr []byte) (validator stakingtypes.Validator, found bool) {

	// Decode address
	addr, _ := sdk.ConsAddressFromBech32(string(consAddr))

	// TODO: Handle error
	return keeper.stakingKeeper.GetValidatorByConsAddr(ctx, addr)
}
