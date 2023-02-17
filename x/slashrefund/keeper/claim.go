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

	// 0. Get refund and refund pool
	zeroCoins := sdk.NewCoins(sdk.NewCoin(k.AllowedTokens(ctx)[0], sdk.NewInt(0)))

	refund, found := k.GetRefund(ctx, delAddr, valAddr)
	if !found {
		return zeroCoins, types.ErrNoRefundForAddress
	}

	refundPool, found := k.GetRefundPool(ctx, valAddr)
	if !found {
		panic("found refund but not the refund pool")
	}

	shares := refund.Shares
	if refund.Shares.GT(refundPool.Shares) {
		shares = refundPool.Shares

	}

	// 2. Remove shares and tokens from redundPool
	refundPool, drawnAmt := k.RemoveRefPoolTokensAndShares(ctx, refundPool, shares)

	// 3. Send coins to delegator
	refundCoins = sdk.NewCoins(sdk.NewCoin(k.AllowedTokens(ctx)[0], drawnAmt))
	senderModule := types.ModuleName
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, senderModule, delAddr, refundCoins)
	if err != nil {
		k.AddRefPoolTokensAndShares(ctx, refundPool, sdk.NewCoin(k.AllowedTokens(ctx)[0], drawnAmt))
		return zeroCoins, err
	}

	k.RemoveRefund(ctx, refund)
	if refundPool.Shares.IsZero() {
		k.RemoveRefundPool(ctx, valAddr)
	}

	return refundCoins, nil
}
