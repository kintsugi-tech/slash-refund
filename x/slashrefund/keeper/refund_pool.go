package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// SetRefundPool set a specific refundPool in the store from its index
func (k Keeper) SetRefundPool(ctx sdk.Context, refundPool types.RefundPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	b := k.cdc.MustMarshal(&refundPool)
	store.Set(types.RefundPoolKey(
		refundPool.OperatorAddress,
	), b)
}

// GetRefundPool returns a refundPool from its index
func (k Keeper) GetRefundPool(
	ctx sdk.Context,
	operatorAddress string,

) (val types.RefundPool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))

	b := store.Get(types.RefundPoolKey(
		operatorAddress,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveRefundPool removes a refundPool from the store
func (k Keeper) RemoveRefundPool(
	ctx sdk.Context,
	operatorAddress string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	store.Delete(types.RefundPoolKey(
		operatorAddress,
	))
}

// GetAllRefundPool returns all refundPool
func (k Keeper) GetAllRefundPool(ctx sdk.Context) (list []types.RefundPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RefundPool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
