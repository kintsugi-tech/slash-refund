package keeper

import (
	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// -------------------------------------------------------------------------------------------------
// Deposit
// -------------------------------------------------------------------------------------------------

// Deposit implements the state transition logic for a deposit. It checks if a pool for the
// validator exists or create it. Sended tokens are deposited into the module and shares are
// created for the depositor to take into account its balance inside the module.
func (k Keeper) Deposit(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	depCoin sdk.Coin,
	validator stakingtypes.Validator,
) (newShares sdk.Dec, err error) {

	// Check if a validator has zero token but shares. This situation can arise due to slashing
	// of the considered validator. Deposit for these validators aren't allowed.
	if validator.InvalidExRate() {
		// Return zero shares and an error
		return sdk.ZeroDec(), types.ErrDepositorShareExRateInvalid
	}

	valAddr := validator.GetOperator()

	// Check if the deposit pool exists or create it.
	var deposit types.Deposit
	depPool, found := k.GetDepositPool(ctx, valAddr)
	if !found {
		depPool = types.NewDepositPool(
			valAddr,
			sdk.NewCoin(depCoin.Denom, sdk.ZeroInt()),
			sdk.ZeroDec(),
		)

		// If the pool does not exists no deposit can exists.
		deposit = types.NewDeposit(depAddr, valAddr, sdk.ZeroDec())
	} else {
		// Check if the deposit exists or create it.
		deposit, found = k.GetDeposit(ctx, depAddr, valAddr)
		if !found {
			// If a previous deposit does not exist initialize one with zero shares.
			deposit = types.NewDeposit(depAddr, valAddr, sdk.ZeroDec())
		}
	}

	// Send deposited tokens to the slashrefund module.
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, 
		depAddr, 
		types.ModuleName, 
		sdk.NewCoins(depCoin),
	); err != nil {
		return sdk.Dec{}, err
	}

	// Deposited tokens are treated as pool shares, similarly to the staking module.
	newShares = k.AddDepPoolTokensAndShares(ctx, depPool, depCoin)

	deposit.Shares = deposit.Shares.Add(newShares)
	k.SetDeposit(ctx, deposit)

	return sdk.NewDec(depCoin.Amount.Int64()), nil
}

// GetDeposit returns a deposit from its indices: depAddr & valAddr
func (k Keeper) GetDeposit(
	ctx sdk.Context, 
	depAddr sdk.AccAddress, 
	valAddr sdk.ValAddress,
) (deposit types.Deposit, found bool) {

	moduleStore := ctx.KVStore(k.storeKey)
	store := prefix.NewStore(moduleStore, types.KeyPrefix(types.DepositKeyPrefix))
	key := types.DepositKey(depAddr, valAddr)
	b := store.Get(key)
	if b == nil {
		return deposit, false
	}

	k.cdc.MustUnmarshal(b, &deposit)
	return deposit, true
}

// SetDeposit set a specific deposit in the store from its index
func (k Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	b := k.cdc.MustMarshal(&deposit)
	store.Set(types.DepositKey(
		sdk.MustAccAddressFromBech32(deposit.DepositorAddress), 
		deposit.MustGetValidatorAddr(),
	), b)
}

// RemoveDeposit removes a deposit from the store
func (k Keeper) RemoveDeposit(ctx sdk.Context, deposit types.Deposit) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	store.Delete(types.DepositKey(
		sdk.MustAccAddressFromBech32(deposit.DepositorAddress),
		deposit.MustGetValidatorAddr(),
	))
}

// GetAllDeposit returns all deposits
func (k Keeper) GetAllDeposit(ctx sdk.Context) (list []types.Deposit) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Deposit
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetDeposits returns all deposits of a specific validator
func (k Keeper) GetValidatorDeposits(
	ctx sdk.Context, 
	valAddr sdk.ValAddress,
) (deposits []types.Deposit) {

	for _, deposit := range k.GetAllDeposit(ctx) {
		if deposit.ValidatorAddress == valAddr.String() {
			deposits = append(deposits, deposit)
		}
	}
	return deposits
}

// -------------------------------------------------------------------------------------------------
// Deposit pool
// -------------------------------------------------------------------------------------------------

// SetDepositPool set a specific depositPool in the store from its index
func (k Keeper) SetDepositPool(ctx sdk.Context, depositPool types.DepositPool) {
	valAddr, err := sdk.ValAddressFromBech32(depositPool.OperatorAddress)
	if err != nil {
		panic(err)
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositPoolKeyPrefix))
	b := k.cdc.MustMarshal(&depositPool)
	store.Set(types.DepositPoolKey(valAddr), b)
}

// GetDepositPool returns a depositPool from its index
func (k Keeper) GetDepositPool(
	ctx sdk.Context,
	validatorAddress sdk.ValAddress,
) (val types.DepositPool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositPoolKeyPrefix))

	b := store.Get(types.DepositPoolKey(validatorAddress))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDepositPool removes a depositPool from the store
func (k Keeper) RemoveDepositPool(
	ctx sdk.Context,
	operatorAddress sdk.ValAddress,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositPoolKeyPrefix))
	store.Delete(types.DepositPoolKey(
		operatorAddress,
	))
}

// GetAllDepositPool returns all depositPool
func (k Keeper) GetAllDepositPool(ctx sdk.Context) (list []types.DepositPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositPoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DepositPool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// AddDepPoolTokensAndShares adds the tokens and associated shares to the pool balance given a pool 
// and an amount of tokens.
func (k Keeper) AddDepPoolTokensAndShares(
	ctx sdk.Context,
	depositPool types.DepositPool,
	tokensToAdd sdk.Coin,
) (addedShares sdk.Dec) {

	var issuedShares sdk.Dec
	if depositPool.Shares.IsZero() {
		issuedShares = sdk.NewDecFromInt(tokensToAdd.Amount)
	} else {
		shares, err := depositPool.SharesFromTokens(tokensToAdd)
		if err != nil { panic(err) }
		issuedShares = shares
	}

	depositPool.Tokens = depositPool.Tokens.Add(tokensToAdd)
	depositPool.Shares = depositPool.Shares.Add(issuedShares)

	k.SetDepositPool(ctx, depositPool)

	return issuedShares
}

// Removes shares and associated tokens from a deposit pool. Return an error
// if the requested amounts are not available.
func (k Keeper) RemoveDepPoolTokensAndShares(
	ctx sdk.Context,
	depositPool types.DepositPool,
	sharesToRemove sdk.Dec,
) (types.DepositPool, sdkmath.Int) {
	var removedTokensAmt sdkmath.Int
	var remainingTokensAmt sdkmath.Int

	remainingShares := depositPool.Shares.Sub(sharesToRemove)

	if remainingShares.IsZero() {
		// Last delegation share gets any trimmings.
		removedTokensAmt = depositPool.Tokens.Amount
		remainingTokensAmt = sdk.ZeroInt()
	} else {
		removedTokensAmt = depositPool.TokensFromShares(sharesToRemove).TruncateInt()
		remainingTokensAmt = depositPool.Tokens.Amount.Sub(removedTokensAmt)
		if remainingTokensAmt.IsNegative() {
			panic("attempting to remove more tokens than available in the pool")
		}
	}

	depositPool.Tokens.Amount = remainingTokensAmt
	depositPool.Shares = remainingShares
	k.SetDepositPool(ctx, depositPool)

	return depositPool, removedTokensAmt
}
