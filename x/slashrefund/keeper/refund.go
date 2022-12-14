package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// SetRefund set a specific refund in the store from its index
func (k Keeper) SetRefund(ctx sdk.Context, refund types.Refund) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundKeyPrefix))
	b := k.cdc.MustMarshal(&refund)
	store.Set(types.RefundKey(
		refund.MustGetDelegatorAddr(),
		refund.MustGetValidatorAddr(),
	), b)
}

// GetRefund returns a refund from its index
func (k Keeper) GetRefund(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (refund types.Refund, found bool) {
	moduleStore := ctx.KVStore(k.storeKey)
	store := prefix.NewStore(moduleStore, types.KeyPrefix(types.RefundKeyPrefix))
	key := types.RefundKey(delAddr, valAddr)
	b := store.Get(key)
	if b == nil {
		return refund, false
	}

	k.cdc.MustUnmarshal(b, &refund)
	return refund, true
}

// RemoveRefund removes a refund from the store
func (k Keeper) RemoveRefund(
	ctx sdk.Context,
	refund types.Refund,
) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundKeyPrefix))
	store.Delete(types.RefundKey(
		refund.MustGetDelegatorAddr(),
		refund.MustGetValidatorAddr(),
	))
}

// GetAllRefund returns all refund
func (k Keeper) GetAllRefund(ctx sdk.Context) (list []types.Refund) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Refund
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
