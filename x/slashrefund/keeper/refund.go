package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (keeper Keeper) GetValidatorByConsAddrBytes(ctx sdk.Context, consAddr []byte) (validator stakingtypes.Validator, found bool) {

	// Decode address
	addr, _ := sdk.ConsAddressFromBech32(string(consAddr))

	// TODO: Handle error
	return keeper.stakingKeeper.GetValidatorByConsAddr(ctx, addr)
}

func (keeper Keeper) IterateValidatorDelegations(ctx sdk.Context, consAddr []byte) {

	validator, _ := keeper.GetValidatorByConsAddrBytes(ctx, consAddr)

	// TODO: Handle not found validator
	delegations := keeper.stakingKeeper.GetValidatorDelegations(ctx, validator.GetOperator())

	for _, delegation := range delegations {
		keeper.Logger(ctx).Error("delegation", "addr", delegation.DelegatorAddress, "shares", delegation.Shares)
	}
}
