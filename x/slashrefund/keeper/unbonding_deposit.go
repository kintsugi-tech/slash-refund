package keeper

import (
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// Sets a specific unbondingDeposit in the store from its index.
func (k Keeper) SetUnbondingDeposit(ctx sdk.Context, unbondingDeposit types.UnbondingDeposit) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&unbondingDeposit)
	depAddr := sdk.MustAccAddressFromBech32(unbondingDeposit.DepositorAddress)
	valAddr, err := sdk.ValAddressFromBech32(unbondingDeposit.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetUBDKey(depAddr, valAddr)
	store.Set(key, b)

	// Store also reverse order for indexing purpose.
	// NOTE: key2 is validator-depositor;
	//       none value will be stored with this key because this key is used for key-rearrangement only,
	//       in order to obtain from this key the actual depositor-validator key.
	key2 := types.GetUBDByValIndexKey(valAddr, depAddr)
	store.Set(key2, []byte{})
}

// Returns an unbondingDeposit from its index
func (k Keeper) GetUnbondingDeposit(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	valAddr sdk.ValAddress,

) (ubd types.UnbondingDeposit, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(depAddr, valAddr)
	value := store.Get(key)

	if value == nil {
		return ubd, false
	}

	k.cdc.MustUnmarshal(value, &ubd)
	return ubd, true
}

// RemoveUnbondingDeposit removes a unbondingDeposit from the store
func (k Keeper) RemoveUnbondingDeposit(
	ctx sdk.Context,
	ubd types.UnbondingDeposit,
) {
	depAddr := sdk.MustAccAddressFromBech32(ubd.DepositorAddress)
	valAddr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	store := ctx.KVStore(k.storeKey)

	key := types.GetUBDKey(depAddr, valAddr)
	store.Delete(key)

	key2 := types.GetUBDByValIndexKey(valAddr, depAddr)
	store.Delete(key2)
}

// GetAllUnbondingDeposit returns all unbondingDeposit
func (k Keeper) GetAllUnbondingDeposit(ctx sdk.Context) (ubds []types.UnbondingDeposit) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetUBDsKeyPrefix())
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var ubd types.UnbondingDeposit
		k.cdc.MustUnmarshal(iterator.Value(), &ubd)
		ubds = append(ubds, ubd)
	}
	return ubds
}

// GetUnbondingDeposit returns a unbondingDeposit from its index
func (k Keeper) GetUnbondingDepositByValIndexKey(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	depAddr sdk.AccAddress,
) (ubd types.UnbondingDeposit, found bool) {

	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDByValIndexKey(valAddr, depAddr)
	b := store.Get(key)
	if b == nil {
		return ubd, false
	}

	k.cdc.MustUnmarshal(b, &ubd)
	return ubd, true
}

// GetUnbondingDelegationsFromValidator returns all unbonding delegations from a
// particular validator.
func (k Keeper) GetUnbondingDepositsFromValidator(
	ctx sdk.Context, 
	valAddr sdk.ValAddress,
	) (ubds []types.UnbondingDeposit) {

	store := ctx.KVStore(k.storeKey)
	keyspace := types.GetUBDsByValIndexKey(valAddr)
	iterator := sdk.KVStorePrefixIterator(store, keyspace)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {

		//reset ubd
		var ubd types.UnbondingDeposit

		// rearrange key
		key2 := iterator.Key()
		key := types.GetUBDKeyFromValIndexKey(key2)

		// access store and append value
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &ubd)
		ubds = append(ubds, ubd)
	}

	return ubds
}

// Adds an entry to the unbonding deposit at the given addresses.
// It creates the unbonding deposit if it does not exist.
func (k Keeper) SetUnbondingDepositEntry(
	ctx sdk.Context, 
	depositorAddr sdk.AccAddress, 
	validatorAddr sdk.ValAddress,
	creationHeight int64, 
	minTime time.Time, 
	balance math.Int,
) types.UnbondingDeposit {
	ubd, found := k.GetUnbondingDeposit(ctx, depositorAddr, validatorAddr)
	if found {
		ubd.AddEntry(creationHeight, minTime, balance)
	} else {
		ubd = types.NewUnbondingDeposit(depositorAddr, validatorAddr, creationHeight, minTime, balance)
	}

	k.SetUnbondingDeposit(ctx, ubd)

	return ubd
}

// InsertUBDQueue inserts an unbonding deposit to the appropriate timeslice
// in the unbonding queue.
func (k Keeper) InsertUBDQueue(ctx sdk.Context, ubd types.UnbondingDeposit, completionTime time.Time) {

	// dvPair indicates the pair of delegator and validator
	dvPair := types.DVPair{DepositorAddress: ubd.DepositorAddress, ValidatorAddress: ubd.ValidatorAddress}

	// timeSlice is a slice of dvPair elements, linked to a given unbonding completionTime
	timeSlice := k.GetUBDQueueTimeSlice(ctx, completionTime)

	// append the new dvPair to the timeSlice and set the udated timeSlice in the unbonding queue
	if len(timeSlice) == 0 {
		k.SetUBDQueueTimeSlice(ctx, completionTime, []types.DVPair{dvPair})
	} else {
		timeSlice = append(timeSlice, dvPair)
		k.SetUBDQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

func (k Keeper) GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []types.DVPair) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetUnbondingDepositTimeKey(timestamp))
	if bz == nil {
		return []types.DVPair{}
	}

	pairs := types.DVPairs{}
	k.cdc.MustUnmarshal(bz, &pairs)

	return pairs.Pairs
}

// SetUBDQueueTimeSlice sets a specific unbonding queue timeslice.
func (k Keeper) SetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time, dvpairs []types.DVPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&types.DVPairs{Pairs: dvpairs})
	store.Set(types.GetUnbondingDepositTimeKey(timestamp), bz)
}

// UBDQueueIterator returns all the unbonding queue timeslices from time 0 until endTime.
func (k Keeper) UBDQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnbondingQueueKey,
		sdk.InclusiveEndBytes(types.GetUnbondingDepositTimeKey(endTime)))
}

// DequeueAllMatureUBDQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue.
func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []types.DVPair) {
	store := ctx.KVStore(k.storeKey)

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UBDQueueIterator(ctx, currTime)
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := types.DVPairs{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		matureUnbonds = append(matureUnbonds, timeslice.Pairs...)

		store.Delete(unbondingTimesliceIterator.Key())
	}

	return matureUnbonds
}

// CompleteUnbonding completes the unbonding of all mature entries in the
// retrieved unbonding deposit object and returns the total unbonding balance
// or an error upon failure.
func (k Keeper) CompleteUnbonding(ctx sdk.Context, depAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	ubd, found := k.GetUnbondingDeposit(ctx, depAddr, valAddr)
	if !found {
		return nil, types.ErrNoUnbondingDeposit
	}

	//TODO: generalize refundDenom with all the AllowedTokens
	refundDenom := k.AllowedTokens(ctx)[0]
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	depositorAddress, err := sdk.AccAddressFromBech32(ubd.DepositorAddress)
	if err != nil {
		return nil, err
	}

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			i--

			// track withdraw only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewCoin(refundDenom, entry.Balance)
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(
					ctx, types.ModuleName, depositorAddress, sdk.NewCoins(amt),
				); err != nil {
					return nil, err
				}

				balances = balances.Add(amt)
			}
		}
	}

	// set the unbonding deposit or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingDeposit(ctx, ubd)
	} else {
		k.SetUnbondingDeposit(ctx, ubd)
	}

	return balances, nil
}
