package keeper

import (
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// Withdraw implements the state transition logic associated to a valid amount of tokens that a user
// wants to withdraw from the module.
func (k Keeper) Withdraw(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	tokens sdk.Coin,
) (sdk.Coin, time.Time, error) {

	deposit, found := k.GetDeposit(ctx, depAddr, valAddr)
	if !found {
		return sdk.NewCoin(tokens.Denom, sdk.NewInt(0)), time.Time{}, types.ErrNoDepositForAddress
	}

	// Get the deposit pool. If at this point it can't be found, then panic is called.
	// This is done because a deposit pool is removed from the store alongside with
	// its deposits when the deposit pool is completely used during refund generation.
	// When a withdraw is done, the deposit pool and the deposit are updated and
	// removed if empty. A situation in which a deposit is found but no linked deposit
	// pool can be found is the result of a serious malfunction thus panic is called.
	depPool, found := k.GetDepositPool(ctx, valAddr)
	if !found {
		panic("found deposit but not the deposit pool")
	}

	// Check if requested amount is valid and returns associated shares.
	witShares, err := k.ComputeAssociatedShares(ctx, deposit, depPool, tokens)
	if err != nil {
		return sdk.NewCoin(tokens.Denom, sdk.NewInt(0)), time.Time{}, err
	}

	if k.HasMaxUnbondingDepositEntries(ctx, depAddr, valAddr) {
		return sdk.NewCoin(tokens.Denom, sdk.NewInt(0)), time.Time{}, types.ErrMaxUnbondingDepositEntries
	}

	witAmt, err := k.Unbond(ctx, deposit, depPool, valAddr, witShares)
	if err != nil {
		return sdk.NewCoin(tokens.Denom, sdk.NewInt(0)), time.Time{}, err
	}

	// Compute time at which withdrawn tokens will complete the unbonding.
	completionTime := ctx.BlockHeader().Time.Add(k.stakingKeeper.UnbondingTime(ctx))

	ubd := k.SetUnbondingDepositEntry(ctx, depAddr, valAddr, ctx.BlockHeight(), completionTime, witAmt)

	k.InsertUBDQueue(ctx, ubd, completionTime)

	return sdk.NewCoin(k.AllowedTokens(ctx)[0], witAmt), completionTime, nil
}

// Returns user's shares associated with desired withdrawal tokens if available,
// or an error.
func (k Keeper) ComputeAssociatedShares(
	ctx sdk.Context,
	deposit types.Deposit,
	depPool types.DepositPool,
	tokens sdk.Coin,
) (shares sdk.Dec, err error) {

	// Compute shares from desired withdrawal tokens.
	shares, err = depPool.SharesFromTokens(tokens)
	if err != nil {
		return sdk.NewDec(0), err
	}

	// Compute rounded down shares from desired withdrawal tokens.
	sharesTruncated, err := depPool.SharesFromTokensTruncated(tokens)
	if err != nil {
		return sdk.NewDec(0), err
	}

	// Check if desired withdrawal tokens converted to truncated shares are greater than actual
	// total of depositor shares.
	depositorShares := deposit.GetShares()
	if sharesTruncated.GT(depositorShares) {
		return sdk.NewDec(0), sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token amount")
	}

	// Cap shares (not-truncated) at total depositor shares.
	if shares.GT(depositorShares) {
		shares = depositorShares
	}

	// Ensure that the pool has enough shares to be removed.
	if depPool.Shares.LT(shares) {
		return sdk.NewDec(0), sdkerrors.Wrap(types.ErrNotEnoughDepositShares, deposit.Shares.String())
	}

	return shares, nil
}

// -------------------------------------------------------------------------------------------------
// Unbond
// -------------------------------------------------------------------------------------------------

func (k Keeper) Unbond(
	ctx sdk.Context,
	deposit types.Deposit,
	depPool types.DepositPool,
	valAddr sdk.ValAddress,
	shares sdk.Dec,
) (issuedTokensAmt sdk.Int, err error) {

	// Subtract shares from deposit.
	deposit.Shares = deposit.Shares.Sub(shares)

	// Remove the deposit if zero or set a new deposit.
	if deposit.Shares.IsZero() {
		k.RemoveDeposit(ctx, deposit)
	} else {
		k.SetDeposit(ctx, deposit)
	}

	depPool, issuedTokensAmt = k.RemoveDepPoolTokensAndShares(ctx, depPool, shares)

	if depPool.Shares.IsZero() {
		k.RemoveDepositPool(ctx, valAddr)
	}

	return issuedTokensAmt, nil
}

// Checks if a user has already requested the maximum number of allowed unbonding deposits for a
// specific validator in the considered timeframe.
func (k Keeper) HasMaxUnbondingDepositEntries(ctx sdk.Context, depAddr sdk.AccAddress, valAddr sdk.ValAddress) bool {
	ubd, found := k.GetUnbondingDeposit(ctx, depAddr, valAddr)
	if !found {
		return false
	}

	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
}

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

// Returns all unbonding deposits associated to a particular validator.
func (k Keeper) GetUnbondingDepositsFromValidator(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
) (ubds []types.UnbondingDeposit) {

	store := ctx.KVStore(k.storeKey)
	keyspace := types.GetUBDsByValIndexKey(valAddr)
	iterator := sdk.KVStorePrefixIterator(store, keyspace)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {

		// reset ubd
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

// InsertUBDQueue inserts an unbonding deposit to the appropriate timeslice in the unbonding queue.
func (k Keeper) InsertUBDQueue(
	ctx sdk.Context,
	ubd types.UnbondingDeposit,
	completionTime time.Time,
) {

	// dvPair indicates the pair of delegator and validator
	dvPair := types.DVPair{
		DepositorAddress: ubd.DepositorAddress,
		ValidatorAddress: ubd.ValidatorAddress,
	}

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

// Returns all unbonding queue timeslices from time 0 until endTime.
func (k Keeper) UBDQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(
		types.UnbondingQueueKey,
		sdk.InclusiveEndBytes(types.GetUnbondingDepositTimeKey(endTime)),
	)
}

// Returns a concatenated list of all timeslices inclusively previous to currTime, and deletes the
// timeslices from the queue.
func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []types.DVPair) {
	store := ctx.KVStore(k.storeKey)

	// Gets an iterator for all timeslices from time 0 until the current Blockheader time
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

// Completes the unbonding of all mature entries in the
// retrieved unbonding deposit object and returns the total unbonding balance
// or an error upon failure.
func (k Keeper) CompleteUnbonding(ctx sdk.Context, currTime time.Time, depAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	ubd, found := k.GetUnbondingDeposit(ctx, depAddr, valAddr)
	if !found {
		return nil, types.ErrNoUnbondingDeposit
	}

	//TODO: generalize refundDenom with all the AllowedTokens
	refundDenom := k.AllowedTokens(ctx)[0]
	balances := sdk.NewCoins()

	// Loop through all entries and complete mature unbonding entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(currTime) {
			ubd.RemoveEntry(int64(i))
			i--

			// Track withdraw only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewCoin(refundDenom, entry.Balance)
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(
					ctx, types.ModuleName, depAddr, sdk.NewCoins(amt),
				); err != nil {
					return nil, err
				}

				balances = balances.Add(amt)
			}
		}
	}

	// Set the unbonding deposit or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingDeposit(ctx, ubd)
	} else {
		k.SetUnbondingDeposit(ctx, ubd)
	}

	return balances, nil
}