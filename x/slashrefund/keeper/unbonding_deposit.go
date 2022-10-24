package keeper

import (
	"time"

	"cosmossdk.io/math"

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
	delegatorAddress sdk.AccAddress,
	validatorAddress sdk.ValAddress,

) (val types.UnbondingDeposit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingDepositKeyPrefix))

	b := store.Get(types.UnbondingDepositKey(
		delegatorAddress.String(),
		validatorAddress.String(),
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

// SetUnbondingDepositEntry adds an entry to the unbonding deposit at
// the given addresses. It creates the unbonding deposit if it does not exist.
func (k Keeper) SetUnbondingDepositEntry(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, balance math.Int,
) types.UnbondingDeposit {
	ubd, found := k.GetUnbondingDeposit(ctx, delegatorAddr, validatorAddr)
	if found {
		ubd.AddEntry(creationHeight, minTime, balance)
	} else {
		ubd = types.NewUnbondingDeposit(delegatorAddr, validatorAddr, creationHeight, minTime, balance)
	}

	k.SetUnbondingDeposit(ctx, ubd)

	return ubd
}

// InsertUBDQueue inserts an unbonding deposit to the appropriate timeslice
// in the unbonding queue.
func (k Keeper) InsertUBDQueue(ctx sdk.Context, ubd types.UnbondingDeposit, completionTime time.Time) {

	// TODO scaffold DVPair

	dvPair := types.DVPair{DelegatorAddress: ubd.DelegatorAddress, ValidatorAddress: ubd.ValidatorAddress}

	timeSlice := k.GetUBDQueueTimeSlice(ctx, completionTime)
	if len(timeSlice) == 0 {
		k.SetUBDQueueTimeSlice(ctx, completionTime, []types.DVPair{dvPair})
	} else {
		timeSlice = append(timeSlice, dvPair)
		k.SetUBDQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}
