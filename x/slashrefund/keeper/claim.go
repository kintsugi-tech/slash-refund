package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (k Keeper) Claim(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
) (refundCoins sdk.Coins, err error) {

	refund, found := k.GetRefund(ctx, delAddr, valAddr)
	if !found {
		return sdk.NewCoins(sdk.NewCoin(k.AllowedTokens(ctx)[0], sdk.NewInt(0))), types.ErrNoRefundForAddress
	}

	// Get the refund pool. If at this point it can't be found, then panic is called.
	// This is done because a refund cannot exists without the linked refund pool.
	// When refund il claimed it is removed from the store, and if it was the last
	// refund linked to the refund pool, then also the refund pool is removed. A
	// situation in which a refund is found but no linked refund pool can be found
	// is the result of a serious malfunction thus panic is called.
	refundPool, found := k.GetRefundPool(ctx, valAddr)
	if !found {
		panic("found refund but not the refund pool")
	}

	shares := refund.Shares
	if refund.Shares.GT(refundPool.Shares) {
		shares = refundPool.Shares

	}

	refundPool, drawnAmt := k.RemoveRefPoolTokensAndShares(ctx, refundPool, shares)

	refundCoins = sdk.NewCoins(sdk.NewCoin(k.AllowedTokens(ctx)[0], drawnAmt))
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delAddr, refundCoins)
	if err != nil {
		k.AddRefPoolTokensAndShares(ctx, refundPool, sdk.NewCoin(k.AllowedTokens(ctx)[0], drawnAmt))
		return sdk.NewCoins(sdk.NewCoin(k.AllowedTokens(ctx)[0], sdk.NewInt(0))), err
	}

	k.RemoveRefund(ctx, refund)
	if refundPool.Shares.IsZero() {
		k.RemoveRefundPool(ctx, valAddr)
	}

	return refundCoins, nil
}
