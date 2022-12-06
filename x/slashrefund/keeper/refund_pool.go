package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// SetRefundPool set a specific refundPool in the store from its index
func (k Keeper) SetRefundPool(ctx sdk.Context, refundPool types.RefundPool) {

	valAddr, err := sdk.ValAddressFromBech32(refundPool.OperatorAddress)
	if err != nil {
		panic(err)
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	b := k.cdc.MustMarshal(&refundPool)
	store.Set(types.RefundPoolKey(valAddr), b)
}

// GetRefundPool returns a refundPool from its index
func (k Keeper) GetRefundPool(
	ctx sdk.Context,
	valAddr sdk.ValAddress,

) (refPool types.RefundPool, found bool) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	b := store.Get(types.RefundPoolKey(valAddr))
	if b == nil {
		return refPool, false
	}
	k.cdc.MustUnmarshal(b, &refPool)
	return refPool, true
}

// RemoveRefundPool removes a refundPool from the store
func (k Keeper) RemoveRefundPool(
	ctx sdk.Context,
	valAddr sdk.ValAddress,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	store.Delete(types.RefundPoolKey(valAddr))
}

// GetAllRefundPool returns all refundPool
func (k Keeper) GetAllRefundPool(ctx sdk.Context) (list []types.RefundPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var refPool types.RefundPool
		k.cdc.MustUnmarshal(iterator.Value(), &refPool)
		list = append(list, refPool)
	}

	return
}

func (k Keeper) AddRefPoolTokensAndShares(
	ctx sdk.Context,
	refundPool types.RefundPool,
	tokensToAdd sdk.Coin,
) (addedShares sdk.Dec) {

	var issuedShares sdk.Dec
	if refundPool.Shares.IsZero() {
		issuedShares = sdk.NewDecFromInt(tokensToAdd.Amount)
	} else {
		// TODO: we have to manage post slashing send of tokens. We have to put zero shares when  tokens -> 0
		shares, err := refundPool.SharesFromTokens(tokensToAdd)
		if err != nil {
			panic(err)
		}
		issuedShares = shares
	}

	refundPool.Tokens = refundPool.Tokens.Add(tokensToAdd)
	refundPool.Shares = refundPool.Shares.Add(issuedShares)

	k.SetRefundPool(ctx, refundPool)

	return issuedShares
}

func (k Keeper) RemoveRefPoolTokensAndShares(
	ctx sdk.Context,
	refundPool types.RefundPool,
	sharesToRemove sdk.Dec,
) (types.RefundPool, sdk.Int) {

	var issuedTokensAmt sdk.Int

	remainingShares := refundPool.Shares.Sub(sharesToRemove)

	if remainingShares.IsZero() {
		// last share gets any trimmings
		issuedTokensAmt = refundPool.Tokens.Amount
		// TODO: generalize it considering AllowedTokens param
		refundPool.Tokens.Amount = sdk.ZeroInt()
	} else {
		// leave excess tokens in the deposit pool
		// however fully use all the depositor shares
		issuedTokensAmt = refundPool.TokensFromShares(sharesToRemove).TruncateInt()
		refundPool.Tokens.Amount = refundPool.Tokens.Amount.Sub(issuedTokensAmt)
		if refundPool.Tokens.Amount.IsNegative() {
			panic("attempting to remove more tokens than available in validator")
		}
	}

	refundPool.Shares = remainingShares
	k.SetRefundPool(ctx, refundPool)

	return refundPool, issuedTokensAmt
}
