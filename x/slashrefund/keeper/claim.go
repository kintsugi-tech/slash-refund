package keeper

import (
	//"fmt"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (k Keeper) Claim(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	shares sdk.Dec,
) (refundCoins sdk.Coins, err error) {

	// 0. Get refund and refund pool
	// Since this function is normally called after ValidateClaimAmount,
	// refund and refundPool should be available as ValidateClaimAmount
	// checks for their presence.
	zeroCoins := sdk.NewCoins(sdk.NewCoin(k.AllowedTokensList(ctx)[0], sdk.NewInt(0)))
	refund, found := k.GetRefund(ctx, delAddr, valAddr)
	if !found {
		return zeroCoins, types.ErrNoRefundForAddress
	}
	refundPool, found := k.GetRefundPool(ctx, valAddr)
	if !found {
		return zeroCoins, types.ErrNoRefundPoolForValidator
	}
	// ensure that we have enough shares to remove
	if refund.Shares.LT(shares) {
		return zeroCoins, sdkerrors.Wrap(types.ErrNotEnoughRefundShares, refund.Shares.String())
	}

	// 1. Remove shares from refund
	refund.Shares = refund.Shares.Sub(shares)
	// remove the deposit if zero or set a new doposit
	if refund.Shares.IsZero() {
		k.RemoveRefund(ctx, refund)
	} else {
		k.SetRefund(ctx, refund)
	}

	// 2. Remove shares and tokens from redundPool
	refundPool, drawnAmt := k.RemoveRefPoolTokensAndShares(ctx, refundPool, shares)
	if refundPool.Shares.IsZero() {
		k.RemoveRefundPool(ctx, valAddr)
	}

	// 3. Send coins to delegator
	refundCoin := sdk.NewCoin(k.AllowedTokensList(ctx)[0], drawnAmt)
	refundCoins = sdk.NewCoins(refundCoin)
	senderModule := types.ModuleName
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, senderModule, delAddr, refundCoins)
	if err != nil {
		return zeroCoins, err
	}

	logger := k.Logger(ctx)
	logger.Error(fmt.Sprintf("Claim completed: %s%s, delegator %s, validator %s", refundCoin.Amount.String(), refundCoin.Denom, delAddr.String(), valAddr.String()))
	return refundCoins, nil
}

func (k Keeper) ValidateClaimAmount(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	tokens sdk.Coin,
) (shares sdk.Dec, err error) {

	logger := k.Logger(ctx)
	logger.Error(fmt.Sprintf("Validating claim amount: %s%s, delegator %s, validator %s", tokens.Amount.String(), tokens.Denom, delAddr.String(), valAddr.String()))

	isValid, err := k.CheckAllowedTokens(ctx, tokens.Denom)
	if !isValid {
		logger.Error("ERROR: Invalid denom")
		return sdk.NewDec(0), err
	}

	if tokens.Amount.IsZero() {
		logger.Error("ERROR: Zero claim")
		return sdk.NewDec(0), types.ErrZeroWithdraw
	}

	refund, found := k.GetRefund(ctx, delAddr, valAddr)
	if !found {
		logger.Error("ERROR: No refund found")
		return sdk.NewDec(0), types.ErrNoRefundForAddress
	}

	refPool, found := k.GetRefundPool(ctx, valAddr)
	if !found {
		logger.Error("ERROR: No refund pool found")
		return sdk.NewDec(0), types.ErrNoRefundPoolForValidator
	}

	// compute shares from wanted withdraw tokens
	shares, err = refPool.SharesFromTokens(tokens)
	if err != nil {
		logger.Error("ERROR: refPool.SharesFromTokens: ErrZeroTokensQuotient.")
		return sdk.NewDec(0), err
	}

	// compute shares from wanted withdraw tokens, rounded down
	sharesTruncated, err := refPool.SharesFromTokensTruncated(tokens)
	if err != nil {
		logger.Error("ERROR: refPool.SharesFromTokensTruncated: ErrZeroTokensQuotient.")
		return sdk.NewDec(0), err
	}

	// check if wanted tokens converted to truncated shares are greater than actual total of delegator shares
	delegatorShares := refund.GetShares()
	if sharesTruncated.GT(delegatorShares) {
		logger.Error("ERROR: Invalid claim amount.")
		return sdk.NewDec(0), sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token amount")
	}

	// cap shares (not-truncated) at total depositor shares
	if shares.GT(delegatorShares) {
		shares = delegatorShares
	}

	logger.Error("Validating claim amount: validated.")

	return shares, nil
}
