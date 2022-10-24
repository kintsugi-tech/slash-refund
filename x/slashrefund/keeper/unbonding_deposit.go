package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// SetUnbondingDeposit set a specific unbondingDeposit in the store from its index
func (k Keeper) SetUnbondingDeposit(ctx sdk.Context, unbondingDeposit types.UnbondingDeposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingDepositKeyPrefix))
	b := k.cdc.MustMarshal(&unbondingDeposit)
	store.Set(types.UnbondingDepositKey(
		unbondingDeposit.DelegatorAddress,
		unbondingDeposit.ValidatorAddress,
	), b)
}

// GetUnbondingDeposit returns a unbondingDeposit from its index
func (k Keeper) GetUnbondingDeposit(
	ctx sdk.Context,
	delegatorAddress string,
	validatorAddress string,

) (val types.UnbondingDeposit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingDepositKeyPrefix))

	b := store.Get(types.UnbondingDepositKey(
		delegatorAddress,
		validatorAddress,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUnbondingDeposit removes a unbondingDeposit from the store
func (k Keeper) RemoveUnbondingDeposit(
	ctx sdk.Context,
	delegatorAddress string,
	validatorAddress string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingDepositKeyPrefix))
	store.Delete(types.UnbondingDepositKey(
		delegatorAddress,
		validatorAddress,
	))
}

// GetAllUnbondingDeposit returns all unbondingDeposit
func (k Keeper) GetAllUnbondingDeposit(ctx sdk.Context) (list []types.UnbondingDeposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingDepositKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UnbondingDeposit
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
