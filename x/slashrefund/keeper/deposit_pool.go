package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// SetDepositPool set a specific depositPool in the store from its index
func (k Keeper) SetDepositPool(ctx sdk.Context, depositPool types.DepositPool) {
	valOperAddr, err := sdk.ValAddressFromBech32(depositPool.OperatorAddress)
	if err != nil {
		panic(err)
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositPoolKeyPrefix))
	b := k.cdc.MustMarshal(&depositPool)
	store.Set(types.DepositPoolKey(
		valOperAddr,
	), b)
}

// GetDepositPool returns a depositPool from its index
func (k Keeper) GetDepositPool(
	ctx sdk.Context,
	operatorAddress sdk.ValAddress,

) (val types.DepositPool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositPoolKeyPrefix))

	b := store.Get(types.DepositPoolKey(
		operatorAddress,
	))
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

func (k Keeper) AddDepPoolTokensAndShares(
	ctx sdk.Context,
	depositPool types.DepositPool,
	tokensToAdd sdk.Coin,
) (addedShares sdk.Dec) {

	var issuedShares sdk.Dec
	if depositPool.Shares.IsZero() {
		issuedShares = sdk.NewDecFromInt(tokensToAdd.Amount)
	} else {
		// TODO: we have to manage post slashing send of tokens. We have to put zero shares when  tokens -> 0
		shares, err := depositPool.SharesFromTokens(tokensToAdd)
		if err != nil {
			panic(err)
		}
		issuedShares = shares

	}

	depositPool.Tokens = depositPool.Tokens.Add(tokensToAdd)
	depositPool.Shares = depositPool.Shares.Add(issuedShares)

	k.SetDepositPool(ctx, depositPool)

	return issuedShares
}

func (k Keeper) RemoveDepPoolTokensAndShares(
	ctx sdk.Context,
	depositPool types.DepositPool,
	sharesToRemove sdk.Dec,
) (types.DepositPool, sdk.Int) {

	var issuedTokensAmt sdk.Int

	remainingShares := depositPool.Shares.Sub(sharesToRemove)

	if remainingShares.IsZero() {
		// last delegation share gets any trimmings
		issuedTokensAmt = depositPool.Tokens.Amount
		// TODO: generalize it considering AllowedTokens param
		depositPool.Tokens.Amount = sdk.ZeroInt()
	} else {
		// leave excess tokens in the deposit pool
		// however fully use all the depositor shares
		issuedTokensAmt = depositPool.TokensFromShares(sharesToRemove).TruncateInt()
		depositPool.Tokens.Amount = depositPool.Tokens.Amount.Sub(issuedTokensAmt)
		if depositPool.Tokens.Amount.IsNegative() {
			panic("attempting to remove more tokens than available in validator")
		}
	}

	depositPool.Shares = remainingShares
	k.SetDepositPool(ctx, depositPool)

	return depositPool, issuedTokensAmt
}
