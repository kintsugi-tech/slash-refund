package keeper

import (
	"encoding/binary"
	"time"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// GetUnbondingDepositCount get the total number of unbondingDeposit
func (k Keeper) GetUnbondingDepositCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.UnbondingDepositCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetUnbondingDepositCount set the total number of unbondingDeposit
func (k Keeper) SetUnbondingDepositCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.UnbondingDepositCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendUnbondingDeposit appends a unbondingDeposit in the store with a new id and update the count
func (k Keeper) AppendUnbondingDeposit(
	ctx sdk.Context,
	unbondingDeposit types.UnbondingDeposit,
) uint64 {
	// Create the unbondingDeposit
	count := k.GetUnbondingDepositCount(ctx)

	// Set the ID of the appended value
	unbondingDeposit.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingDepositKey))
	appendedValue := k.cdc.MustMarshal(&unbondingDeposit)
	store.Set(GetUnbondingDepositIDBytes(unbondingDeposit.Id), appendedValue)

	// Update unbondingDeposit count
	k.SetUnbondingDepositCount(ctx, count+1)

	return count
}

// SetUnbondingDeposit set a specific unbondingDeposit in the store
func (k Keeper) SetUnbondingDeposit(ctx sdk.Context, unbondingDeposit types.UnbondingDeposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingDepositKey))
	b := k.cdc.MustMarshal(&unbondingDeposit)
	store.Set(GetUnbondingDepositIDBytes(unbondingDeposit.Id), b)
}

// GetUnbondingDeposit returns a unbondingDeposit from its id
func (k Keeper) GetUnbondingDeposit(ctx sdk.Context, id uint64) (val types.UnbondingDeposit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingDepositKey))
	b := store.Get(GetUnbondingDepositIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUnbondingDeposit removes a unbondingDeposit from the store
func (k Keeper) RemoveUnbondingDeposit(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingDepositKey))
	store.Delete(GetUnbondingDepositIDBytes(id))
}

// GetAllUnbondingDeposit returns all unbondingDeposit
func (k Keeper) GetAllUnbondingDeposit(ctx sdk.Context) (list []types.UnbondingDeposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingDepositKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UnbondingDeposit
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetUnbondingDepositIDBytes returns the byte representation of the ID
func GetUnbondingDepositIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetUnbondingDepositIDFromBytes returns ID in uint64 format from a byte array
func GetUnbondingDepositIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

/// CUSTOM IMPLEMENTATIONS
// GetUnbondingDeposit returns a unbondingDeposit from its id
func (k Keeper) SendUnbondedTokens(ctx sdk.Context) {

	logger := k.Logger(ctx)
	logger.Error("Bella ziiiii")

	candidate_unbonding, isFound := k.GetUnbondingDeposit(ctx, uint64(0))

	if isFound {
		if candidate_unbonding.UnbondingStart.After(ctx.BlockTime().Add(120 * time.Second)) {
		}
		/*
		for i := 1; uint64(i) <= k.GetUnbondingDepositCount(ctx); i++ {

			// Se unbondati
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.Coins{msg.Amount})
			if err != nil {
				return nil, err
			}
		}
		*/
	}
}